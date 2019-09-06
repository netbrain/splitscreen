//go:generate sh -c "SSPATH=$PWD/../../../cmd/splitscreen go run ../../../cmd/splitscreen/main.go -generate view"

// ALL OTHER MORTALS SHOULD USE go:generate splitscreen -generate view

package view

import (
	"context"
	"fmt"
	"github.com/netbrain/splitscreen/example/todo/domain/todo"
)

type Todo struct {
	Title   string
	Content string
	Done    bool
}

type TodoView struct {
	todos map[string]*Todo
}

func NewTodoView() *TodoView {
	return &TodoView{
		todos: map[string]*Todo{},
	}
}

func (v *TodoView) All() (todos []*Todo) {
	for _, todo := range v.todos {
		todos = append(todos, todo)
	}
	return
}

func (v *TodoView) OnCreatedEvent(ctx context.Context, event *todo.CreatedEvent) error {
	v.todos[event.AggregateID] = &Todo{
		Title:   event.Title,
		Content: event.Content,
	}
	return nil
}

func (v *TodoView) OnEditedEvent(ctx context.Context, event *todo.EditedEvent) error {
	todo, ok := v.todos[event.AggregateID]
	if !ok {
		return fmt.Errorf("no such todo")
	}
	todo.Content = event.Content
	todo.Title = event.Title
	return nil
}

func (v *TodoView) OnCompletedEvent(ctx context.Context, event *todo.CompletedEvent) error {
	todo, ok := v.todos[event.AggregateID]
	if !ok {
		return fmt.Errorf("no such todo")
	}
	todo.Done = true
	return nil
}
