package todo

import (
	"context"
	"strings"

	"github.com/netbrain/splitscreen/cqrs"
)

func (a TodoAggregate) titleAndContent(d string) (string, string) {
	data := strings.Split(d, "\n")
	var title, content string
	title = data[0]
	if len(data) > 1 {
		content = data[1]
	}
	return title, content
}

func (a TodoAggregate) HandleCreateTodoCommand(ctx context.Context, cmd CreateTodoCommand) (cqrs.Message, error) {
	title, content := a.titleAndContent(cmd.Data)
	return NewTodoCreatedEventMessage(ctx, TodoCreatedEvent{
		Title:   title,
		Content: content,
	}, cmd.AggregateID), nil
}

func (a TodoAggregate) HandleEditTodoCommand(ctx context.Context, cmd EditTodoCommand) (cqrs.Message, error) {
	title, content := a.titleAndContent(cmd.Data)
	return NewTodoEditedEventMessage(ctx, TodoEditedEvent{
		Title:   title,
		Content: content,
	}, cmd.AggregateID), nil
}

func (a TodoAggregate) HandleCompleteTodoCommand(ctx context.Context, cmd CompleteTodoCommand) (cqrs.Message, error) {
	return NewTodoCompletedEventMessage(ctx, TodoCompletedEvent{}, cmd.AggregateID), nil
}

func (a *TodoAggregate) ApplyTodoCreatedEvent(ctx context.Context, event TodoCreatedEvent) error {
	a.Title = event.Title
	a.Content = event.Content
	return nil
}

func (a *TodoAggregate) ApplyTodoEditedEvent(ctx context.Context, event TodoEditedEvent) error {
	a.Title = event.Title
	a.Content = event.Content
	return nil
}

func (a *TodoAggregate) ApplyTodoCompletedEvent(ctx context.Context, event TodoCompletedEvent) error {
	a.Done = true
	return nil
}
