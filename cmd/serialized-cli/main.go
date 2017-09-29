package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/alecthomas/kingpin"
	serialized "github.com/marcusolsson/serialized-go"
	uuid "github.com/satori/go.uuid"
)

var (
	app = kingpin.New("serialized-cli", "Interact with the Serialized.io API from the command-line.").Version("0.1.0")

	store                = app.Command("store", "Store a new event.")
	storeType            = store.Flag("type", "Type of event.").Short('t').Required().String()
	storeID              = store.Flag("id", "ID of event.").String()
	storeData            = store.Flag("data", "Event data.").Short('d').Required().String()
	storeExpectedVersion = store.Flag("expected-version", "Version number for optimistic concurrency control.").Short('v').Int64()

	aggregate      = app.Command("aggregate", "Display an aggregate.")
	aggregateID    = aggregate.Arg("id", "ID of aggregate.").Required().String()
	aggregateType  = aggregate.Flag("type", "Type of aggregate.").Short('t').String()
	aggregateLimit = aggregate.Flag("limit", "Max number of events to show in preview.").Short('l').Default("10").Int()

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
		id := *storeID
		if id == "" {
			id = uuid.NewV4().String()
		}

		event := &serialized.Event{
			Type: *storeType,
			ID:   id,
			Data: []byte(*storeData),
		}

		err := client.Store(context.Background(), *storeType, id, *storeExpectedVersion, event)
		if err != nil {
			kingpin.Fatalf("unable to store event: %s", err)
		}
	case aggregate.FullCommand():
		agg, err := client.LoadAggregate(context.Background(), *aggregateType, *aggregateID)
		if err != nil {
			kingpin.Fatalf("unable to load aggregate: %s", err)
		}

		w := tabwriter.NewWriter(os.Stdout, 5, 4, 1, ' ', 0)
		fmt.Fprintln(w, "TYPE:", "\t", agg.Type)
		fmt.Fprintln(w, "ID:", "\t", agg.ID)
		fmt.Fprintln(w, "VERSION:", "\t", agg.Version)
		fmt.Fprintln(w)
		fmt.Fprintf(w, "Showing the %d most recent events:\n", *aggregateLimit)
		fmt.Fprintln(w)

		w.Flush()

		fmt.Fprintln(w, "ID:", "\t", "Type:", "\t", "Data:")

		events := agg.Events
		if len(events) > *aggregateLimit {
			events = events[len(events)-*aggregateLimit:]
		}
		for _, e := range events {
			fmt.Fprintln(w, e.ID, "\t", e.Type, "\t", string(e.Data))
		}
		w.Flush()

	case feed.FullCommand():
		ctx := context.Background()

		if *feedCurrent {
			seq, err := client.FeedSequenceNumber(ctx, *feedName)
			if err != nil {
				kingpin.Fatalf("unable to get sequence number: %s", err)
			}
			fmt.Println(seq)
			return
		}

		err := client.Feed(ctx, *feedName, *feedSince, func(e *serialized.FeedEntry) {
			for _, ev := range e.Events {
				var buf bytes.Buffer
				if err := json.Compact(&buf, ev.Data); err != nil {
					kingpin.Fatalf("unable to format event data: %s", err)
				}
				fmt.Println(buf.String())
			}
		})
		if err != nil {
			kingpin.Fatalf("unable to get feed: %s", err)
		}
	case feeds.FullCommand():
		feeds, err := client.Feeds(context.Background())
		if err != nil {
			kingpin.Fatalf("unable to get feeds: %s", err)
		}
		for _, f := range feeds {
			fmt.Println(f)
		}
	}
}
