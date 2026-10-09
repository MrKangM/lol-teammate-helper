// Package connector keeps a WebSocket connection to the League client alive and
// forwards its events to the dispatcher.
package connector

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"

	"lol-teammate-helper/internal/config"
	"lol-teammate-helper/internal/dispatch"
	"lol-teammate-helper/internal/lcu"
	"lol-teammate-helper/internal/service"
	"lol-teammate-helper/internal/types"
)

const (
	wsURLTemplate = "wss://127.0.0.1:%d"

	// Subscribing to specific events avoids receiving every LCU event.
	champSelectTopic = "OnJsonApiEvent_lol-champ-select_v1_session"
	gameflowTopic    = "OnJsonApiEvent_lol-gameflow_v1_session"
	eventPrefix      = "OnJsonApiEvent"

	opSubscribe = 5

	minBackoff = time.Second
	maxBackoff = 10 * time.Second
)

// Run detects the League client and keeps a websocket connected to it until ctx
// is cancelled. Credentials are re-detected on every attempt because the client
// picks a new port and token each time it starts.
func Run(ctx context.Context) {
	backoff := minBackoff
	for ctx.Err() == nil {
		creds, err := lcu.Detect()
		if err != nil {
			slog.Debug("league client not detected", "err", err)
		} else {
			if config.Update(creds.Port, creds.Token, creds.Region) {
				service.Shared().ResetCaches()
			}
			start := time.Now()
			if err := session(ctx, creds.Port); err != nil {
				slog.Warn("websocket session ended", "err", err)
			}
			if time.Since(start) > 30*time.Second {
				backoff = minBackoff // it was a healthy connection; retry promptly
			}
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(backoff):
		}
		if backoff *= 2; backoff > maxBackoff {
			backoff = maxBackoff
		}
	}
}

func session(ctx context.Context, port int) error {
	cfg, ok := config.Instance()
	if !ok {
		return fmt.Errorf("config not initialised")
	}

	dialer := websocket.Dialer{
		TLSClientConfig:  &tls.Config{InsecureSkipVerify: true},
		HandshakeTimeout: 10 * time.Second,
	}
	headers := http.Header{}
	headers.Set("Authorization", cfg.Token)

	conn, _, err := dialer.DialContext(ctx, fmt.Sprintf(wsURLTemplate, port), headers)
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}
	defer conn.Close()

	// Unblock ReadMessage when the app is shutting down.
	go func() {
		<-ctx.Done()
		_ = conn.Close()
	}()

	for _, topic := range []string{champSelectTopic, gameflowTopic} {
		if err := conn.WriteJSON([]interface{}{opSubscribe, topic}); err != nil {
			return fmt.Errorf("subscribe %s: %w", topic, err)
		}
	}
	slog.Info("connected to league client", "port", port)

	// A single worker handles events in order. While it is busy only the latest
	// event per resource is kept, since each update supersedes the previous one.
	q := newEventQueue()
	done := make(chan struct{})
	go func() {
		defer close(done)
		for msg := range q.out() {
			dispatch.EventHandler(msg.Uri, msg.EventType, msg.Data)
		}
	}()
	defer func() {
		q.close()
		<-done
	}()

	// Pick up a game that was already in progress before we connected.
	go bootstrap(cfg, q)

	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			return fmt.Errorf("read: %w", err)
		}
		msg, ok := parseMessage(raw)
		if !ok {
			continue
		}
		q.push(msg)
	}
}

// bootstrap feeds the current state of each resource into the queue as if it
// had just been pushed, so nothing is missed when the app starts mid-game.
func bootstrap(cfg *config.AppConfig, q *eventQueue) {
	for _, uri := range []string{dispatch.GameflowURI, dispatch.ChampSelectURI} {
		body, err := cfg.SendHttpRequest(uri, http.MethodGet)
		if err != nil {
			continue // e.g. 404 while not in champ select
		}
		q.push(types.WSMessageType{Uri: uri, EventType: "Update", Data: body})
	}
}

// parseMessage decodes a WAMP-style frame [opcode, topic, payload].
func parseMessage(raw []byte) (types.WSMessageType, bool) {
	var frame []json.RawMessage
	if err := json.Unmarshal(raw, &frame); err != nil || len(frame) < 3 {
		return types.WSMessageType{}, false
	}

	var topic string
	if err := json.Unmarshal(frame[1], &topic); err != nil || !strings.HasPrefix(topic, eventPrefix) {
		return types.WSMessageType{}, false
	}

	var msg types.WSMessageType
	if err := json.Unmarshal(frame[2], &msg); err != nil {
		slog.Warn("decode event payload failed", "topic", topic, "err", err)
		return types.WSMessageType{}, false
	}
	return msg, true
}
