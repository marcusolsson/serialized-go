package serialized

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFeed(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
  "entries": [
    {
      "sequenceNumber": 12314,
      "aggregateId": "22c3780f-6dcb-440f-8532-6693be83f21c",
      "timestamp": 1503386583474,
      "events": [
        {
          "eventId": "f2c8bfc1-c702-4f1a-b295-ef113ed7c8be",
          "eventType": "PaymentProcessed",
          "data": {
            "paymentMethod": "CARD",
            "amount": 1000,
            "currency": "SEK"
          }
        }
      ]
    }
  ],
  "hasMore": false
}`))
	}))

	c := NewClient(
		WithBaseURL(ts.URL),
	)

	f, err := c.Feed("payment")
	if err != nil {
		t.Fatal(err)
	}

	if len(f.Entries) != 1 {
		t.Fatalf("unexpected number of entries = %d; want = %d", len(f.Entries), 1)
	}
	entry := f.Entries[0]

	if len(entry.Events) != 1 {
		t.Fatalf("unexpected number of events = %d; want = %d", len(entry.Events), 1)
	}
	event := entry.Events[0]

	var pp testPaymentProcessed
	json.Unmarshal(event.Data, &pp)

	if pp.Amount != 1000 {
		t.Errorf("incorrect amount = %d; want = %d", pp.Amount, 1000)
	}
}
