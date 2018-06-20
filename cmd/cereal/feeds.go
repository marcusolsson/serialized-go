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
)

func feedsGetHandler(c *serialized.Client, feed string, since int64, showCurrent bool) error {
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

func feedsListHandler(c *serialized.Client) error {
	feeds, err := c.Feeds(context.Background())
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)

	fmt.Fprintln(w, strings.Join([]string{"TYPE", "AGGREGATES", "BATCHES", "EVENTS"}, "\t"))

	for _, f := range feeds {
		fmt.Fprintln(w, fmt.Sprintf("%s\t%d\t%d\t%d",
			f.AggregateType, f.AggregateCount, f.BatchCount, f.EventCount))
	}

	return w.Flush()
}
