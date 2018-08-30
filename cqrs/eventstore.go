package cqrs

import (
	"sync"
	"fmt"
)

type EventStore interface {
	Store(e *Event) error
	Load(id string, typ AggregateType) []*Event
}

type MemoryEventStore struct {
	events []*Event
	mutex sync.Mutex
}


func (m *MemoryEventStore) Store(e *Event) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for _,pe := range m.events {
		if pe.Aggregate.ID == e.Aggregate.ID && pe.Aggregate.Version == e.Aggregate.Version + 1 {
			return fmt.Errorf("version conflict")
		}
	}
	m.events = append(m.events,e)
	return nil
}

func (m *MemoryEventStore) Load(id string, typ AggregateType) []*Event {
	var out []*Event
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for _,pe := range m.events {
		if pe.Aggregate.ID == id && pe.Aggregate.Type == typ {
			out = append(out,pe)
		}
	}
	return out
}