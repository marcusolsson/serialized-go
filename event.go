package serialized

import "encoding/json"

// Event represents a Serialized.io event.
type Event struct {
	ID            string          `json:"eventId"`
	Type          string          `json:"eventType"`
	Data          json.RawMessage `json:"data,omitempty"`
	EncryptedData string          `json:"encryptedData,omitempty"`
}
