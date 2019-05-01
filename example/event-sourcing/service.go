package main

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/google/uuid"
	"github.com/marcusolsson/serialized-go"
)

type SerializedEventService struct {
	client *serialized.Client
}

func (s *SerializedEventService) SaveEvents(ctx context.Context, id OrderID, version int64, events ...OrderEvent) error {
	var res []*serialized.Event
	for _, ev := range events {
		b, err := json.Marshal(ev)
		if err != nil {
			return err
		}

		res = append(res, &serialized.Event{
			ID:   uuid.New().String(),
			Type: reflect.TypeOf(ev).Name(),
			Data: json.RawMessage(b),
		})
	}
	return s.client.Store(ctx, "order", string(id), version, res...)
}

func (s *SerializedEventService) Load(ctx context.Context, id OrderID) (OrderState, error) {
	agg, err := s.client.LoadAggregate(ctx, "order", string(id))
	if err != nil {
		return OrderState{}, err
	}
	return buildState(OrderID(agg.ID), agg.Version, agg.Events), nil
}
