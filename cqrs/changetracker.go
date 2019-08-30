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
	for {
		changes := c.changes
		if changes == nil {
			return nil
		}
		c.changes = nil
		if err := c.app.Store(ctx, changes...); err != nil {
			return err
		}
		if err := c.app.Emit(ctx, changes...); err != nil {
			return err
		}
	}
}
