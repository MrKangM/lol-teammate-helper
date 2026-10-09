// Package diag keeps a small in-memory picture of the connection to the League
// client and the most recent raw events, so the UI can show what is going on
// without anyone having to dig through log files.
package diag

import (
	"sync"
	"time"
)

const (
	maxEvents       = 40
	maxEventPayload = 30000
)

// Event is one raw LCU event as received.
type Event struct {
	Time      time.Time `json:"time"`
	URI       string    `json:"uri"`
	Type      string    `json:"type"`
	Bytes     int       `json:"bytes"`
	Truncated bool      `json:"truncated"`
	Data      string    `json:"data"`
}

// Snapshot is what the diagnostics page displays.
type Snapshot struct {
	Connected   bool      `json:"connected"`
	Port        int       `json:"port"`
	Region      string    `json:"region"`
	LastError   string    `json:"lastError"`
	ConnectedAt time.Time `json:"connectedAt"`
	LogPath     string    `json:"logPath"`
	Phase       string    `json:"phase"`
	Events      []Event   `json:"events"` // newest first
}

var state struct {
	mu   sync.Mutex
	snap Snapshot
	ring []Event
}

// SetLogPath records where the log file lives.
func SetLogPath(path string) {
	state.mu.Lock()
	state.snap.LogPath = path
	state.mu.Unlock()
}

// SetConnected marks a successful connection.
func SetConnected(port int, region string) {
	state.mu.Lock()
	state.snap.Connected = true
	state.snap.Port = port
	state.snap.Region = region
	state.snap.LastError = ""
	state.snap.ConnectedAt = time.Now()
	state.mu.Unlock()
}

// SetDisconnected marks the connection as lost, with the reason.
func SetDisconnected(reason string) {
	state.mu.Lock()
	state.snap.Connected = false
	state.snap.LastError = reason
	state.mu.Unlock()
}

// SetPhase records the latest gameflow phase.
func SetPhase(phase string) {
	state.mu.Lock()
	state.snap.Phase = phase
	state.mu.Unlock()
}

// RecordEvent stores a raw event in the ring buffer.
func RecordEvent(uri, eventType string, data []byte) {
	ev := Event{Time: time.Now(), URI: uri, Type: eventType, Bytes: len(data)}
	if len(data) > maxEventPayload {
		ev.Data, ev.Truncated = string(data[:maxEventPayload]), true
	} else {
		ev.Data = string(data)
	}

	state.mu.Lock()
	state.ring = append(state.ring, ev)
	if len(state.ring) > maxEvents {
		state.ring = state.ring[len(state.ring)-maxEvents:]
	}
	state.mu.Unlock()
}

// Get returns a copy of the current state, newest events first.
func Get() Snapshot {
	state.mu.Lock()
	defer state.mu.Unlock()

	out := state.snap
	out.Events = make([]Event, len(state.ring))
	for i, ev := range state.ring {
		out.Events[len(state.ring)-1-i] = ev
	}
	return out
}

// Reset clears everything except the log path (used by tests).
func Reset() {
	state.mu.Lock()
	path := state.snap.LogPath
	state.snap = Snapshot{LogPath: path}
	state.ring = nil
	state.mu.Unlock()
}
