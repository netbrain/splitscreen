package cqrs

import "context"

type CommandImpl interface {
	Dispatch(context.Context,string,int) error
}

type CommandHandlerFunc func(*Command) error

type CommandHandler interface {
	Handle(context.Context, *Command) error
}

type CommandType string

type Command struct {
	ID            string
	CorrelationID string
	CausationID   string
	Type          CommandType
	Aggregate     AggregateID
	Impl          CommandImpl
}
