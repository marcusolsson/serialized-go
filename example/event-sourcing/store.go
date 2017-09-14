package main

import (
	"context"
	"log"
	"os"

	serialized "github.com/marcusolsson/serialized-go"
	uuid "github.com/satori/go.uuid"
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

	orderID := uuid.NewV4().String()

	orderClient := OrderClient{client: client}

	// Create ...
	initState := OrderState{AggregateID: orderID}
	order := Order{State: initState}

	// ... and place a new order.
	ope := order.Place("klarna", 1000)
	if err := orderClient.Save(ctx, initState.AggregateID, initState.Version, ope); err != nil {
		log.Fatal(err)
	}

	// Load ...
	toCancelState := orderClient.Load(ctx, orderID)
	orderToCancel := Order{State: toCancelState}

	// ... and cancel order.
	oce := orderToCancel.Cancel("DOA")
	if err := orderClient.Save(ctx, initState.AggregateID, initState.Version, oce); err != nil {
		log.Fatal(err)
	}
}
