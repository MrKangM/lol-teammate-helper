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

func TestOfferLatestKeepsNewest(t *testing.T) {
	ch := make(chan types.WSMessageType, 1)
	offerLatest(ch, types.WSMessageType{Uri: "old"})
	offerLatest(ch, types.WSMessageType{Uri: "new"})
	if got := <-ch; got.Uri != "new" {
		t.Fatalf("got %q, want new", got.Uri)
	}
}
