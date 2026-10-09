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

	// Subscribing to a specific event avoids receiving every LCU event.
	champSelectTopic = "OnJsonApiEvent_lol-champ-select_v1_session"
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

	if err := conn.WriteJSON([]interface{}{opSubscribe, champSelectTopic}); err != nil {
		return fmt.Errorf("subscribe: %w", err)
	}
	slog.Info("connected to league client", "port", port)

	// A single worker handles events in order; only the latest pending event is
	// kept while it is busy, since each champ select update supersedes the last.
	events := make(chan types.WSMessageType, 1)
	done := make(chan struct{})
	go func() {
		defer close(done)
		for msg := range events {
			dispatch.EventHandler(msg.Uri, msg.EventType, msg.Data)
		}
	}()
	defer func() {
		close(events)
		<-done
	}()

	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			return fmt.Errorf("read: %w", err)
		}
		msg, ok := parseMessage(raw)
		if !ok {
			continue
		}
		offerLatest(events, msg)
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

// offerLatest queues msg, replacing a stale queued message if the worker is behind.
func offerLatest(ch chan types.WSMessageType, msg types.WSMessageType) {
	for {
		select {
		case ch <- msg:
			return
		default:
		}
		select {
		case <-ch: // drop the stale one
		default:
		}
	}
}
