package cqrs

import (
	"context"
	"fmt"
)

func DispatchMessage(ctx context.Context, msg Message) error {
	if msg.Meta().MessageType.IsEvent() {
		return dispatchEvent(ctx, msg)
	}

	if msg.Meta().MessageType.IsCommand() {
		return dispatchCommand(ctx, msg)
	}

	return fmt.Errorf("unknown instance type")
}

func dispatchEvent(ctx context.Context, msg Message) error {
	if !msg.Meta().Replayed() {
		ct := ChangeTrackerFromContext(ctx)
		if ct == nil {
			return fmt.Errorf("no change tracker on context")
		}
		if err := ct.TrackChange(msg); err != nil {
			return err
		}
	}
	return nil
}

func dispatchCommand(ctx context.Context, msg Message) error {
	aggr := GetAggregate(msg.Meta().AggregateType)
	if aggr == nil {
		return fmt.Errorf("unknown aggregate")
	}
	return aggr.Handle(ctx, msg)
}
