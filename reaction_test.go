package serialized

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestListReactions(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := loadJSON("testdata/reaction_list_response.json")
		if err != nil {
			t.Fatal(err)
		}
		w.Write(b)
	}))

	c := NewClient(
		WithBaseURL(ts.URL),
	)

	r, err := c.ListReactionDefinitions(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if len(r) != 1 {
		t.Errorf("unexpected number of reactions = %d; want = %d", len(r), 1)
	}
}

func TestCreateReaction(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		want, err := loadJSON("testdata/reaction_create_request.json")
		if err != nil {
			t.Fatal(err)
		}
		assertEqualJSON(t, got, want)

		w.WriteHeader(http.StatusOK)
	}))

	c := NewClient(
		WithBaseURL(ts.URL),
	)

	r := &ReactionDefinition{
		Name:               "payment-processed-email-reaction",
		Feed:               "payment",
		ReactOnEventType:   "PaymentProcessed",
		CancelOnEventTypes: []string{"OrderCanceledEvent"},
		TriggerTimeField:   "my.event.data.field",
		Offset:             "PT1H",
		Action: &Action{
			ActionType: ActionTypeHTTPPost,
		},
	}

	if err := c.CreateReactionDefinition(context.Background(), r); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteReaction(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))

	c := NewClient(
		WithBaseURL(ts.URL),
	)

	if err := c.DeleteReactionDefinition(context.Background(), "payment-processed-email-reaction"); err != nil {
		t.Fatal(err)
	}
}

func TestGetReaction(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := loadJSON("testdata/reaction_get_response.json")
		if err != nil {
			t.Fatal(err)
		}
		w.Write(b)
	}))

	c := NewClient(
		WithBaseURL(ts.URL),
	)

	want := &ReactionDefinition{
		Name:               "payment-processed-email-reaction",
		Feed:               "payment",
		ReactOnEventType:   "PaymentProcessed",
		CancelOnEventTypes: []string{"OrderCanceledEvent"},
		TriggerTimeField:   "my.event.data.field",
		Offset:             "PT1H",
		Action: &Action{
			ActionType: ActionTypeHTTPPost,
		},
	}

	got, err := c.ReactionDefinition(context.Background(), want.Name)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got = %v; want = %v", got, want)
	}
}
