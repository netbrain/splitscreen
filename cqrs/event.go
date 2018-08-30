package cqrs

import (
	"context"
)

type EventImpl interface{
	Apply(context.Context, *Aggregate, *Command) error
}

type EventHandlerFunc func(*Event) error

type EventHandler interface {
	Apply(context.Context, *Event) error
}

type EventType string

type Event struct {
	ID            string
	CorrelationID string
	CausationID   string
	Type          EventType
	Aggregate     AggregateID
	Impl          EventImpl
}







