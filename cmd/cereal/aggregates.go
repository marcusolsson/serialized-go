package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	serialized "github.com/marcusolsson/serialized-go"
)

func aggregatesGetHandler(c *serialized.Client, aggType, aggID string, limit int) error {
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

func aggregatesDeleteHandler(c *serialized.Client, aggType string) error {
	ctx := context.Background()

	token, err := c.RequestDeleteAggregateByType(ctx, aggType)
	if err != nil {
		return err
	}

	fmt.Printf("WARNING: This will permanently delete all aggregates of type %q, including all associated events.\n", aggType)
	fmt.Printf("Are you sure you want to delete all aggregates of type %q? (yes/no): ", aggType)

	var answer string
	_, err = fmt.Scan(&answer)
	if err != nil {
		return err
	}

	if yesOrNo(answer) {
		if err := c.DeleteAggregateByType(ctx, aggType, token); err != nil {
			return err
		}
		fmt.Println("Successfully deleted aggregates.")
	} else {
		fmt.Println("Canceled. No changes were made.")
	}

	return nil
}

func yesOrNo(str string) bool {
	return strings.ToLower(str) == "yes"
}
