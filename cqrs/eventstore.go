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

func (m *MemoryEventStore) Store(ctx context.Context, events ...Message) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for _, event := range events {
		rawEv, err := NewRawMessage(ctx,event)
		if err != nil {
			return err
		}
		m.events = append(m.events, rawEv)
	}
	return nil
}

func (m *MemoryEventStore) Load(ctx context.Context, id string, typ AggregateType) <-chan *EventLoadResult {
	app := FromContext(ctx)
	out := make(chan *EventLoadResult)
	go func(){
		defer close(out)
		m.mutex.Lock()
		defer m.mutex.Unlock()
		for _, pe := range m.events {
			if pe.Meta().AggregateID == id && pe.Meta().AggregateType == typ {
				pe.Replay = true
				impl := app.GetMessage(ctx,pe.MessageType)
				if err := pe.ToImplementation(ctx,impl); err != nil {
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
