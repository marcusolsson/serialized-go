package serialized

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestFeedList(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := loadJSON("testdata/feed_list_response.json")
		if err != nil {
			t.Fatal(err)
		}
		w.Write(b)
	}))

	c := NewClient(
		WithBaseURL(ts.URL),
	)

	fis, err := c.Feeds(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if len(fis) != 1 {
		t.Fatalf("unexpected number of entries = %d; want = %d", len(fis), 1)
	}
	got := fis[0]

	want := FeedInfo{
		AggregateType:  "payment",
		AggregateCount: 1337,
		BatchCount:     7331,
		EventCount:     9977,
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got = %v; want = %v", got, want)
	}
}

func TestFeed(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := loadJSON("testdata/feed_log_response.json")
		if err != nil {
			t.Fatal(err)
		}
		w.Write(b)
	}))

	c := NewClient(
		WithBaseURL(ts.URL),
	)

	f, err := c.feed(context.Background(), "payment", 0)
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
