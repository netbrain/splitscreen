package cqrs

import (
	"context"
)

type AggregateHandleFunc func(ctx context.Context, msg Message) error

type AggregateRoot interface {
	Registerable
	Meta() *AggregateMeta
	Handle(ctx context.Context, msg Message) error
}

type AggregateType string

type AggregateMeta struct {
	AggregateID   string        `json:"aggregateId"`
	AggregateType AggregateType `json:"aggregateType"`
	loaded        bool
}
