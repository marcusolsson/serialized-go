package serialized

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testPaymentProcessed struct {
	PaymentMethod string `json:"paymentMethod"`
	Amount        int    `json:"amount"`
	Currency      string `json:"currency"`
}

func TestStore(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		want, err := loadJSON("testdata/event_store_request.json")
		if err != nil {
			t.Fatal(err)
		}
		assertEqualJSON(t, got, want)
	}))

	c := NewClient(
		WithBaseURL(ts.URL),
	)

	pp := testPaymentProcessed{
		PaymentMethod: "CARD",
		Amount:        1000,
		Currency:      "SEK",
	}

	ev := &Event{
		ID:            "f2c8bfc1-c702-4f1a-b295-ef113ed7c8be",
		Type:          "PaymentProcessed",
		Data:          mustMarshal(pp),
		EncryptedData: "string",
	}

	if err := c.Store(context.Background(), "payment", "2c3cf88c-ee88-427e-818a-ab0267511c84", 1, ev); err != nil {
		t.Fatal(err)
	}
}

func TestAggregateExists(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))

	c := NewClient(
		WithBaseURL(ts.URL),
	)

	var (
		aggType = "payment"
		aggID   = "22c3780f-6dcb-440f-8532-6693be83f21c"
	)

	exists, err := c.AggregateExists(context.Background(), aggType, aggID)
	if err != nil {
		t.Fatal(err)
	}

	if !exists {
		t.Fatal("aggregate should exist")
	}
}

func TestLoadAggregate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := loadJSON("testdata/event_load_response.json")
		if err != nil {
			t.Fatal(err)
		}
		w.Write(b)
	}))

	c := NewClient(
		WithBaseURL(ts.URL),
	)

	var (
		aggType    = "payment"
		aggID      = "22c3780f-6dcb-440f-8532-6693be83f21c"
		aggVersion = 1
	)

	agg, err := c.LoadAggregate(context.Background(), aggType, aggID)
	if err != nil {
		t.Fatal(err)
	}

	if agg.Type != aggType {
		t.Errorf("unexpected type = %s; want = %s", agg.Type, aggType)
	}
	if agg.ID != aggID {
		t.Errorf("unexpected ID = %s; want = %s", agg.ID, aggID)
	}
	if agg.Version != aggVersion {
		t.Errorf("unexpected version = %d; want = %d", agg.Version, aggVersion)
	}
	if len(agg.Events) != 1 {
		t.Errorf("unexpected number of events = %d; want = %d", len(agg.Events), 1)
	}
}
