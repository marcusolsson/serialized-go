package main

import (
	"encoding/json"

	"github.com/marcusolsson/serialized-go"
)

type OrderState struct {
	ID      OrderID
	Version int64

	CustomerID     CustomerID
	Status         OrderStatus
	Amount         Amount
	CancelReason   string
	TrackingNumber TrackingNumber
}

func buildState(id OrderID, version int64, events []*serialized.Event) OrderState {
	state := OrderState{
		ID:      id,
		Version: version,
	}

	for _, e := range events {
		switch e.Type {
		case "OrderPlacedEvent":
			var ev OrderPlacedEvent
			json.Unmarshal(e.Data, &ev)

			state.CustomerID = ev.CustomerID
			state.Status = OrderStatusPlaced
			state.Amount = ev.Amount
		case "PaymentReceivedEvent":
			var ev PaymentReceivedEvent
			json.Unmarshal(e.Data, &ev)

			state.Amount = state.Amount - ev.AmountPaid
		case "OrderFullyPaidEvent":
			var ev OrderFullyPaidEvent
			json.Unmarshal(e.Data, &ev)

			state.Status = OrderStatusPaid
		case "OrderCancelledEvent":
			var ev OrderCancelledEvent
			json.Unmarshal(e.Data, &ev)

			state.Status = OrderStatusCancelled
			state.CancelReason = ev.Reason
		case "OrderShippedEvent":
			var ev OrderShippedEvent
			json.Unmarshal(e.Data, &ev)

			state.Status = OrderStatusShipped
			state.TrackingNumber = ev.TrackingNumber
		}
	}

	return state
}
