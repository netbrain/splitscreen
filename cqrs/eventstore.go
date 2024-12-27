package cqrs

import (
	"context"
	"sync"
)

type EventLoadResult struct {
	Message Message
	Err     error
}

type EventStore interface {
	Store(ctx context.Context, events ...Message) error
	Load(ctx context.Context, id string, typ AggregateType, until ... UntilFunc) <-chan *EventLoadResult
	LoadAggregate(ctx context.Context, meta *AggregateMeta, dst AggregateRoot, until ... UntilFunc) error
	Replay(ctx context.Context) <-chan *EventLoadResult
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

func (m *MemoryEventStore) Load(ctx context.Context, id string, typ AggregateType, until ...UntilFunc) <-chan *EventLoadResult {
	out := make(chan *EventLoadResult)
	go func() {
		defer close(out)
		m.mutex.Lock()
		defer m.mutex.Unlock()
		for _, pe := range m.events {
			if pe.Meta().AggregateID == id && pe.Meta().AggregateType == typ {
				if len(until) > 0 && !until[0](pe){
					break
				}

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

func (m *MemoryEventStore) Replay(ctx context.Context) <-chan *EventLoadResult {
	panic("implement me")
}

func (m *MemoryEventStore) LoadAggregate(ctx context.Context, meta *AggregateMeta, dst AggregateRoot, until ... UntilFunc) error {
	return LoadAggregate(ctx, m, meta, dst, until...)
}


