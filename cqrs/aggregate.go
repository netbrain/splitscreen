package cqrs

import (
	"fmt"
	"context"
)

type AggregateType string

type AggregateID struct {
	ID string
	Version int
	Type AggregateType
}

type Aggregate struct {
	*AggregateID
	Impl           AggregateImpl
	uncommitted    []*Event
	uncommittedMap map[string]*Event
	replayMode     bool
}

type AggregateImpl interface {
	CommandHandler
	EventHandler
}

func (a *Aggregate) ReplayMode() bool {
	return a.replayMode
}

func (a *Aggregate) Apply(ctx context.Context,e *Event) error {
	if a.Version == 0 {
		a.ID = e.ID
	}
	if !a.ReplayMode() {
		if a.Version != e.Aggregate.Version {
			return fmt.Errorf("version mismatch during apply")
		}
		a.addEventToChangeList(e)
	}else{
		a.Version++
		if a.Version != e.Aggregate.Version {
			return fmt.Errorf("version mismatch during replay")
		}
	}
	return a.Impl.Apply(ctx,e)
}


func (a *Aggregate) addEventToChangeList(e *Event) {
	if a.uncommittedMap == nil {
		a.uncommittedMap = make(map[string]*Event)
	}
	if _,ok := a.uncommittedMap[e.ID]; !ok {
		a.uncommittedMap[e.ID] = e
		a.uncommitted = append(a.uncommitted, e)
	}
}

func (a *Aggregate) Pending(e *Event) bool {
	return a.uncommittedMap != nil && a.uncommittedMap[e.ID] != nil
}

func (a *Aggregate) Commit() error {
	for _,e := range a.uncommitted {
		e.Aggregate.Version = a.Version + 1
		if err := eventstore.Store(e); err != nil {
			return err
		}
		a.Version++
		BroadcastEvent(e)
	}
	return nil
}

func (a *Aggregate) Handle(ctx context.Context,c *Command) error {
	if a.Version != c.Aggregate.Version {
		return fmt.Errorf("command/aggregate version mismatch: %d != %d",c.Aggregate.Version, a.Version)
	}
	return a.Impl.Handle(ctx,c)
}
