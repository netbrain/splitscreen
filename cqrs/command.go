package cqrs

type CommandImpl interface {
	Dispatch(string,int) error
}

type CommandHandlerFunc func(*Command) error

type CommandHandler interface {
	Handle(*Command) error
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
