package cqrs

import (
	"context"
	"github.com/netbrain/splitscreen/cqrs/json"
	"net/http"
)

type contextKey int

const cqrsContextKey contextKey = 0

type App struct {
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

func New(a *App) *App {
	def := &App{
		Serializer:       json.NewSerializer(),
		Deserializer:     json.NewDeserializer(),
		EventStore:       NewMemoryEventStore(),
		MessageBus:       NewLocalMessageBus(),
		IDGenerator:      NewDefaultIDGenerator(),
		AggregateFactory: NewAggregateFactory(),
		MessageFactory:   NewMessageFactory(),
		ViewRepository:   NewViewRepository(),
		ChangeTracker: 	  NewChangeTracker(),
	}
	if a == nil {
		return def
	}
	if a.Serializer == nil {
		a.Serializer = def.Serializer
	}
	if a.Deserializer == nil {
		a.Deserializer = def.Deserializer
	}
	if a.EventStore == nil {
		a.EventStore = def.EventStore
	}
	if a.MessageBus == nil {
		a.MessageBus = def.MessageBus
	}
	if a.IDGenerator == nil {
		a.IDGenerator = def.IDGenerator
	}
	if a.AggregateFactory == nil {
		a.AggregateFactory = def.AggregateFactory
	}
	if a.MessageFactory == nil {
		a.MessageFactory = def.MessageFactory
	}
	if a.ViewRepository == nil {
		a.ViewRepository = def.ViewRepository
	}
	if a.ChangeTracker == nil {
		a.ChangeTracker = def.ChangeTracker
	}
	return a
}

func (c *App) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(c.NewContext(r.Context()))
		next.ServeHTTP(w, r)
		if err := c.CommitChanges(r.Context()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

func (c *App) NewContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, cqrsContextKey, c)
}

func FromContext(ctx context.Context) *App {
	/*if c, ok := ctx.Value(cqrsContextKey).(*App); ok {
		return c
	}
	return nil*/
	return ctx.Value(cqrsContextKey).(*App)
}
