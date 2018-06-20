package main

import (
	"os"

	"github.com/alecthomas/kingpin"
	serialized "github.com/marcusolsson/serialized-go"
)

func main() {
	var (
		app = kingpin.New("serialized-cli", "Interact with the Serialized.io API from the command-line.").Version("0.1.0")

		events = app.Command("events", "Event commands.")

		eventsStore                = events.Command("store", "Store a new event.")
		eventsStoreAggType         = eventsStore.Flag("agg-type", "Type of aggregate.").Short('a').Required().String()
		eventsStoreAggID           = eventsStore.Flag("agg-id", "ID of aggregate.").String()
		eventsStoreEventType       = eventsStore.Flag("event-type", "Type of event.").Short('e').Required().String()
		eventsStoreEventID         = eventsStore.Flag("event-id", "ID of event.").String()
		eventsStoreData            = eventsStore.Flag("data", "Event data.").Short('d').Required().String()
		eventsStoreExpectedVersion = eventsStore.Flag("expected-version", "Version number for optimistic concurrency control.").Int64()

		aggregates = app.Command("aggregates", "Aggregate commands.")

		aggregatesGet        = aggregates.Command("get", "Show aggregate")
		aggregatesGetID      = aggregatesGet.Arg("id", "ID of aggregate.").Required().String()
		aggregatesGetType    = aggregatesGet.Flag("type", "Type of aggregate.").Short('t').Required().String()
		aggregatesGetLimit   = aggregatesGet.Flag("limit", "Max number of events to show in preview.").Short('l').Default("10").Int()
		aggregatesDelete     = aggregates.Command("delete", "Delete aggregates of a given type")
		aggregatesDeleteType = aggregatesDelete.Flag("type", "Type of aggregate.").Short('t').Required().String()

		feeds = app.Command("feeds", "Feed commands.")

		feedsGet        = feeds.Command("get", "Show feed.")
		feedsGetName    = feedsGet.Arg("name", "Name of feed.").Required().String()
		feedsGetSince   = feedsGet.Flag("since", "Sequence number to start from.").Short('s').Int64()
		feedsGetCurrent = feedsGet.Flag("current", "Return current sequence number at head for a given feed.").Short('c').Bool()
		feedsList       = feeds.Command("list", "List all existing feeds.")

		projections = app.Command("projections", "Projection commands.")

		projectionsSingle               = projections.Command("single", "Single projection commands.")
		projectionsSingleGet            = projectionsSingle.Command("get", "Show projection.")
		projectionsSingleGetName        = projectionsSingleGet.Arg("name", "Name of the projection.").Required().String()
		projectionsSingleGetAggregateID = projectionsSingleGet.Flag("agg-id", "ID of aggregate.").Required().String()

		projectionsAggregated        = projections.Command("aggregated", "Aggregated projection commands.")
		projectionsAggregatedGet     = projectionsAggregated.Command("get", "Show aggregated projection.")
		projectionsAggregatedGetName = projectionsAggregatedGet.Arg("name", "Name of the aggregated projection.").Required().String()
		projectionsAggregatedList    = projectionsAggregated.Command("list", "List aggregated projections.")

		projectionsDefinitions           = projections.Command("definitions", "Projection definitions commands.")
		projectionsDefinitionsGet        = projectionsDefinitions.Command("get", "Show projection definition.")
		projectionsDefinitionsGetName    = projectionsDefinitionsGet.Arg("name", "Name of the projection definition.").Required().String()
		projectionsDefinitionsDelete     = projectionsDefinitions.Command("delete", "Delete a projection definition.")
		projectionsDefinitionsDeleteName = projectionsDefinitionsDelete.Arg("name", "Name of the projection definition.").Required().String()
		projectionsDefinitionsList       = projectionsDefinitions.Command("list", "List projection definitions.")

		reactions = app.Command("reactions", "Reaction commands.")

		reactionsDefinitions           = reactions.Command("definitions", "Reaction commands.")
		reactionsDefinitionsGet        = reactionsDefinitions.Command("get", "Show reaction definition.")
		reactionsDefinitionsGetName    = reactionsDefinitionsGet.Arg("name", "Name of the reaction definition").Required().String()
		reactionsDefinitionsDelete     = reactionsDefinitions.Command("delete", "Delete a reaction definition.")
		reactionsDefinitionsDeleteName = reactionsDefinitionsDelete.Arg("name", "Name of the reaction definition.").Required().String()
		reactionsDefinitionsList       = reactionsDefinitions.Command("list", "List reaction definitions.")
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
	// Events
	case eventsStore.FullCommand():
		kingpin.FatalIfError(
			eventsStoreHandler(client, *eventsStoreAggType, *eventsStoreAggID, *eventsStoreEventType, *eventsStoreEventID, *eventsStoreData, *eventsStoreExpectedVersion),
			"unable to store event")

		// Aggregates
	case aggregatesGet.FullCommand():
		kingpin.FatalIfError(
			aggregatesGetHandler(client, *aggregatesGetType, *aggregatesGetID, *aggregatesGetLimit),
			"unable to get aggregate")
	case aggregatesDelete.FullCommand():
		kingpin.FatalIfError(
			aggregatesDeleteHandler(client, *aggregatesDeleteType),
			"unable to delete aggregate")

		// Projections
	case projectionsSingleGet.FullCommand():
		kingpin.FatalIfError(
			projectionsSingleGetHandler(client, *projectionsSingleGetName, *projectionsSingleGetAggregateID),
			"unable to get single projection")
	case projectionsAggregatedGet.FullCommand():
		kingpin.FatalIfError(
			projectionsAggregatedGetHandler(client, *projectionsAggregatedGetName),
			"unable to get aggregated projection")
	case projectionsAggregatedList.FullCommand():
		kingpin.FatalIfError(
			projectionsAggregatedListHandler(client),
			"unable to list aggregated projections")
	case projectionsDefinitionsGet.FullCommand():
		kingpin.FatalIfError(
			projectionsDefinitionsGetHandler(client, *projectionsDefinitionsGetName),
			"unable to get projection definition")
	case projectionsDefinitionsDelete.FullCommand():
		kingpin.FatalIfError(
			projectionsDefinitionsDeleteHandler(client, *projectionsDefinitionsDeleteName),
			"unable to delete projection definition")
	case projectionsDefinitionsList.FullCommand():
		kingpin.FatalIfError(
			projectionsDefinitionsListHandler(client),
			"unable to list projection definitions")

		// Feeds
	case feedsGet.FullCommand():
		kingpin.FatalIfError(
			feedsGetHandler(client, *feedsGetName, *feedsGetSince, *feedsGetCurrent),
			"unable to get feed")
	case feedsList.FullCommand():
		kingpin.FatalIfError(
			feedsListHandler(client),
			"unable to list feeds")

		// Reactions
	case reactionsDefinitionsGet.FullCommand():
		kingpin.FatalIfError(
			reactionsDefinitionsGetHandler(client, *reactionsDefinitionsGetName),
			"unable to get reaction definition")
	case reactionsDefinitionsDelete.FullCommand():
		kingpin.FatalIfError(
			reactionsDefinitionsDeleteHandler(client, *reactionsDefinitionsDeleteName),
			"unable to delete reaction definition")
	case reactionsDefinitionsList.FullCommand():
		kingpin.FatalIfError(
			reactionsDefinitionsListHandler(client),
			"unable to list reaction definitions")
	}
}
