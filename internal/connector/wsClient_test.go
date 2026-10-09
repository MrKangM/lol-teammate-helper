package connector

import (
	"testing"

	"lol-teammate-helper/internal/types"
)

func TestParseMessage(t *testing.T) {
	raw := []byte(`[8,"OnJsonApiEvent_lol-champ-select_v1_session",{"data":{"id":"x"},"eventType":"Update","uri":"/lol-champ-select/v1/session"}]`)
	msg, ok := parseMessage(raw)
	if !ok || msg.Uri != "/lol-champ-select/v1/session" || msg.EventType != "Update" || len(msg.Data) == 0 {
		t.Fatalf("unexpected result: %+v ok=%v", msg, ok)
	}
}

func TestParseMessageIgnoresOthers(t *testing.T) {
	for _, raw := range []string{`not json`, `[8]`, `[8,"SomethingElse",{}]`, `[8,"OnJsonApiEvent",5]`} {
		if _, ok := parseMessage([]byte(raw)); ok {
			t.Errorf("expected %q to be ignored", raw)
		}
	}
}

func TestEventQueueKeepsLatestPerURI(t *testing.T) {
	q := newEventQueue()
	q.push(types.WSMessageType{Uri: "a", EventType: "old"})
	q.push(types.WSMessageType{Uri: "b", EventType: "b1"})
	q.push(types.WSMessageType{Uri: "a", EventType: "new"})

	first, _ := q.pop()
	second, _ := q.pop()
	if first.Uri != "a" || first.EventType != "new" || second.Uri != "b" {
		t.Fatalf("got %+v then %+v", first, second)
	}
	if _, ok := q.pop(); ok {
		t.Fatal("queue should be empty")
	}
}

func TestEventQueueDrainsOnClose(t *testing.T) {
	q := newEventQueue()
	q.push(types.WSMessageType{Uri: "a"})
	q.close()

	var got int
	for range q.out() {
		got++
	}
	if got != 1 {
		t.Fatalf("delivered %d messages, want 1", got)
	}
}
