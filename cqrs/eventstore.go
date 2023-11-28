package cqrs

import (
	"context"
	"sync"
	"time"
)

type EventLoadResult struct {
	Message Message
	Err     error
}

type EventStore interface {
	Store(ctx context.Context, events ...Message) error
	Load(ctx context.Context, id string, typ AggregateType) <-chan *EventLoadResult
	LoadAggregate(ctx context.Context, meta *AggregateMeta, dst AggregateRoot) error
	LoadAggregateUntilTime(ctx context.Context, meta *AggregateMeta, dst AggregateRoot, time time.Time) error
}

type MemoryEventStore struct {
	Serializer
	Deserializer
	MessageFactory
	events []*RawMessage
	mutex  sync.Mutex
}

func NewMemoryEventStore(s Serializer, d Deserializer, m MessageFactory) *MemoryEventStore {
	return &MemoryEventStore{
		Serializer:     s,
		Deserializer:   d,
		MessageFactory: m,
	}
}

func (m *MemoryEventStore) Store(ctx context.Context, events ...Message) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for _, event := range events {
		rawEv, err := NewRawMessage(m, event)
		if err != nil {
			return err
		}
		m.events = append(m.events, rawEv)
	}
	return nil
}

func (m *MemoryEventStore) Load(ctx context.Context, id string, typ AggregateType) <-chan *EventLoadResult {
	out := make(chan *EventLoadResult)
	go func() {
		defer close(out)
		m.mutex.Lock()
		defer m.mutex.Unlock()
		for _, pe := range m.events {
			if pe.Meta().AggregateID == id && pe.Meta().AggregateType == typ {
				pe.Replay = true
				impl := m.NewMessage(pe.MessageType)
				if err := pe.ToImplementation(m, impl); err != nil {
					out <- &EventLoadResult{
						Err: err,
					}
					return
				}
				out <- &EventLoadResult{
					Message: impl,
				}
			}
		}
	}()
	return out

}

func (m *MemoryEventStore) LoadAggregate(ctx context.Context, meta *AggregateMeta, dst AggregateRoot) error {
	return LoadAggregate(ctx, m, meta, dst)
}

func (m *MemoryEventStore) LoadAggregateUntilTime(ctx context.Context, meta *AggregateMeta, dst AggregateRoot, time time.Time) error {
	return LoadAggregateUntilTime(ctx, m, meta, dst, time)
}
