package cqrs

import (
	"context"
	"fmt"
	"github.com/netbrain/splitscreen/cqrs/json"
	"net/http"
)

type contextKey int

const (
	cqrsContextKey contextKey = iota
	changeTrackerContextKey
)

type DispatchFlags int

const (
	CustomAggregateID DispatchFlags = 1 << iota //when provided with a custom aggregate id, should not fail loading aggregate
)

type Dispatcher interface {
	DispatchMessage(ctx context.Context, msg Message) error
}

type App struct {
	Serializer
	Deserializer
	EventStore
	MessageBus
	IDGenerator
	AggregateFactory
	MessageFactory
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
		Serializer:       serializer,
		Deserializer:     deserializer,
		EventStore:       eventStore,
		MessageBus:       NewLocalMessageBus(),
		IDGenerator:      idGen,
		AggregateFactory: aggregateFactory,
		MessageFactory:   messageFactory,
		ViewRepository:   NewViewRepository(),
		ManagerRepository: NewManagerRepository(),
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
	return context.WithValue(context.WithValue(ctx, changeTrackerContextKey, NewChangeTracker(a)), cqrsContextKey, a)
}

func (a *App) DispatchMessage(ctx context.Context, msg Message, opts ...DispatchFlags) error {
	//TODO should we separate dispatch event/command, only provide public method for commands maybe?
	if msg.Meta().MessageType.IsEvent() {
		return a.dispatchEvent(ctx, msg)
	}

	if msg.Meta().MessageType.IsCommand() {
		var f DispatchFlags
		for i := range opts {
			f |= opts[i]
		}
		return a.dispatchCommand(ctx, msg,f)
	}

	return fmt.Errorf("unknown instance type")
}

func (a *App) dispatchEvent(ctx context.Context, msg Message) error {
	if !msg.Meta().Replay {
		changeTracker := ChangeTrackerFromContext(ctx)
		if err := changeTracker.TrackChange(msg); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) dispatchCommand(ctx context.Context, msg Message, opts DispatchFlags) error {
	aggr := a.GetAggregate(msg.Meta().AggregateType)
	if aggr == nil {
		return fmt.Errorf("unknown aggregate")
	}

	if opts&(CustomAggregateID) == 0 {
		if msg.Meta().AggregateID != "" {
			if err := a.LoadAggregate(ctx, msg.Meta().AggregateMeta, aggr); err != nil {
				return err
			}
		} else {
			msg.Meta().AggregateID = a.NewID()
		}
	}

	return aggr.Handle(ctx, msg)
}

func FromContext(ctx context.Context) *App {
	return ctx.Value(cqrsContextKey).(*App)
}

func ChangeTrackerFromContext(ctx context.Context) *ChangeTracker {
	return ctx.Value(changeTrackerContextKey).(*ChangeTracker)
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
