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

func projectionsSingleGetHandler(c *serialized.Client, projName, aggID string) error {
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

func projectionsAggregatedGetHandler(c *serialized.Client, projName string) error {
	proj, err := c.AggregatedProjection(context.Background(), projName)
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

func projectionsAggregatedListHandler(c *serialized.Client) error {
	projs, err := c.ListAggregatedProjections(context.Background())
	if err != nil {
		return err
	}

	if len(projs) == 0 {
		fmt.Println("No aggregated projections found.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)

	fmt.Fprintln(w, strings.Join([]string{"ID", "DATA"}, "\t"))
	for _, p := range projs {
		fmt.Fprintln(w, strings.Join([]string{p.ID, string(p.Data)}, "\t"))
	}

	return w.Flush()
}

func projectionsDefinitionsGetHandler(c *serialized.Client, name string) error {
	proj, err := c.ProjectionDefinition(context.Background(), name)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)

	var level int

	fmt.Fprintf(w, kv("Name", proj.Name, level))
	fmt.Fprintf(w, kv("Feed", proj.Feed, level))
	fmt.Fprintf(w, kv("Handlers", "", level))

	w.Flush()

	for i, h := range proj.Handlers {
		level++
		fmt.Fprintf(w, kv(h.EventType, "", level))

		w.Flush()

		for _, f := range h.Functions {
			level++
			fmt.Fprintf(w, kv("Function", f.Function, level))
			fmt.Fprintf(w, kv("Target selector", f.TargetSelector, level))
			fmt.Fprintf(w, kv("Event selector", f.EventSelector, level))
			fmt.Fprintf(w, kv("Target filter", f.TargetFilter, level))
			fmt.Fprintf(w, kv("Event filter", f.EventFilter, level))
			fmt.Fprintf(w, kv("Raw data", fmt.Sprintf("%s", f.RawData), level))

			w.Flush()

			if i < len(proj.Handlers)-1 {
				fmt.Println()
			}
		}
	}

	return nil
}

func kv(key, val string, indent int) string {
	return fmt.Sprintf("%s%s:\t%s\n", strings.Repeat("  ", indent), key, val)
}

func projectionsDefinitionsDeleteHandler(c *serialized.Client, name string) error {
	if err := c.DeleteProjectionDefinition(context.Background(), name); err != nil {
		return err
	}

	fmt.Printf("projection definition %q deleted\n", name)

	return nil
}

func projectionsDefinitionsListHandler(c *serialized.Client) error {
	defs, err := c.ListProjectionDefinitions(context.Background())
	if err != nil {
		return err
	}

	if len(defs) == 0 {
		fmt.Println("No projection definitions found.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)

	fmt.Fprintln(w, strings.Join([]string{"NAME", "FEED"}, "\t"))
	for _, d := range defs {
		fmt.Fprintln(w, strings.Join([]string{d.Name, d.Feed}, "\t"))
	}

	return w.Flush()
}
