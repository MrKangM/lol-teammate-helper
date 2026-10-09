package diag

import (
	"strings"
	"testing"
)

func TestRingBufferKeepsNewestFirst(t *testing.T) {
	Reset()
	for i := 0; i < maxEvents+5; i++ {
		RecordEvent("/u", string(rune('a'+i%26)), []byte("x"))
	}
	got := Get().Events
	if len(got) != maxEvents {
		t.Fatalf("len = %d, want %d", len(got), maxEvents)
	}
	if got[0].Time.Before(got[len(got)-1].Time) {
		t.Fatal("events should be newest first")
	}
}

func TestEventTruncation(t *testing.T) {
	Reset()
	RecordEvent("/u", "Update", []byte(strings.Repeat("a", maxEventPayload+10)))
	ev := Get().Events[0]
	if !ev.Truncated || len(ev.Data) != maxEventPayload || ev.Bytes != maxEventPayload+10 {
		t.Fatalf("unexpected event %+v", ev.Bytes)
	}
}

func TestConnectionState(t *testing.T) {
	Reset()
	SetDisconnected("boom")
	if s := Get(); s.Connected || s.LastError != "boom" {
		t.Fatalf("got %+v", s)
	}
	SetConnected(1234, "艾欧尼亚")
	if s := Get(); !s.Connected || s.Port != 1234 || s.LastError != "" {
		t.Fatalf("got %+v", s)
	}
}
