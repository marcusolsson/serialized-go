package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/alecthomas/kingpin"
	serialized "github.com/marcusolsson/serialized-go"
	uuid "github.com/satori/go.uuid"
)

var (
	app = kingpin.New("serialized-cli", "Interact with the Serialized.io API from the command-line.").Version("0.1.0")

	store                = app.Command("store", "Store a new event.")
	storeAggType         = store.Flag("agg-type", "Type of aggregate.").Required().String()
	storeAggID           = store.Flag("agg-id", "ID of aggregate.").String()
	storeEventType       = store.Flag("event-type", "Type of event.").Required().String()
	storeEventID         = store.Flag("event-id", "ID of event.").String()
	storeData            = store.Flag("data", "Event data.").Short('d').Required().String()
	storeExpectedVersion = store.Flag("expected-version", "Version number for optimistic concurrency control.").Int64()

	aggregate      = app.Command("aggregate", "Display an aggregate.")
	aggregateID    = aggregate.Arg("id", "ID of aggregate.").Required().String()
	aggregateType  = aggregate.Flag("type", "Type of aggregate.").Short('t').Required().String()
	aggregateLimit = aggregate.Flag("limit", "Max number of events to show in preview.").Short('l').Default("10").Int()

	projection      = app.Command("projection", "Display a projection.")
	projectionName  = projection.Arg("name", "Name of the projection.").Required().String()
	projectionAggID = projection.Flag("agg-id", "ID of aggregate.").Required().String()

	feed        = app.Command("feed", "Display the feed.")
	feedName    = feed.Arg("name", "Name of feed.").Required().String()
	feedSince   = feed.Flag("since", "Sequence number to start from.").Short('s').Int64()
	feedCurrent = feed.Flag("current", "Return current sequence number at head for a given feed.").Short('c').Bool()

	feeds = app.Command("feeds", "List all existing feeds.")
)

func main() {
	var (
		accessKey       = os.Getenv("SERIALIZED_ACCESS_KEY")
		secretAccessKey = os.Getenv("SERIALIZED_SECRET_ACCESS_KEY")
	)

	client := serialized.NewClient(
		serialized.WithAccessKey(accessKey),
		serialized.WithSecretAccessKey(secretAccessKey),
	)

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case store.FullCommand():
		eventID := *storeEventID
		if eventID == "" {
			id, err := uuid.NewV4()
			if err != nil {
				kingpin.Fatalf("Unable to create uuid: %s", err)
			}
			eventID = id.String()
		}

		aggID := *storeAggID
		if aggID == "" {
			id, err := uuid.NewV4()
			if err != nil {
				kingpin.Fatalf("Unable to create uuid: %s", err)
			}
			aggID = id.String()
		}

		event := &serialized.Event{
			Type: *storeEventType,
			ID:   eventID,
			Data: []byte(*storeData),
		}

		err := client.Store(context.Background(), *storeAggType, aggID, *storeExpectedVersion, event)
		if err != nil {
			kingpin.Fatalf("Unable to store event: %s", err)
		}

	case aggregate.FullCommand():
		agg, err := client.LoadAggregate(context.Background(), *aggregateType, *aggregateID)
		if err != nil {
			kingpin.Fatalf("Unable to load aggregate: %s", err)
		}

		w := tabwriter.NewWriter(os.Stdout, 5, 4, 1, ' ', 0)
		fmt.Fprintln(w, "Type:", "\t", agg.Type)
		fmt.Fprintln(w, "ID:", "\t", agg.ID)
		fmt.Fprintln(w, "Version:", "\t", agg.Version)
		fmt.Fprintln(w)
		fmt.Fprintf(w, "Showing the %d most recent events:\n", *aggregateLimit)
		fmt.Fprintln(w)

		w.Flush()

		fmt.Fprintln(w, "ID", "\t", "TYPE", "\t", "DATA")

		events := agg.Events
		if len(events) > *aggregateLimit {
			events = events[len(events)-*aggregateLimit:]
		}
		for _, e := range events {
			var buf bytes.Buffer
			if err := json.Compact(&buf, e.Data); err != nil {
				kingpin.Fatalf("Unable to show aggregate data: %s", err)
			}
			fmt.Fprintln(w, e.ID, "\t", e.Type, "\t", buf.String())
		}
		w.Flush()

	case projection.FullCommand():
		ctx := context.Background()

		proj, err := client.Projection(ctx, *projectionName, *projectionAggID)
		if err != nil {
			kingpin.Fatalf("Unable to load projection: %s", err)
		}

		var buf bytes.Buffer
		if err := json.Indent(&buf, proj.Data, "", "  "); err != nil {
			kingpin.Fatalf("Unable to show projection data: %s", err)
		}
		fmt.Println(buf.String())

	case feed.FullCommand():
		ctx := context.Background()

		if *feedCurrent {
			seq, err := client.FeedSequenceNumber(ctx, *feedName)
			if err != nil {
				kingpin.Fatalf("Unable to get sequence number: %s", err)
			}
			fmt.Println(seq)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)
		fmt.Fprintln(w, strings.Join([]string{"TIMESTAMP", "AGGREGATE ID", "EVENT TYPE"}, "\t"))

		err := client.Feed(ctx, *feedName, *feedSince, func(e *serialized.FeedEntry) {
			ts := time.Unix(e.Timestamp/1000, 0)
			for _, ev := range e.Events {
				var buf bytes.Buffer
				if err := json.Compact(&buf, ev.Data); err != nil {
					kingpin.Fatalf("unable to format event data: %s", err)
				}
				fmt.Fprintln(w, strings.Join([]string{ts.Format(time.RFC1123Z), e.AggregateID, ev.Type}, "\t"))
				w.Flush()
			}
		})
		if err != nil {
			kingpin.Fatalf("Unable to get feed: %s", err)
		}

	case feeds.FullCommand():
		feeds, err := client.Feeds(context.Background())
		if err != nil {
			kingpin.Fatalf("Unable to get feeds: %s", err)
		}
		for _, f := range feeds {
			fmt.Println(f)
		}
	}
}
