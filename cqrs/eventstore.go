package cqrs

import (
	"context"
	"sync"
)

type EventLoadResult struct {
	Message Message
	Err error
}

type EventStore interface {
	Store(ctx context.Context, events ...Message) error
	Load(ctx context.Context,id string, typ AggregateType) <-chan *EventLoadResult
}

type MemoryEventStore struct {
	events []*RawMessage
	mutex  sync.Mutex
}

func NewMemoryEventStore() *MemoryEventStore {
	return &MemoryEventStore{}
}

func (m *MemoryEventStore) Store(_ context.Context, events ...Message) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for _, event := range events {
		rawEv, err := NewRawMessage(event)
		if err != nil {
			return err
		}
		m.events = append(m.events, rawEv)
	}
	return nil
}

func (m *MemoryEventStore) Load(_ context.Context, id string, typ AggregateType) <-chan *EventLoadResult {
	out := make(chan *EventLoadResult)
	go func(){
		defer close(out)
		m.mutex.Lock()
		defer m.mutex.Unlock()
		for _, pe := range m.events {
			if pe.Meta().AggregateID == id && pe.Meta().AggregateType == typ {
				pe.Replay = true
				impl := GetMessage(pe.MessageType)
				if err := pe.ToImplementation(impl); err != nil {
					out<-&EventLoadResult{
						Err:     err,
					}
					return
				}
				out<-&EventLoadResult{
					Message:impl,
				}
			}
		}
	}()
	return out

}
