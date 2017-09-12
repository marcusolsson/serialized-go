package serialized

import "encoding/json"

// Event holds a Serialized.io event.
type Event struct {
	ID            string          `json:"eventId"`
	Type          string          `json:"eventType"`
	Data          json.RawMessage `json:"data,omitempty"`
	EncryptedData string          `json:"encryptedData,omitempty"`
}

// NewEvent is a helper function for marshaling event data to JSON.
func NewEvent(eventID, eventType string, data interface{}) Event {
	b, _ := json.Marshal(data)
	return Event{
		Type: eventType,
		ID:   eventID,
		Data: b,
	}
}
