// +build integration

package test

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"

	serialized "github.com/marcusolsson/serialized-go"
)

func TestAPI(t *testing.T) {
	client := serialized.NewClient(
		serialized.WithAccessKey(os.Getenv("SERIALIZED_ACCESS_KEY")),
		serialized.WithSecretAccessKey(os.Getenv("SERIALIZED_SECRET_ACCESS_KEY")),
	)

	ctx := context.Background()

	aggType := "payment_" + uuid.New().String()
	aggID := uuid.New().String()

	eventData, _ := json.Marshal(map[string]interface{}{
		"paymentMethod": "CARD",
		"amount":        1000,
		"currency":      "SEK",
	})

	event := &serialized.Event{
		ID:            uuid.New().String(),
		Type:          "PaymentProcessed",
		Data:          eventData,
		EncryptedData: "string",
	}

	if err := client.Store(ctx, aggType, aggID, 0, event); err != nil {
		t.Fatal(err)
	}

	exists, err := client.AggregateExists(ctx, aggType, aggID)
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatalf("aggregate does not exist")
	}

	agg, err := client.LoadAggregate(ctx, aggType, aggID)
	if err != nil {
		t.Fatal(err)
	}

	if agg.Type != aggType {
		t.Errorf("agg.Type = %q; want = %q", agg.Type, aggType)
	}
	if agg.ID != aggID {
		t.Errorf("agg.ID = %q; want = %q", agg.ID, aggID)
	}
	if agg.Version != 1 {
		t.Errorf("agg.Version = %q; want = %q", agg.Version, 1)
	}
	if len(agg.Events) != 1 {
		t.Errorf("number of events = %d; want = %d", len(agg.Events), 1)
	}

	ev := agg.Events[0]

	if ev.ID != event.ID {
		t.Errorf("unexpected id in event = %q; want = %q", string(ev.ID), string(event.ID))
	}
	if ev.Type != event.Type {
		t.Errorf("unexpected type in event = %q; want = %q", string(ev.Type), string(event.Type))
	}
	if !equalsJSON(ev.Data, event.Data) {
		t.Errorf("unexpected data in event = %q; want = %q", string(ev.Data), string(event.Data))
	}
	if ev.EncryptedData != event.EncryptedData {
		t.Errorf("unexpected encrypted data in event = %q; want = %q", ev.EncryptedData, event.EncryptedData)
	}

	// Feed
	feeds, err := client.Feeds(ctx)
	if err != nil {
		t.Fatal(err)
	}

	var foundFeed bool
	for _, f := range feeds {
		if f.AggregateType != aggType {
			continue
		}

		foundFeed = true

		if f.AggregateCount != 1 {
			t.Errorf("feed aggregate count = %d; want = %d", f.AggregateCount, 1)
		}
		if f.BatchCount != 1 {
			t.Errorf("feed batch count = %d; want = %d", f.BatchCount, 1)
		}
		if f.EventCount != 1 {
			t.Errorf("feed event count = %d; want = %d", f.EventCount, 1)
		}
	}

	if !foundFeed {
		t.Errorf("feed was not found = %q", aggType)
	}

	seq, err := client.FeedSequenceNumber(ctx, aggType)
	if err != nil {
		t.Fatal(err)
	}
	if seq != 1 {
		t.Errorf("sequence number = %d; want = %d", seq, 1)
	}

	feedctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = client.Feed(feedctx, aggType, 0, func(e *serialized.FeedEntry) {
		if e.AggregateID != aggID {
			t.Errorf("feed entry aggregate id = %q; want = %q", e.AggregateID, aggID)
		}
		if len(e.Events) != 1 {
			t.Errorf("feed entry events = %d; want = %d", len(e.Events), 1)
		}
		if e.SequenceNumber != 1 {
			t.Errorf("feed entry sequence number = %d; want = %d", e.SequenceNumber, 1)
		}
		cancel()
	})
	if err != context.Canceled {
		t.Fatal(err)
	}

	tok, err := client.RequestDeleteAggregateByType(ctx, aggType)
	if err != nil {
		t.Fatal(err)
	}

	if err := client.DeleteAggregateByType(ctx, aggType, tok); err != nil {
		t.Fatal(err)
	}

	exists, err = client.AggregateExists(ctx, aggType, aggID)
	if err != nil {
		t.Fatal(err)
	}
	if exists {
		t.Fatalf("aggregate was not properly cleaned up")
	}
}
