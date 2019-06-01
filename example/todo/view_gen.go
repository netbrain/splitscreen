// CODE GENERATED, DO NOT MODIFY

package todo

import (
	"context"
	"github.com/netbrain/splitscreen/cqrs"
)

const TodoViewType cqrs.ViewType = "TodoView"

func (v *TodoView) Register(app *cqrs.App) {
	app.RegisterView(TodoViewType, v)

	app.Subscribe(func(ctx context.Context, msg cqrs.Message) error {
		return v.OnTodoCreatedEvent(ctx, msg.(*TodoCreatedEvent))
	}, TodoCreatedEventType)

	app.Subscribe(func(ctx context.Context, msg cqrs.Message) error {
		return v.OnTodoEditedEvent(ctx, msg.(*TodoEditedEvent))
	}, TodoEditedEventType)

	app.Subscribe(func(ctx context.Context, msg cqrs.Message) error {
		return v.OnTodoCompletedEvent(ctx, msg.(*TodoCompletedEvent))
	}, TodoCompletedEventType)

}
