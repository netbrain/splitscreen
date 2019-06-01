// CODE GENERATED, DO NOT MODIFY

package todo

import (
	"context"
	"fmt"

	"github.com/netbrain/splitscreen/cqrs"
)

const (
	TodoAggregateType cqrs.AggregateType = "TodoAggregate"

	CreateTodoCommandType cqrs.MessageType = "CreateTodoCommand"

	EditTodoCommandType cqrs.MessageType = "EditTodoCommand"

	CompleteTodoCommandType cqrs.MessageType = "CompleteTodoCommand"

	TodoCreatedEventType cqrs.MessageType = "TodoCreatedEvent"

	TodoEditedEventType cqrs.MessageType = "TodoEditedEvent"

	TodoCompletedEventType cqrs.MessageType = "TodoCompletedEvent"
)

func (a *TodoAggregate) Register(app *cqrs.App) {
	app.RegisterAggregate(TodoAggregateType, func() cqrs.AggregateRoot {
		return &TodoAggregate{AggregateMeta: &cqrs.AggregateMeta{}}
	})

	app.RegisterMessage(func() cqrs.Message {
		return &TodoCreatedEvent{
			MessageMeta: &cqrs.MessageMeta{
				AggregateMeta: &cqrs.AggregateMeta{
					AggregateType: TodoAggregateType,
				},
				MessageType: TodoCreatedEventType,
			},
		}
	})

	app.RegisterMessage(func() cqrs.Message {
		return &TodoEditedEvent{
			MessageMeta: &cqrs.MessageMeta{
				AggregateMeta: &cqrs.AggregateMeta{
					AggregateType: TodoAggregateType,
				},
				MessageType: TodoEditedEventType,
			},
		}
	})

	app.RegisterMessage(func() cqrs.Message {
		return &TodoCompletedEvent{
			MessageMeta: &cqrs.MessageMeta{
				AggregateMeta: &cqrs.AggregateMeta{
					AggregateType: TodoAggregateType,
				},
				MessageType: TodoCompletedEventType,
			},
		}
	})

	app.RegisterMessage(func() cqrs.Message {
		return &CreateTodoCommand{
			MessageMeta: &cqrs.MessageMeta{
				AggregateMeta: &cqrs.AggregateMeta{
					AggregateType: TodoAggregateType,
				},
				MessageType: CreateTodoCommandType,
			},
		}
	})

	app.RegisterMessage(func() cqrs.Message {
		return &EditTodoCommand{
			MessageMeta: &cqrs.MessageMeta{
				AggregateMeta: &cqrs.AggregateMeta{
					AggregateType: TodoAggregateType,
				},
				MessageType: EditTodoCommandType,
			},
		}
	})

	app.RegisterMessage(func() cqrs.Message {
		return &CompleteTodoCommand{
			MessageMeta: &cqrs.MessageMeta{
				AggregateMeta: &cqrs.AggregateMeta{
					AggregateType: TodoAggregateType,
				},
				MessageType: CompleteTodoCommandType,
			},
		}
	})

}

func (a *TodoAggregate) Meta() *cqrs.AggregateMeta {
	return a.AggregateMeta
}

func (a *TodoAggregate) Handle(ctx context.Context, msg cqrs.Message) (err error) {
	var event cqrs.Message
	switch msg.Meta().MessageType {

	case CreateTodoCommandType:
		event, err = a.HandleCreateTodoCommand(ctx, *msg.(*CreateTodoCommand))

	case EditTodoCommandType:
		event, err = a.HandleEditTodoCommand(ctx, *msg.(*EditTodoCommand))

	case CompleteTodoCommandType:
		event, err = a.HandleCompleteTodoCommand(ctx, *msg.(*CompleteTodoCommand))

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

func (a *TodoAggregate) Apply(ctx context.Context, msg cqrs.Message) (err error) {
	switch msg.Meta().MessageType {

	case TodoCreatedEventType:
		err = a.ApplyTodoCreatedEvent(ctx, *msg.(*TodoCreatedEvent))

	case TodoEditedEventType:
		err = a.ApplyTodoEditedEvent(ctx, *msg.(*TodoEditedEvent))

	case TodoCompletedEventType:
		err = a.ApplyTodoCompletedEvent(ctx, *msg.(*TodoCompletedEvent))

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

func NewCreateTodoCommandMessage(ctx context.Context, data CreateTodoCommand, aggregateId ...string) *CreateTodoCommand {
	data.MessageMeta = cqrs.FromContext(ctx).NewMessage(CreateTodoCommandType, aggregateId...).Meta()
	return &data
}

func (e *CreateTodoCommand) Meta() *cqrs.MessageMeta {
	return e.MessageMeta
}

func NewEditTodoCommandMessage(ctx context.Context, data EditTodoCommand, aggregateId ...string) *EditTodoCommand {
	data.MessageMeta = cqrs.FromContext(ctx).NewMessage(EditTodoCommandType, aggregateId...).Meta()
	return &data
}

func (e *EditTodoCommand) Meta() *cqrs.MessageMeta {
	return e.MessageMeta
}

func NewCompleteTodoCommandMessage(ctx context.Context, data CompleteTodoCommand, aggregateId ...string) *CompleteTodoCommand {
	data.MessageMeta = cqrs.FromContext(ctx).NewMessage(CompleteTodoCommandType, aggregateId...).Meta()
	return &data
}

func (e *CompleteTodoCommand) Meta() *cqrs.MessageMeta {
	return e.MessageMeta
}

func NewTodoCreatedEventMessage(ctx context.Context, data TodoCreatedEvent, aggregateId ...string) *TodoCreatedEvent {
	data.MessageMeta = cqrs.FromContext(ctx).NewMessage(TodoCreatedEventType, aggregateId...).Meta()
	return &data
}

func (e *TodoCreatedEvent) Meta() *cqrs.MessageMeta {
	return e.MessageMeta
}

func NewTodoEditedEventMessage(ctx context.Context, data TodoEditedEvent, aggregateId ...string) *TodoEditedEvent {
	data.MessageMeta = cqrs.FromContext(ctx).NewMessage(TodoEditedEventType, aggregateId...).Meta()
	return &data
}

func (e *TodoEditedEvent) Meta() *cqrs.MessageMeta {
	return e.MessageMeta
}

func NewTodoCompletedEventMessage(ctx context.Context, data TodoCompletedEvent, aggregateId ...string) *TodoCompletedEvent {
	data.MessageMeta = cqrs.FromContext(ctx).NewMessage(TodoCompletedEventType, aggregateId...).Meta()
	return &data
}

func (e *TodoCompletedEvent) Meta() *cqrs.MessageMeta {
	return e.MessageMeta
}
