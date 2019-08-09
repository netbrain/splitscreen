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
		return v.OnEditedEvent(ctx, msg.(*todo.EditedEvent))
	}, todo.EditedEventType)

}
