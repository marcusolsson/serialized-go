package serialized

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProjectionListDefinitions(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := loadJSON("testdata/projection_list_definitions_response.json")
		if err != nil {
			t.Fatal(err)
		}
		w.Write(b)
	}))

	c := NewClient(
		WithBaseURL(ts.URL),
	)

	defs, err := c.ListProjectionDefinitions(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if len(defs) != 1 {
		t.Fatalf("unexpected number of definitions = %d; want = %d", len(defs), 1)
	}
	if want := "orders"; defs[0].Name != want {
		t.Fatalf("want = %s; got = %s", want, defs[0].Name)
	}
	if want := "order"; defs[0].Feed != want {
		t.Fatalf("want = %s; got = %s", want, defs[0].Feed)
	}
}

func TestProjectionCreateDefinition(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		want, err := loadJSON("testdata/projection_create_definition_request.json")
		if err != nil {
			t.Fatal(err)
		}
		assertEqualJSON(t, got, want)
	}))

	c := NewClient(
		WithBaseURL(ts.URL),
	)

	def := &ProjectionDefinition{
		Name: "orders",
		Feed: "order",
		Handlers: []*EventHandler{
			{
				EventType: "OrderCancelledEvent",
				Functions: []*Function{
					{
						Function:       "inc",
						TargetSelector: "$.projection.orders[?]",
						EventSelector:  "$.event[?]",
						TargetFilter:   "@.orderId == $.event.orderId",
						EventFilter:    "@.orderAmount > 4000",
					},
				},
			},
		},
	}

	if err := c.CreateProjectionDefinition(context.Background(), def); err != nil {
		t.Fatal(err)
	}
}

func TestProjectionGetSingle(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := loadJSON("testdata/projection_get_single_response.json")
		if err != nil {
			t.Fatal(err)
		}
		w.Write(b)
	}))

	c := NewClient(
		WithBaseURL(ts.URL),
	)

	proj, err := c.Projection(context.Background(), "orders", "foo")
	if err != nil {
		t.Fatal(err)
	}

	if want := "string (uuid)"; want != proj.ID {
		t.Fatalf("want = %s; got = %s", want, proj.ID)
	}
	if want := []byte(`{"field":"data"}`); !bytes.Equal(want, proj.Data) {
		t.Fatalf("want = %s; got = %s", want, string(proj.Data))
	}
}
