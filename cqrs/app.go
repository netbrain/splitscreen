package cqrs

import (
	"context"
	"github.com/netbrain/splitscreen/cqrs/json"
	"net/http"
)

type contextKey int

const (
	cqrsContextKey contextKey = iota
	changeTrackerContextKey
)

type App struct {
	Serializer
	Deserializer
	Dispatcher
	EventStore
	MessageBus
	IDGenerator
	AggregateFactory
	MessageFactory
	ChangeTrackerFactory
	*ViewRepository
	*ManagerRepository
}

func New(a *App) *App {
	idGen := NewDefaultIDGenerator()
	serializer, deserializer := json.New()
	aggregateFactory := NewDefaultAggregateFactory()
	messageFactory := NewDefaultMessageFactory(idGen)
	eventStore := NewMemoryEventStore(serializer, deserializer, messageFactory)
	def := &App{
		Serializer:           serializer,
		Deserializer:         deserializer,
		Dispatcher:           NewDefaultDispatcher(),
		ChangeTrackerFactory: NewDefaultChangeTrackerFactory(),
		EventStore:           eventStore,
		MessageBus:           NewLocalMessageBus(),
		IDGenerator:          idGen,
		AggregateFactory:     aggregateFactory,
		MessageFactory:       messageFactory,
		ViewRepository:       NewViewRepository(),
		ManagerRepository:    NewManagerRepository(),
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
	if a.ManagerRepository == nil {
		a.ManagerRepository = def.ManagerRepository
	}
	if a.Dispatcher == nil {
		a.Dispatcher = def.Dispatcher
	}
	if a.ChangeTrackerFactory == nil {
		a.ChangeTrackerFactory = def.ChangeTrackerFactory
	}
	return a
}

func (a *App) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(a.NewContext(r.Context()))
		next.ServeHTTP(w, r)

		ct := ChangeTrackerFromContext(r.Context())
		if err := ct.CommitChanges(r.Context()); err != nil {
			panic(err)
		}
	})
}

func (a *App) NewContext(ctx context.Context) context.Context {
	return context.WithValue(context.WithValue(ctx, changeTrackerContextKey, a.ChangeTrackerFactory.NewChangeTracker()), cqrsContextKey, a)
}

func FromContext(ctx context.Context) *App {
	return ctx.Value(cqrsContextKey).(*App)
}

func ChangeTrackerFromContext(ctx context.Context) ChangeTracker {
	return ctx.Value(changeTrackerContextKey).(ChangeTracker)
}

func Handle(ctx context.Context, msg Message) error {
	app := FromContext(ctx)
	aggr := app.GetAggregate(msg.Meta().AggregateType)
	err := app.LoadAggregate(ctx, msg.Meta().AggregateMeta, aggr)
	if err != nil {
		return err
	}

	return aggr.Handle(ctx, msg)
}
