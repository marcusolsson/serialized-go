package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

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

	ctx := context.Background()

	err := client.Feed(ctx, "order", 0, func(entry *serialized.FeedEntry) {
		log.Printf("Processing entry with sequenceNumber: %d", entry.SequenceNumber)

		for _, event := range entry.Events {
			var data map[string]interface{}
			if err := json.Unmarshal(event.Data, &data); err != nil {
				log.Printf("Unable to unmarshal event data: %v", err)
			}

			switch event.Type {
			case "OrderPlacedEvent":
				log.Printf("An order with ID [%s] was placed by customer [%v]\n", entry.AggregateID, data["customerId"])
			case "OrderPaidEvent":
				log.Printf("The order with ID [%s] was paid, amountPaid: %v, amountLeft: %v\n", entry.AggregateID, data["amountPaid"], data["amountLeft"])
			case "OrderShippedEvent":
				log.Printf("The order with ID [%s] was shipped, trackingNumber: %v\n", entry.AggregateID, data["trackingNumber"])
			case "OrderCancelledEvent":
				log.Printf("The order with ID [%s] was cancelled, reason: %v\n", entry.AggregateID, data["reason"])
			case "PaymentReceivedEvent":
				log.Printf("The order with ID [%s] received payment: %v\n", entry.AggregateID, data["amountPaid"])
			case "OrderFullyPaid":
				log.Printf("The order with ID [%s] is fully paid\n", entry.AggregateID)
			default:
				log.Printf("Don't know how to handle events of type: %s", event.Type)
			}
		}
	})
	if err != nil {
		log.Fatal(err)
	}
}
