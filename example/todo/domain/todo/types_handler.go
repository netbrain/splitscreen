package todo

import (
	"context"
	"strings"

	"github.com/netbrain/splitscreen/cqrs"
)

func (a Aggregate) titleAndContent(d string) (string, string) {
	data := strings.Split(d, "\n")
	var title, content string
	title = data[0]
	if len(data) > 1 {
		content = data[1]
	}
	return title, content
}

func (a Aggregate) HandleCreateCommand(ctx context.Context, cmd CreateCommand) (cqrs.Message, error) {
	title, content := a.titleAndContent(cmd.Data)
	return NewCreatedEventMessage(ctx, CreatedEvent{
		Title:   title,
		Content: content,
	}, cmd.AggregateID), nil
}

func (a Aggregate) HandleEditCommand(ctx context.Context, cmd EditCommand) (cqrs.Message, error) {
	title, content := a.titleAndContent(cmd.Data)
	return NewEditedEventMessage(ctx, EditedEvent{
		Title:   title,
		Content: content,
	}, cmd.AggregateID), nil
}

func (a Aggregate) HandleCompleteCommand(ctx context.Context, cmd CompleteCommand) (cqrs.Message, error) {
	return NewCompletedEventMessage(ctx, CompletedEvent{}, cmd.AggregateID), nil
}

func (a *Aggregate) ApplyCreatedEvent(ctx context.Context, event CreatedEvent) error {
	a.Title = event.Title
	a.Content = event.Content
	return nil
}

func (a *Aggregate) ApplyEditedEvent(ctx context.Context, event EditedEvent) error {
	a.Title = event.Title
	a.Content = event.Content
	return nil
}

func (a *Aggregate) ApplyCompletedEvent(ctx context.Context, event CompletedEvent) error {
	a.Done = true
	return nil
}
