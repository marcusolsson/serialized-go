package main

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"
	"time"

	serialized "github.com/marcusolsson/serialized-go"
	"github.com/spf13/cobra"
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

	var since int
	var current bool

	var cmdAggregate = &cobra.Command{
		Use:   "aggregate [type] [id]",
		Short: "Display an aggregate",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {

			agg, err := client.LoadAggregate(args[0], args[1])
			if err != nil {
				log.Fatal(err)
			}

			w := tabwriter.NewWriter(os.Stdout, 5, 4, 1, ' ', 0)
			fmt.Fprintln(w, "TYPE:", "\t", agg.Type)
			fmt.Fprintln(w, "ID:", "\t", agg.ID)
			fmt.Fprintln(w, "VERSION:", "\t", agg.Version)
			fmt.Fprintln(w)
			fmt.Fprintln(w, "Showing the 10 most recent events:")
			fmt.Fprintln(w)

			w.Flush()

			fmt.Fprintln(w, "ID:", "\t", "Type:", "\t", "Data:")

			events := agg.Events
			if len(events) > 5 {
				events = events[:5]
			}

			for _, e := range events {
				fmt.Fprintln(w, e.ID, "\t", e.Type, "\t", string(e.Data))
			}

			w.Flush()
		},
	}

	var cmdFeed = &cobra.Command{
		Use:   "feed [name]",
		Short: "Display the feed output",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if current {
				seq, err := client.FeedSequenceNumber(args[0])
				if err != nil {
					log.Fatal(err)
				}

				fmt.Println(seq)

				return
			}

			feed, err := client.Feed(args[0], since)
			if err != nil {
				log.Fatal(err)
			}

			w := tabwriter.NewWriter(os.Stdout, 5, 4, 1, ' ', 0)
			fmt.Fprintln(w, "SEQUENCE\tAGGREGATE ID\tNUM EVENTS\tTIMESTAMP")

			for _, e := range feed.Entries {
				t := time.Unix(e.Timestamp/1000, 0)
				fmt.Fprintf(w, "%d\t%s\t%d\t%s\n", e.SequenceNumber, e.AggregateID, len(e.Events), t.Format(time.RFC3339))
			}
			w.Flush()
		},
	}

	var cmdFeeds = &cobra.Command{
		Use:   "feeds",
		Short: "List the available fields",
		Run: func(cmd *cobra.Command, args []string) {
			feeds, err := client.Feeds()
			if err != nil {
				log.Fatal(err)
			}
			for _, f := range feeds {
				fmt.Println(f)
			}
		},
	}

	cmdFeed.Flags().IntVarP(&since, "since", "s", 0, "Optional sequence number to start from.")
	cmdFeed.Flags().BoolVarP(&current, "current", "c", false, "Return current sequence number at head for a given feed.")

	var rootCmd = &cobra.Command{Use: "serialized-cli"}
	rootCmd.AddCommand(cmdAggregate, cmdFeed, cmdFeeds)
	rootCmd.Execute()
}
