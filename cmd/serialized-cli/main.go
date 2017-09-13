package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cobra"

	serialized "github.com/marcusolsson/serialized-go"
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

	var since int64
	var current bool
	var maxNumEvents int

	var (
		eventType       string
		eventID         string
		eventData       string
		expectedVersion int64
	)

	var cmdStore = &cobra.Command{
		Use:   "store [type] [id]",
		Short: "Store a new event",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if eventID == "" {
				eventID = uuid.NewV4().String()
			}
			if eventType == "" {
				fmt.Println("event type was not specified")
				os.Exit(1)
			}
			if eventData == "" {
				fmt.Println("event data was empty")
				os.Exit(1)
			}

			event := serialized.Event{
				Type: eventType,
				ID:   eventID,
				Data: []byte(eventData),
			}

			if err := client.Store(args[0], args[1], expectedVersion, event); err != nil {
				fmt.Println("unable to store event:", err)
				os.Exit(1)
			}
		},
	}

	var cmdAggregate = &cobra.Command{
		Use:   "aggregate [type] [id]",
		Short: "Display an aggregate",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			agg, err := client.LoadAggregate(args[0], args[1])
			if err != nil {
				fmt.Println("unable to load aggregate:", err)
				os.Exit(1)
			}

			w := tabwriter.NewWriter(os.Stdout, 5, 4, 1, ' ', 0)
			fmt.Fprintln(w, "TYPE:", "\t", agg.Type)
			fmt.Fprintln(w, "ID:", "\t", agg.ID)
			fmt.Fprintln(w, "VERSION:", "\t", agg.Version)
			fmt.Fprintln(w)
			fmt.Fprintf(w, "Showing the %d most recent events:\n", maxNumEvents)
			fmt.Fprintln(w)

			w.Flush()

			fmt.Fprintln(w, "ID:", "\t", "Type:", "\t", "Data:")

			events := agg.Events
			if len(events) > maxNumEvents {
				events = events[len(events)-maxNumEvents:]
			}
			for _, e := range events {
				fmt.Fprintln(w, e.ID, "\t", e.Type, "\t", string(e.Data))
			}
			w.Flush()
		},
	}

	var cmdFeed = &cobra.Command{
		Use:   "feed [name]",
		Short: "Display the feed",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if current {
				seq, err := client.FeedSequenceNumber(args[0])
				if err != nil {
					fmt.Println("unable to get sequence number:", err)
					os.Exit(1)
				}
				fmt.Println(seq)
				return
			}

			feed, err := client.Feed(args[0], since)
			if err != nil {
				fmt.Println("unable to get feed:", err)
				os.Exit(1)
			}

			w := tabwriter.NewWriter(os.Stdout, 5, 5, 2, ' ', 0)
			fmt.Fprintln(w, "SEQUENCE", "\t", "AGGREGATE", "\t", "EVENTS")

			for _, e := range feed.Entries {
				fmt.Fprintf(w, "%d\t%s\t%s\n", e.SequenceNumber, e.AggregateID, e.Events[0].Data)

				if len(e.Events) > 1 {
					for _, ev := range e.Events[1:] {
						fmt.Fprintf(w, "\t\t%s\n", string(ev.Data))
					}
				}
			}
			w.Flush()
		},
	}

	var cmdFeeds = &cobra.Command{
		Use:   "feeds",
		Short: "List all existing feeds",
		Run: func(cmd *cobra.Command, args []string) {
			feeds, err := client.Feeds()
			if err != nil {
				fmt.Println("unable to list feeds:", err)
				os.Exit(1)
			}
			for _, f := range feeds {
				fmt.Println(f)
			}
		},
	}

	cmdStore.Flags().StringVarP(&eventID, "id", "i", "", "Optional event ID.")
	cmdStore.Flags().StringVarP(&eventType, "type", "t", "", "Event type")
	cmdStore.Flags().StringVarP(&eventData, "data", "d", "", "Event data")
	cmdStore.Flags().Int64VarP(&expectedVersion, "expected-version", "v", 0, "Version number for optimistic concurrency control.")

	cmdAggregate.Flags().IntVarP(&maxNumEvents, "max-events", "m", 10, "Maximum number of events to show in preview.")

	cmdFeed.Flags().Int64VarP(&since, "since", "s", 0, "Sequence number to start from.")
	cmdFeed.Flags().BoolVarP(&current, "current", "c", false, "Return current sequence number at head for a given feed.")

	var rootCmd = &cobra.Command{
		Use:   "serialized-cli",
		Short: "Interact with the Serialized.io API from the command-line.",
	}
	rootCmd.AddCommand(cmdStore, cmdAggregate, cmdFeed, cmdFeeds)
	rootCmd.Execute()
}
