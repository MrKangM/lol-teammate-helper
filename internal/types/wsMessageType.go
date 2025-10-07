package types

import "encoding/json"

type WSMessageType struct {
	Data      json.RawMessage `json:"data"`
	EventType string          `json:"eventType"`
	Uri       string          `json:"uri"`
}
