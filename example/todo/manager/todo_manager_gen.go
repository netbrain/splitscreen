// CODE GENERATED, DO NOT MODIFY

package manager

import (
	context "context"

	"github.com/netbrain/splitscreen/cqrs"
	todo "github.com/netbrain/splitscreen/example/todo/domain/todo"
)

const TodoManagerType cqrs.ManagerType = "TodoManager"

func (v *TodoManager) Register(app *cqrs.App) {
	app.RegisterManager(TodoManagerType, v)

	app.Subscribe(func(ctx context.Context, msg cqrs.Message) error {
		return v.OnCreatedEvent(ctx, msg.(*todo.CreatedEvent))
	}, cqrs.ManagerSubscription, todo.CreatedEventType)

}
