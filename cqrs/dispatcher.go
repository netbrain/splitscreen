package cqrs

import (
	"context"
	"fmt"
)

type Dispatcher interface {
	DispatchMessage(ctx context.Context, msg Message) error
}

func NewDefaultDispatcher(factory AggregateFactory, store EventStore, idGen IDGenerator) *DefaultDispatcher {
	return &DefaultDispatcher{
		AggregateFactory: factory,
		EventStore:       store,
		IDGenerator:      idGen,
	}
}

type DefaultDispatcher struct {
	AggregateFactory
	EventStore
	IDGenerator
}

func (d *DefaultDispatcher) DispatchMessage(ctx context.Context, msg Message) error {
	if msg.Meta().MessageType.IsEvent() {
		return d.dispatchEvent(ctx, msg)
	}

	if msg.Meta().MessageType.IsCommand() {
		return d.dispatchCommand(ctx, msg)
	}

	return fmt.Errorf("unknown instance type")
}

func (d *DefaultDispatcher) dispatchEvent(ctx context.Context, msg Message) error {
	if !msg.Meta().Replay {
		changeTracker := ChangeTrackerFromContext(ctx)
		if err := changeTracker.TrackChange(msg); err != nil {
			return err
		}
	}
	return nil
}

func (d *DefaultDispatcher) dispatchCommand(ctx context.Context, msg Message) error {
	aggr := d.GetAggregate(msg.Meta().AggregateType)
	if aggr == nil {
		return fmt.Errorf("unknown aggregate")
	}

	if msg.Meta().AggregateID != "" {
		if err := d.LoadAggregate(ctx, msg.Meta().AggregateMeta, aggr); err != nil {
			return err
		}
	} else {
		msg.Meta().AggregateID = d.NewID()
	}

	return aggr.Handle(ctx, msg)
}
