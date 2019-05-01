package main

import (
	"context"
	"log"
	"os"

	"github.com/marcusolsson/serialized-go"
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

	store := SerializedEventService{
		client: client,
	}

	ctx := context.Background()

	// ===============================================

	orderID := newOrderID()
	customerID := newCustomerID()

	// Create order ...

	orderInitState := buildState(orderID, 0, nil)
	order := Order{
		CustomerID: customerID,
		Status:     orderInitState.Status,
	}

	orderPlaced, _ := order.Place(4321)

	if err := store.SaveEvents(ctx, orderInitState.ID, orderInitState.Version, orderPlaced); err != nil {
		log.Fatal(err)
	}

	// -----------------------------------------------

	// Load ...
	orderToCancelState, _ := store.Load(ctx, orderID)
	order = Order{
		CustomerID: customerID,
		Status:     orderToCancelState.Status,
		Amount:     orderToCancelState.Amount,
	}

	// ... and cancel order
	orderCancelled, _ := order.Cancel("DOA")

	if err := store.SaveEvents(ctx, orderToCancelState.ID, orderToCancelState.Version, orderCancelled); err != nil {
		log.Fatal(err)
	}

	// ===============================================

	orderID2 := newOrderID()

	// Create

	orderInitState1 := buildState(orderID2, 0, nil)
	order1 := Order{
		CustomerID: customerID,
		Status:     orderInitState1.Status,
		Amount:     orderInitState1.Amount,
	}

	orderPlaced1, _ := order1.Place(1234)

	if err := store.SaveEvents(ctx, orderInitState1.ID, orderInitState1.Version, orderPlaced1); err != nil {
		log.Fatal(err)
	}

	// -----------------------------------------------

	orderToPayState, _ := store.Load(ctx, orderID2)
	orderToPay := Order{
		CustomerID: customerID,
		Status:     orderToPayState.Status,
		Amount:     orderToPayState.Amount,
	}

	events, _ := orderToPay.Pay(1234)

	if err := store.SaveEvents(ctx, orderToPayState.ID, orderToPayState.Version, events...); err != nil {
		log.Fatal(err)
	}

	// -----------------------------------------------

	orderToShipState, _ := store.Load(ctx, orderID2)
	orderToShip := Order{
		CustomerID: customerID,
		Status:     orderToShipState.Status,
		Amount:     orderToShipState.Amount,
	}

	trackingNumber := newTrackingNumber()
	orderShipped, _ := orderToShip.Ship(trackingNumber)

	if err := store.SaveEvents(ctx, orderToShipState.ID, orderToShipState.Version, orderShipped); err != nil {
		log.Fatal(err)
	}
}
