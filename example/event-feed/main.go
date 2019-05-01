package main

import (
	"context"
	"encoding/json"
	"fmt"
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

	err := client.Feed(context.Background(), "order", 0, func(entry *serialized.FeedEntry) {
		for _, event := range entry.Events {
			switch event.Type {
			case "OrderPlacedEvent":
				var orderPlaced struct {
					CustomerID string `json:"customerId"`
				}
				if err := json.Unmarshal(event.Data, &orderPlaced); err != nil {
					fmt.Printf("Unable to unmarshal event data: %v", err)
				}

				fmt.Printf("An order with ID %s was placed by %s\n", entry.AggregateID, orderPlaced.CustomerID)
			case "OrderPaidEvent":
				fmt.Printf("The order with ID %s was paid\n", entry.AggregateID)
			case "OrderCancelledEvent":
				fmt.Printf("The order with ID %s was cancelled\n", entry.AggregateID)
			default:
				fmt.Println("Don't know how to handle events of type:", event.Type)
			}
		}
	})
	if err != nil {
		log.Fatal(err)
	}
}
