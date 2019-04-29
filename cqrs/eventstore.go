package cqrs

import (
	"sync"
)

type EventStore interface {
	Store(events ...Message) error
	Load(id string, typ AggregateType) ([]Message, error)
}

type MemoryEventStore struct {
	events []*RawMessage
	mutex  sync.Mutex
}

func NewMemoryEventStore() *MemoryEventStore {
	return &MemoryEventStore{}
}

func (m *MemoryEventStore) Store(events ...Message) error {
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

func (m *MemoryEventStore) Load(id string, typ AggregateType) ([]Message, error) {
	var out []Message
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for _, pe := range m.events {
		if pe.Meta().AggregateID == id && pe.Meta().AggregateType == typ {
			pe.replay = true
			impl := GetMessage(pe.MessageType)
			if err := pe.ToImplementation(impl); err != nil {
				return nil, err
			}
			out = append(out, impl)
		}
	}
	return out, nil
}
