package serialized

import (
	"bytes"
	"encoding/json"
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
		want := []byte(`{"aggregateId":"123","events":[{"eventId":"456","eventType":"PaymentProcessed","data":{"paymentMethod":"CARD","amount":1000,"currency":"SEK"}}]}`)

		got, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}

		got = bytes.TrimSpace(got)

		if !bytes.Equal(got, want) {
			var wantIndent bytes.Buffer
			json.Indent(&wantIndent, want, "", "\t")

			var gotIndent bytes.Buffer
			json.Indent(&gotIndent, got, "", "\t")

			t.Errorf("unexpected request body =\n%s\n\nwant =\n%s", got, want)
		}
	}))

	c := NewClient(
		WithBaseURL(ts.URL),
	)

	pp := testPaymentProcessed{
		PaymentMethod: "CARD",
		Amount:        1000,
		Currency:      "SEK",
	}

	err := c.Store("payment", "123", 1, NewEvent("456", "PaymentProcessed", pp))
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoadAggregate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
  "aggregateId": "22c3780f-6dcb-440f-8532-6693be83f21c",
  "aggregateVersion": 1,
  "aggregateType": "payment",
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
}`))
	}))

	c := NewClient(
		WithBaseURL(ts.URL),
	)

	var (
		aggType    = "payment"
		aggID      = "22c3780f-6dcb-440f-8532-6693be83f21c"
		aggVersion = 1
	)

	agg, err := c.LoadAggregate(aggType, aggID)
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
