// CODE GENERATED, DO NOT MODIFY

package view

import (
	context "context"

	"github.com/netbrain/splitscreen/cqrs"
	todo "github.com/netbrain/splitscreen/example/todo/domain/todo"
)

const TodoViewType cqrs.ViewType = "TodoView"

func (v *TodoView) Register(app *cqrs.App) {
	app.RegisterView(TodoViewType, v)

	app.Subscribe(func(ctx context.Context, msg cqrs.Message) error {
		return v.OnCreatedEvent(ctx, msg.(*todo.CreatedEvent))
	}, cqrs.ViewSubscription, todo.CreatedEventType)

	app.Subscribe(func(ctx context.Context, msg cqrs.Message) error {
		return v.OnEditedEvent(ctx, msg.(*todo.EditedEvent))
	}, cqrs.ViewSubscription, todo.EditedEventType)

	app.Subscribe(func(ctx context.Context, msg cqrs.Message) error {
		return v.OnCompletedEvent(ctx, msg.(*todo.CompletedEvent))
	}, cqrs.ViewSubscription, todo.CompletedEventType)

}
