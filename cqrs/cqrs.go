package cqrs

import (
	"context"
	"net/http"
)

var app *CQRS

func init() {
	app = &CQRS{
		EventStore:       NewMemoryEventStore(),
		MessageBus:       NewLocalMessageBus(),
		IDGenerator:      NewDefaultIDGenerator(),
		AggregateFactory: NewAggregateFactory(),
		MessageFactory:   NewMessageFactory(),
	}
	app.ChangeTracker = NewChangeTracker(app.EventStore, app.MessageBus)
}

type CQRS struct {
	EventStore
	MessageBus
	IDGenerator
	*AggregateFactory
	*MessageFactory
	*ChangeTracker
}

func Store(events ...Message) error {
	return app.Store(events...)
}

func Load(id string, typ AggregateType) ([]Message, error) {
	return app.Load(id, typ)
}

func RegisterAggregate(typ AggregateType, f func() AggregateRoot) {
	app.RegisterAggregate(typ, f)
}

func GetAggregate(typ AggregateType) AggregateRoot {
	return app.GetAggregate(typ)
}

func RegisterMessage(f func() Message) {
	app.RegisterMessage(f)
}

func GetMessage(typ MessageType) Message {
	return app.GetMessage(typ)
}

func NewID() string {
	return app.NewID()
}

func Subscribe(f AggregateHandleFunc, typ ...MessageType) {
	app.Subscribe(f, typ...)
}

func Emit(ctx context.Context, messages ...Message) error {
	return app.Emit(ctx, messages...)
}

func TrackChange(event Message) error {
	return app.TrackChange(event)
}

func CommitChanges(ctx context.Context) error {
	return app.CommitChanges(ctx)
}

func Middleware(next http.Handler) http.Handler {
	return app.Middleware(next)
}

func NewContext(ctx context.Context) context.Context {
	return app.NewContext(ctx)
}

func SetCQRS(a *CQRS) {
	app = a
}

func NewMessage(typ MessageType, aggregateId ...string) Message {
	var id string
	if len(aggregateId) == 0 || aggregateId[0] == "" {
		id = NewID()
	} else {
		id = aggregateId[0]
	}
	msg := GetMessage(typ)
	msg.Meta().AggregateID = id
	return msg
}
