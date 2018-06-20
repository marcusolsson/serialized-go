package main

import (
	"context"

	serialized "github.com/marcusolsson/serialized-go"
	uuid "github.com/satori/go.uuid"
)

func eventsStoreHandler(c *serialized.Client, aggType, aggID, eventType, eventID, data string, version int64) error {
	if eventID == "" {
		id, err := uuid.NewV4()
		if err != nil {
			return err

		}
		eventID = id.String()
	}

	if aggID == "" {
		id, err := uuid.NewV4()
		if err != nil {
			return err
		}
		aggID = id.String()
	}

	event := &serialized.Event{
		Type: eventType,
		ID:   eventID,
		Data: []byte(data),
	}

	return c.Store(context.Background(), aggType, aggID, version, event)
}
