package serialized

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestProjectionListDefinitions(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := loadJSON("testdata/projection_list_definitions_response.json")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write(b); err != nil {
			t.Fatal(err)
		}
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

func TestProjectionGetDefinition(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := loadJSON("testdata/projection_get_definition_response.json")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write(b); err != nil {
			t.Fatal(err)
		}
	}))

	c := NewClient(
		WithBaseURL(ts.URL),
	)

	got, err := c.ProjectionDefinition(context.Background(), "orders")
	if err != nil {
		t.Fatal(err)
	}

	want := &ProjectionDefinition{
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

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got = %v; want = %v", got, want)
	}
}

func TestProjectionDeleteDefinition(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/projections/definitions/foo" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))

	c := NewClient(
		WithBaseURL(ts.URL),
	)

	if err := c.DeleteProjectionDefinition(context.Background(), "foo"); err != nil {
		t.Fatal(err)
	}
}

func TestProjectionGetSingle(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := loadJSON("testdata/projection_get_single_response.json")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write(b); err != nil {
			t.Fatal(err)
		}
	}))

	c := NewClient(
		WithBaseURL(ts.URL),
	)

	proj, err := c.SingleProjection(context.Background(), "orders", "foo")
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

func TestProjectionListSingle(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := loadJSON("testdata/projection_list_single_response.json")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write(b); err != nil {
			t.Fatal(err)
		}
	}))

	c := NewClient(
		WithBaseURL(ts.URL),
	)

	projs, err := c.ListSingleProjections(context.Background(), "foo")
	if err != nil {
		t.Fatal(err)
	}

	if len(projs) != 1 {
		t.Fatalf("unexpected number of definitions = %d; want = %d", len(projs), 1)
	}
	if projs[0].ID != "string (uuid)" {
		t.Fatalf("unexpected projection id: %s", projs[0].ID)
	}
}

func TestProjectionGetAggregated(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := loadJSON("testdata/projection_get_agg_response.json")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write(b); err != nil {
			t.Fatal(err)
		}
	}))

	c := NewClient(
		WithBaseURL(ts.URL),
	)

	proj, err := c.AggregatedProjection(context.Background(), "foo")
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

func TestProjectionListAggregated(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := loadJSON("testdata/projection_list_agg_response.json")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write(b); err != nil {
			t.Fatal(err)
		}
	}))

	c := NewClient(
		WithBaseURL(ts.URL),
	)

	projs, err := c.ListAggregatedProjections(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if len(projs) != 1 {
		t.Fatalf("unexpected number of definitions = %d; want = %d", len(projs), 1)
	}
	if projs[0].ID != "string (uuid)" {
		t.Fatalf("unexpected projection id: %s", projs[0].ID)
	}
}
