package cqrs

import (
	"context"
	"fmt"
	"net/http"
)

var ErrMetaNotPresent = fmt.Errorf("meta not initialized on aggregate")

type AggregateHandleFunc func(ctx context.Context, msg Message) error

type AggregateRoot interface {
	Meta() *AggregateMeta
	Handle(ctx context.Context, msg Message) error
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

	events, err := Load(meta.AggregateID, meta.AggregateType)
	if err != nil {
		return err
	}
	for _, e := range events {
		if err := aggr.Handle(ctx, e); err != nil {
			return err
		}
	}

	*aggrMeta = *meta
	aggrMeta.loaded = true
	return nil
}

type contextKey int

const (
	changeTracker contextKey = iota
)

type ChangeTracker struct {
	store   EventStore
	bus     MessageBus
	changes []Message
}

func NewChangeTracker(store EventStore, bus MessageBus) *ChangeTracker {
	return &ChangeTracker{store: store, bus: bus}
}

func (c *ChangeTracker) TrackChange(event Message) error {
	if !event.Meta().MessageType.IsEvent() {
		return fmt.Errorf("expected event")
	}

	c.changes = append(c.changes, event)
	return nil
}

func (c *ChangeTracker) CommitChanges(ctx context.Context) error {
	err := c.store.Store(c.changes...)
	if err != nil {
		return err
	}
	for _, msg := range c.changes {
		if err := c.bus.Emit(ctx, msg); err != nil {
			return err
		}
	}
	c.changes = nil
	return nil
}

func (c *ChangeTracker) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(c.NewContext(r.Context()))
		next.ServeHTTP(w, r)
		if err := c.CommitChanges(r.Context()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

func (c *ChangeTracker) NewContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, changeTracker, c)
}

func ChangeTrackerFromContext(ctx context.Context) *ChangeTracker {
	if ct, ok := ctx.Value(changeTracker).(*ChangeTracker); ok {
		return ct
	}
	return nil
}
