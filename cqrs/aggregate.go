package cqrs

import (
	"context"
	"fmt"
)

var ErrMetaNotPresent = fmt.Errorf("meta not initialized on aggregate")
var ErrNoID = fmt.Errorf("no id specified on aggregate")
var ErrNoEvents = fmt.Errorf("no events")

type AggregateHandleFunc func(ctx context.Context, msg Message) error

type AggregateRoot interface {
	Meta() *AggregateMeta
	Handle(ctx context.Context, msg Message) error
	Register(app *App)
}

type AggregateType string

type AggregateMeta struct {
	AggregateID   string        `json:"aggregateId"`
	AggregateType AggregateType `json:"aggregateType"`
	loaded        bool
}

func LoadAggregate(ctx context.Context, meta *AggregateMeta, aggr AggregateRoot) error {
	aggrMeta := aggr.Meta()
	if aggrMeta == nil {
		return ErrMetaNotPresent
	}

	if aggrMeta.loaded {
		return nil
	}

	if meta.AggregateID == "" {
		return ErrNoID
	}

	app := FromContext(ctx)
	result := app.Load(ctx, meta.AggregateID, meta.AggregateType)
	var count int
	for e := range result {
		count++
		if e.Err != nil {
			return e.Err
		}
		if err := aggr.Handle(ctx, e.Message); err != nil {
			return err
		}
	}

	if count == 0 {
		return ErrNoEvents
	}

	*aggrMeta = *meta
	aggrMeta.loaded = true
	return nil
}

type ChangeTracker struct {
	changes []Message
}

func NewChangeTracker() *ChangeTracker {
	return &ChangeTracker{}
}

func (c *ChangeTracker) TrackChange(event Message) error {
	if !event.Meta().MessageType.IsEvent() {
		return fmt.Errorf("expected event")
	}

	c.changes = append(c.changes, event)
	return nil
}

func (c *ChangeTracker) CommitChanges(ctx context.Context) error {
	app := FromContext(ctx)
	err := app.Store(ctx, c.changes...)
	if err != nil {
		return err
	}
	for _, msg := range c.changes {
		if err := app.Emit(ctx, msg); err != nil {
			return err
		}
	}
	c.changes = nil
	return nil
}
