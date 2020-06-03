package cqrs

import (
	"context"
	"fmt"
)

type DefaultDispatcher struct {}

func NewDefaultDispatcher() Dispatcher {
	return &DefaultDispatcher{}
}

func (d *DefaultDispatcher) DispatchMessage(ctx context.Context, msg Message, opts ...DispatchFlags) error {
	//TODO should we separate dispatch event/command, only provide public method for commands maybe?
	if msg.Meta().MessageType.IsEvent() {
		return d.dispatchEvent(ctx, msg)
	}

	if msg.Meta().MessageType.IsCommand() {
		var f DispatchFlags
		for i := range opts {
			f |= opts[i]
		}
		return d.dispatchCommand(ctx, msg,f)
	}

	return fmt.Errorf("unknown instance type")
}

func (d *DefaultDispatcher) dispatchEvent(ctx context.Context, msg Message) error {
	if !msg.Meta().Replay {
		changeTracker := ChangeTrackerFromContext(ctx)
		if err := changeTracker.TrackChange(msg); err != nil {
			return err
		}
	}
	return nil
}

func (d *DefaultDispatcher) dispatchCommand(ctx context.Context, msg Message, opts DispatchFlags) error {
	a := FromContext(ctx)
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



