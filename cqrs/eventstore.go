package cqrs

import (
	"sync"
	"fmt"
)

var events []*Event
var mutex sync.Mutex

func Store(e *Event) error {
	mutex.Lock()
	defer mutex.Unlock()
	for _,pe := range events {
		if pe.Aggregate.ID == e.Aggregate.ID && pe.Aggregate.Version == e.Aggregate.Version + 1 {
			return fmt.Errorf("version conflict")
		}
	}
	events = append(events,e)
	return nil
}

func Load(id string, typ AggregateType) []*Event {
	var out []*Event
	mutex.Lock()
	defer mutex.Unlock()
	for _,pe := range events {
		if pe.Aggregate.ID == id && pe.Aggregate.Type == typ {
			out = append(out,pe)
		}
	}
	return out
}