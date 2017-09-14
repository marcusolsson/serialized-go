package main

import (
	"context"
	"encoding/json"
	"log"
	"reflect"

	serialized "github.com/marcusolsson/serialized-go"
	uuid "github.com/satori/go.uuid"
)

type OrderStatus int

const (
	OrderStatusPlaced OrderStatus = iota
	OrderStatusCancelled
)

type OrderState struct {
	AggregateID string
	Version     int64

	OrderStatus OrderStatus
	OrderAmount int64
	Paid        bool
}

func (s *OrderState) ApplyEvent(v interface{}) {
	switch ev := v.(type) {
	case OrderPlacedEvent:
		s.OrderStatus = OrderStatusPlaced
		s.OrderAmount = ev.OrderAmount
	case OrderPaidEvent:
		s.Paid = true
	case OrderCancelledEvent:
		s.OrderStatus = OrderStatusCancelled
	}
}

type Order struct {
	State OrderState
}

func (o *Order) Place(customerID string, amount int64) OrderPlacedEvent {
	return OrderPlacedEvent{
		CustomerID:  customerID,
		OrderAmount: amount,
	}
}

func (o *Order) Cancel(reason string) OrderCancelledEvent {
	return OrderCancelledEvent{
		Reason: reason,
	}
}

type OrderPlacedEvent struct {
	CustomerID  string `json:"customerId"`
	OrderAmount int64  `json:"orderAmount"`
}

type OrderPaidEvent struct{}

type OrderCancelledEvent struct {
	Reason string `json:"reason"`
}

type OrderClient struct {
	client *serialized.Client
}

func (c *OrderClient) Load(ctx context.Context, orderID string) OrderState {
	agg, err := c.client.LoadAggregate(ctx, "order", orderID)
	if err != nil {
		log.Fatal(err)
	}

	var state OrderState

	state.AggregateID = orderID

	for _, ev := range agg.Events {
		switch ev.Type {
		case "OrderPlacedEvent":
			var event OrderPlacedEvent
			json.Unmarshal(ev.Data, &event)
			state.ApplyEvent(event)
		case "OrderPaidEvent":
			var event OrderPaidEvent
			json.Unmarshal(ev.Data, &event)
			state.ApplyEvent(event)
		case "OrderCancelledEvent":
			var event OrderCancelledEvent
			json.Unmarshal(ev.Data, &event)
			state.ApplyEvent(event)
		}
	}

	return state
}

func typeName(v interface{}) string {
	return reflect.TypeOf(v).Name()
}

func (c *OrderClient) Save(ctx context.Context, orderID string, version int64, event interface{}) error {
	return c.client.Store(ctx, "order", orderID, version,
		&serialized.Event{
			ID:   uuid.NewV4().String(),
			Type: typeName(event),
			Data: mustMarshal(event),
		})
}

func mustMarshal(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}
