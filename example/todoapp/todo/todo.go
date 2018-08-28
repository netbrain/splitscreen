package todo

//go:generate ss-aggregate

import (
	"github.com/netbrain/splitscreen/cqrs"
	"fmt"
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


func (t *TodoAggregate) HandleCreateTodoCommand(c *CreateTodoCommand) error {
	if t.Version > 0 {
		return fmt.Errorf("todo already created")
	}
	if len(c.Text) == 0 {
		return fmt.Errorf("invalid todo name %s", c.Text)
	}

	return TodoItemCreatedEvent{
		Text:  c.Text,
	}.Apply(t.Aggregate,c.Command)
}

func (t *TodoAggregate) HandleArchiveTodoCommand(c *ArchiveTodoCommand) error {
	if t.Archived {
		return fmt.Errorf("todo already archived")
	}
	return TodoItemArchivedEvent{}.Apply(t.Aggregate,c.Command)
}

func (t *TodoAggregate) HandleDeleteTodoCommand(c *DeleteTodoCommand) error {
	return TodoItemDeletedEvent{}.Apply(t.Aggregate,c.Command)
}

func (t *TodoAggregate) ApplyTodoItemCreatedEvent(e *TodoItemCreatedEvent) error {
	t.Text = e.Text
	return nil
}

func (t *TodoAggregate) ApplyTodoItemArchivedEvent(e *TodoItemArchivedEvent) error {
	t.Archived = true
	return nil
}

func (t *TodoAggregate) ApplyTodoItemDeletedEvent(e *TodoItemDeletedEvent) error {
	return nil
}
