package todo

//go:generate ss-aggregate

import (
	"github.com/netbrain/splitscreen/cqrs"
	"fmt"
	"context"
)
type CreateTodoCommand struct {
	*cqrs.Command
	Text string
}

type ArchiveTodoCommand struct {
	*cqrs.Command
}

type DeleteTodoCommand struct {
	*cqrs.Command
}

type TodoItemCreatedEvent struct {
	*cqrs.Event
	Text string
}

type TodoItemArchivedEvent struct {
	*cqrs.Event
}

type TodoItemDeletedEvent struct {
	*cqrs.Event
}

type TodoAggregate struct {
	*cqrs.Aggregate
	Text    string
	Archived bool
}


func (t *TodoAggregate) HandleCreateTodoCommand(ctx context.Context, c *CreateTodoCommand) error {
	if t.Version > 0 {
		return fmt.Errorf("todo already created")
	}
	if len(c.Text) == 0 {
		return fmt.Errorf("invalid todo name %s", c.Text)
	}

	return TodoItemCreatedEvent{
		Text:  c.Text,
	}.Apply(ctx,t.Aggregate,c.Command)
}

func (t *TodoAggregate) HandleArchiveTodoCommand(ctx context.Context,c *ArchiveTodoCommand) error {
	if t.Archived {
		return fmt.Errorf("todo already archived")
	}
	return TodoItemArchivedEvent{}.Apply(ctx,t.Aggregate,c.Command)
}

func (t *TodoAggregate) HandleDeleteTodoCommand(ctx context.Context,c *DeleteTodoCommand) error {
	return TodoItemDeletedEvent{}.Apply(ctx,t.Aggregate,c.Command)
}

func (t *TodoAggregate) ApplyTodoItemCreatedEvent(ctx context.Context,e *TodoItemCreatedEvent) error {
	t.Text = e.Text
	return nil
}

func (t *TodoAggregate) ApplyTodoItemArchivedEvent(ctx context.Context,e *TodoItemArchivedEvent) error {
	t.Archived = true
	return nil
}

func (t *TodoAggregate) ApplyTodoItemDeletedEvent(ctx context.Context,e *TodoItemDeletedEvent) error {
	return nil
}
