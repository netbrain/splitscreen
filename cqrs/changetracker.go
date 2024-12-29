package cqrs

import (
	"context"
	"fmt"
)

type ChangeTrackerFactory interface {
	NewChangeTracker() ChangeTracker
}

type DefaultChangeTrackerFactory struct {}

func NewDefaultChangeTrackerFactory() ChangeTrackerFactory {
	return &DefaultChangeTrackerFactory{}
}

func (d *DefaultChangeTrackerFactory) NewChangeTracker() ChangeTracker {
	return NewDefaultChangeTracker()
}

type ChangeTracker interface {
	TrackChange(event Message) error
	CommitChanges(ctx context.Context) error
	Changes() []Message
}

type DefaultChangeTracker struct {
	changes []Message
}

func NewDefaultChangeTracker() *DefaultChangeTracker {
	return &DefaultChangeTracker{}
}

func (c *DefaultChangeTracker) TrackChange(event Message) error {
	if !event.Meta().MessageType.IsEvent() {
		return fmt.Errorf("expected event")
	}
	c.changes = append(c.changes, event)
	return nil
}

func (c *DefaultChangeTracker) CommitChanges(ctx context.Context) error {
	app := FromContext(ctx)
	for i := 0; i < len(c.changes); i++ {
		if err := app.Manage(ctx, c.changes[i]); err != nil {
			return err
		}
	}
	//TODO Store & Emit should be in a single transaction
	if err := app.Store(ctx, c.changes...); err != nil {
		return err
	}
	if err := app.Emit(ctx, c.changes...); err != nil {
		return err
	}
	c.changes = nil
	return nil
}

func (c *DefaultChangeTracker) Changes() []Message {
	return append([]Message{},c.changes...)
}