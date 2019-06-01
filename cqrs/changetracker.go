package cqrs

import (
	"context"
	"fmt"
)

type ChangeTracker struct {
	changes []Message
	app     *App
}

func NewChangeTracker(app *App) *ChangeTracker {
	return &ChangeTracker{
		app: app,
	}
}

func (c *ChangeTracker) TrackChange(event Message) error {
	if !event.Meta().MessageType.IsEvent() {
		return fmt.Errorf("expected event")
	}

	c.changes = append(c.changes, event)
	return nil
}

func (c *ChangeTracker) CommitChanges(ctx context.Context) error {
	err := c.app.Store(ctx, c.changes...)
	if err != nil {
		return err
	}
	for _, msg := range c.changes {
		if err := c.app.Emit(ctx, msg); err != nil {
			return err
		}
	}
	c.changes = nil
	return nil
}
