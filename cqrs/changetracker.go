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
	for i := 0; i < len(c.changes); i++ {
		if err := c.app.Manage(ctx, c.changes[i]); err != nil {
			return err
		}
	}
	if err := c.app.Store(ctx, c.changes...); err != nil {
		return err
	}
	if err := c.app.Emit(ctx, c.changes...); err != nil {
		return err
	}
	c.changes = nil
	return nil
}
