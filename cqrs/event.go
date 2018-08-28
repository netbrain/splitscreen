package cqrs

type EventImpl interface{
	Apply(*Aggregate, *Command) error
}

type EventHandlerFunc func(*Event) error

type EventHandler interface {
	Apply(*Event) error
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
