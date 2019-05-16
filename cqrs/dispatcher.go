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
	if !msg.Meta().Replay {
		app := FromContext(ctx)
		if err := app.TrackChange(msg); err != nil {
			return err
		}
	}
	return nil
}

func dispatchCommand(ctx context.Context, msg Message) error {
	app := FromContext(ctx)
	aggr := app.GetAggregate(msg.Meta().AggregateType)
	if aggr == nil {
		return fmt.Errorf("unknown aggregate")
	}

	if msg.Meta().AggregateID != ""{
		if err := LoadAggregate(ctx,msg.Meta().AggregateMeta,aggr); err != nil {
			return err
		}
	}
	return aggr.Handle(ctx, msg)
}
