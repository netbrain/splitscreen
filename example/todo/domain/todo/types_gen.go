// CODE GENERATED, DO NOT MODIFY

package todo

import (
	"context"
	"fmt"

	"github.com/netbrain/splitscreen/cqrs"
)

const (
	AggregateType cqrs.AggregateType = "todo.Aggregate"

	CreateCommandType cqrs.MessageType = "todo.CreateCommand"

	EditCommandType cqrs.MessageType = "todo.EditCommand"

	CompleteCommandType cqrs.MessageType = "todo.CompleteCommand"

	CreatedEventType cqrs.MessageType = "todo.CreatedEvent"

	EditedEventType cqrs.MessageType = "todo.EditedEvent"

	CompletedEventType cqrs.MessageType = "todo.CompletedEvent"
)

func (a *Aggregate) Register(app *cqrs.App) {
	app.RegisterAggregate(AggregateType, func() cqrs.AggregateRoot {
		return &Aggregate{AggregateMeta: &cqrs.AggregateMeta{}}
	})

	app.RegisterMessage(func() cqrs.Message {
		return &CreatedEvent{
			MessageMeta: &cqrs.MessageMeta{
				AggregateMeta: &cqrs.AggregateMeta{
					AggregateType: AggregateType,
				},
				MessageType: CreatedEventType,
			},
		}
	})

	app.RegisterMessage(func() cqrs.Message {
		return &EditedEvent{
			MessageMeta: &cqrs.MessageMeta{
				AggregateMeta: &cqrs.AggregateMeta{
					AggregateType: AggregateType,
				},
				MessageType: EditedEventType,
			},
		}
	})

	app.RegisterMessage(func() cqrs.Message {
		return &CompletedEvent{
			MessageMeta: &cqrs.MessageMeta{
				AggregateMeta: &cqrs.AggregateMeta{
					AggregateType: AggregateType,
				},
				MessageType: CompletedEventType,
			},
		}
	})

	app.RegisterMessage(func() cqrs.Message {
		return &CreateCommand{
			MessageMeta: &cqrs.MessageMeta{
				AggregateMeta: &cqrs.AggregateMeta{
					AggregateType: AggregateType,
				},
				MessageType: CreateCommandType,
			},
		}
	})

	app.RegisterMessage(func() cqrs.Message {
		return &EditCommand{
			MessageMeta: &cqrs.MessageMeta{
				AggregateMeta: &cqrs.AggregateMeta{
					AggregateType: AggregateType,
				},
				MessageType: EditCommandType,
			},
		}
	})

	app.RegisterMessage(func() cqrs.Message {
		return &CompleteCommand{
			MessageMeta: &cqrs.MessageMeta{
				AggregateMeta: &cqrs.AggregateMeta{
					AggregateType: AggregateType,
				},
				MessageType: CompleteCommandType,
			},
		}
	})

}

func (a *Aggregate) Meta() *cqrs.AggregateMeta {
	return a.AggregateMeta
}

func (a *Aggregate) Handle(ctx context.Context, msg cqrs.Message) (err error) {
	var event cqrs.Message
	switch msg.Meta().MessageType {

	case CreateCommandType:
		event, err = a.HandleCreateCommand(ctx, *msg.(*CreateCommand))

	case EditCommandType:
		event, err = a.HandleEditCommand(ctx, *msg.(*EditCommand))

	case CompleteCommandType:
		event, err = a.HandleCompleteCommand(ctx, *msg.(*CompleteCommand))

	default:
		return a.Apply(ctx, msg)
	}
	if err != nil {
		return err
	}
	meta := event.Meta()
	newMeta := cqrs.FromContext(ctx).NewMessage(meta.MessageType, msg.Meta().AggregateID).Meta()
	*meta = *newMeta
	return a.Apply(ctx, event)
}

func (a *Aggregate) Apply(ctx context.Context, msg cqrs.Message) (err error) {
	switch msg.Meta().MessageType {

	case CreatedEventType:
		err = a.ApplyCreatedEvent(ctx, *msg.(*CreatedEvent))

	case EditedEventType:
		err = a.ApplyEditedEvent(ctx, *msg.(*EditedEvent))

	case CompletedEventType:
		err = a.ApplyCompletedEvent(ctx, *msg.(*CompletedEvent))

	default:
		return fmt.Errorf("unknown message type")
	}

	if err != nil {
		return err
	}

	if msg.Meta().Replay {
		return nil
	}
	return cqrs.FromContext(ctx).DispatchMessage(ctx, msg)
}

func NewCreateCommandMessage(ctx context.Context, data CreateCommand, aggregateId ...string) *CreateCommand {
	data.MessageMeta = cqrs.FromContext(ctx).NewMessage(CreateCommandType, aggregateId...).Meta()
	return &data
}

func (e *CreateCommand) Meta() *cqrs.MessageMeta {
	return e.MessageMeta
}

func NewEditCommandMessage(ctx context.Context, data EditCommand, aggregateId ...string) *EditCommand {
	data.MessageMeta = cqrs.FromContext(ctx).NewMessage(EditCommandType, aggregateId...).Meta()
	return &data
}

func (e *EditCommand) Meta() *cqrs.MessageMeta {
	return e.MessageMeta
}

func NewCompleteCommandMessage(ctx context.Context, data CompleteCommand, aggregateId ...string) *CompleteCommand {
	data.MessageMeta = cqrs.FromContext(ctx).NewMessage(CompleteCommandType, aggregateId...).Meta()
	return &data
}

func (e *CompleteCommand) Meta() *cqrs.MessageMeta {
	return e.MessageMeta
}

func NewCreatedEventMessage(ctx context.Context, data CreatedEvent, aggregateId ...string) *CreatedEvent {
	data.MessageMeta = cqrs.FromContext(ctx).NewMessage(CreatedEventType, aggregateId...).Meta()
	return &data
}

func (e *CreatedEvent) Meta() *cqrs.MessageMeta {
	return e.MessageMeta
}

func NewEditedEventMessage(ctx context.Context, data EditedEvent, aggregateId ...string) *EditedEvent {
	data.MessageMeta = cqrs.FromContext(ctx).NewMessage(EditedEventType, aggregateId...).Meta()
	return &data
}

func (e *EditedEvent) Meta() *cqrs.MessageMeta {
	return e.MessageMeta
}

func NewCompletedEventMessage(ctx context.Context, data CompletedEvent, aggregateId ...string) *CompletedEvent {
	data.MessageMeta = cqrs.FromContext(ctx).NewMessage(CompletedEventType, aggregateId...).Meta()
	return &data
}

func (e *CompletedEvent) Meta() *cqrs.MessageMeta {
	return e.MessageMeta
}
