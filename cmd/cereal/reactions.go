package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	serialized "github.com/marcusolsson/serialized-go"
)

func reactionsDefinitionsGetHandler(c *serialized.Client, name string) error {
	def, err := c.ReactionDefinition(context.Background(), name)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)

	var level int

	fmt.Fprintf(w, kv("Name", def.Name, level))
	fmt.Fprintf(w, kv("Feed", def.Feed, level))
	fmt.Fprintf(w, kv("Reacts on", def.ReactOnEventType, level))
	fmt.Fprintf(w, kv("Cancels on", strings.Join(def.CancelOnEventTypes, ", "), level))
	fmt.Fprintf(w, kv("Trigger time field", def.TriggerTimeField, level))
	fmt.Fprintf(w, kv("Offset", def.Offset, level))
	fmt.Fprintf(w, kv("Action", "", level))

	w.Flush()

	level++

	fmt.Fprintf(w, kv("Type", string(def.Action.ActionType), level))
	fmt.Fprintf(w, kv("Target URI", string(def.Action.TargetURI), level))
	fmt.Fprintf(w, kv("Body", string(def.Action.Body), level))

	w.Flush()

	return nil
}

func reactionsDefinitionsDeleteHandler(c *serialized.Client, name string) error {
	if err := c.DeleteReactionDefinition(context.Background(), name); err != nil {
		return err
	}

	fmt.Printf("reaction definition %q deleted\n", name)

	return nil
}

func reactionsDefinitionsListHandler(c *serialized.Client) error {
	defs, err := c.ListReactionDefinitions(context.Background())
	if err != nil {
		return err
	}

	if len(defs) == 0 {
		fmt.Println("No reactions found.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)

	fmt.Fprintln(w, strings.Join([]string{"NAME", "FEED", "REACTS ON", "ACTION", "TARGET URI"}, "\t"))
	for _, d := range defs {
		fmt.Fprintln(w, strings.Join([]string{d.Name, d.Feed, d.ReactOnEventType, string(d.Action.ActionType), d.Action.TargetURI}, "\t"))
	}

	return w.Flush()
}
