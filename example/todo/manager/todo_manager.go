//go:generate sh -c "SSPATH=$PWD/../../../cmd/splitscreen go run ../../../cmd/splitscreen/main.go -generate manager"

// ALL OTHER MORTALS SHOULD USE go:generate splitscreen -generate manager

package manager

import (
	"context"
	"github.com/netbrain/splitscreen/example/todo/domain/todo"
)

type TodoManager struct {

}

func NewTodoManager() *TodoManager {
	return &TodoManager{

	}
}


func (v *TodoManager) OnCreatedEvent(ctx context.Context, event *todo.CreatedEvent) error {
	// TODO's are too simple to have process managers, but still, leaving as example
	// Managers are basically the same as views, except they dont output view model changes. they output new commands
	return nil
}
