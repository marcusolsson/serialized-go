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

func main() {
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
		kingpin.FatalIfError(
			storeEvent(client, *storeAggType, *storeAggID, *storeEventType, *storeEventID, *storeData, *storeExpectedVersion),
			"unable to store event")
	case aggregate.FullCommand():
		kingpin.FatalIfError(
			showAggregate(client, *aggregateType, *aggregateID, *aggregateLimit),
			"unable to get aggregate")
	case projection.FullCommand():
		kingpin.FatalIfError(
			showProjection(client, *projectionName, *projectionAggID),
			"unable to get projection")
	case feed.FullCommand():
		kingpin.FatalIfError(
			showFeed(client, *feedName, *feedSince, *feedCurrent),
			"unable to get feed")
	case feeds.FullCommand():
		kingpin.FatalIfError(
			listFeeds(client),
			"unable to list feeds")
	}
}

func storeEvent(c *serialized.Client, aggType, aggID, eventType, eventID, data string, version int64) error {
	if eventID == "" {
		id, err := uuid.NewV4()
		if err != nil {
			return err

		}
		eventID = id.String()
	}

	if aggID == "" {
		id, err := uuid.NewV4()
		if err != nil {
			return err
		}
		aggID = id.String()
	}

	event := &serialized.Event{
		Type: eventType,
		ID:   eventID,
		Data: []byte(data),
	}

	return c.Store(context.Background(), aggType, aggID, version, event)
}

func showAggregate(c *serialized.Client, aggType, aggID string, limit int) error {
	agg, err := c.LoadAggregate(context.Background(), aggType, aggID)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 5, 4, 1, ' ', 0)
	fmt.Fprintln(w, "Type:", "\t", agg.Type)
	fmt.Fprintln(w, "ID:", "\t", agg.ID)
	fmt.Fprintln(w, "Version:", "\t", agg.Version)
	fmt.Fprintln(w)
	fmt.Fprintf(w, "Showing the %d most recent events:\n", limit)
	fmt.Fprintln(w)

	w.Flush()

	fmt.Fprintln(w, "EVENT ID", "\t", "TYPE", "\t", "DATA")

	events := agg.Events
	if len(events) > limit {
		events = events[len(events)-limit:]
	}
	for _, e := range events {
		var buf bytes.Buffer
		if err := json.Compact(&buf, e.Data); err != nil {
			return err
		}
		fmt.Fprintln(w, e.ID, "\t", e.Type, "\t", buf.String())
	}
	w.Flush()

	return nil
}

func showProjection(c *serialized.Client, projName, aggID string) error {
	proj, err := c.SingleProjection(context.Background(), projName, aggID)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := json.Indent(&buf, proj.Data, "", "  "); err != nil {
		return err
	}

	fmt.Println(buf.String())

	return nil
}

func showFeed(c *serialized.Client, feed string, since int64, showCurrent bool) error {
	ctx := context.Background()

	if showCurrent {
		seq, err := c.FeedSequenceNumber(ctx, feed)
		if err != nil {
			return fmt.Errorf("unable to get sequence number: %s", err)
		}

		fmt.Println(seq)

		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)

	fmt.Fprintln(w, strings.Join([]string{"TIMESTAMP", "AGGREGATE ID", "EVENT TYPE"}, "\t"))

	return c.Feed(ctx, feed, since, func(e *serialized.FeedEntry) {
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
}

func listFeeds(c *serialized.Client) error {
	feeds, err := c.Feeds(context.Background())
	if err != nil {
		return err
	}

	for _, f := range feeds {
		fmt.Println(f)
	}

	return nil
}
