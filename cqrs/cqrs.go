package cqrs

import (
	"context"
	"github.com/netbrain/splitscreen/cqrs/json"
	"net/http"
)

var app *CQRS

func init() {
	app = &CQRS{
		Serializer:       json.NewSerializer(),
		Deserializer:     json.NewDeserializer(),
		EventStore:       NewMemoryEventStore(),
		MessageBus:       NewLocalMessageBus(),
		IDGenerator:      NewDefaultIDGenerator(),
		AggregateFactory: NewAggregateFactory(),
		MessageFactory:   NewMessageFactory(),
		ViewRepository:   NewViewRepository(),
	}
	app.ChangeTracker = NewChangeTracker(app.EventStore, app.MessageBus)
}

type CQRS struct {
	Serializer
	Deserializer
	EventStore
	MessageBus
	IDGenerator
	*AggregateFactory
	*MessageFactory
	*ChangeTracker
	*ViewRepository
}

func SetCQRS(a *CQRS) {
	app = a
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

func RegisterView(typ ViewType, v interface{}) {
	app.RegisterView(typ, v)
}

func GetView(typ ViewType) interface{} {
	return app.GetView(typ)
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

func Serialize(src interface{}) ([]byte, error) {
	return app.Serialize(src)
}

func Deserialize(buf []byte, dst interface{}) error {
	return app.Deserialize(buf, dst)
}
