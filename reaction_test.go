package serialized

import (
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

	r, err := c.ListReactions()
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

		w.WriteHeader(http.StatusCreated)
	}))

	c := NewClient(
		WithBaseURL(ts.URL),
	)

	r := Reaction{
		ID:        "be278b27-8687-42b4-a502-164a6702797c",
		Name:      "PaymentProcessedEmailReaction",
		Feed:      "payment",
		EventType: "PaymentProcessed",
		Delay:     "PT1H",
		Action: Action{
			HTTPMethod: "POST",
			TargetURI:  "https://your-webhook",
			Body:       "A new payment was processed",
			ActionType: "HTTP",
		},
	}

	if err := c.CreateReaction(r); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteReaction(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))

	c := NewClient(
		WithBaseURL(ts.URL),
	)

	if err := c.DeleteReaction("be278b27-8687-42b4-a502-164a6702797c"); err != nil {
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

	want := Reaction{
		ID:        "be278b27-8687-42b4-a502-164a6702797c",
		Name:      "PaymentProcessedEmailReaction",
		Feed:      "payment",
		EventType: "PaymentProcessed",
		Delay:     "PT1H",
		Action: Action{
			HTTPMethod: "POST",
			TargetURI:  "https://your-webhook",
			Body:       "A new payment was processed",
			ActionType: "HTTP",
		},
	}

	got, err := c.Reaction(want.ID)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got = %v; want = %v", got, want)
	}
}
