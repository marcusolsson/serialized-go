package main

import (
	"errors"

	"github.com/google/uuid"
)

type OrderStatus int

const (
	OrderStatusNew OrderStatus = iota
	OrderStatusPlaced
	OrderStatusCancelled
	OrderStatusShipped
	OrderStatusPaid
)

type OrderID string

func newOrderID() OrderID {
	return OrderID(uuid.New().String())
}

type CustomerID string

func newCustomerID() CustomerID {
	return CustomerID(uuid.New().String())
}

type TrackingNumber string

func newTrackingNumber() TrackingNumber {
	return TrackingNumber(uuid.New().String())
}

type Order struct {
	CustomerID CustomerID
	Status     OrderStatus
	Amount     Amount
}

func (o *Order) Place(amount Amount) (OrderPlacedEvent, error) {
	if o.Status != OrderStatusNew {
		return OrderPlacedEvent{}, errors.New("order already placed")
	}
	return OrderPlacedEvent{
		CustomerID: o.CustomerID,
		Amount:     amount,
	}, nil
}

func (o *Order) Pay(amount Amount) ([]OrderEvent, error) {
	if o.Status != OrderStatusPlaced {
		return nil, errors.New("order not placed")
	}
	if !amount.IsPositive() {
		return nil, errors.New("invalid amount")
	}

	var events []OrderEvent
	events = append(events, PaymentReceivedEvent{
		CustomerID: o.CustomerID,
		AmountPaid: amount,
	})

	if amount.LargerThanEq(o.Amount) {
		events = append(events, OrderFullyPaidEvent{
			CustomerID: o.CustomerID,
		})
	}

	return events, nil
}

func (o *Order) Ship(trackingNumber TrackingNumber) (OrderShippedEvent, error) {
	if o.Status != OrderStatusPaid {
		return OrderShippedEvent{}, errors.New("order not paid")
	}
	return OrderShippedEvent{
		CustomerID:     o.CustomerID,
		TrackingNumber: trackingNumber,
	}, nil
}

func (o *Order) Cancel(reason string) (OrderCancelledEvent, error) {
	if o.Status != OrderStatusPlaced {
		return OrderCancelledEvent{}, errors.New("order not placed")
	}
	return OrderCancelledEvent{
		CustomerID: o.CustomerID,
		Reason:     reason,
	}, nil
}
