//go:generate splitscreen -generate view

package todo

import (
	"context"
	"fmt"
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

func (v *TodoView) Todos() (todos []*Todo) {
	for _, todo := range v.todos {
		todos = append(todos, todo)
	}
	return
}

func (v *TodoView) OnTodoCreatedEvent(ctx context.Context, event *TodoCreatedEvent) error {
	v.todos[event.AggregateID] = &Todo{
		Title:   event.Title,
		Content: event.Content,
	}
	return nil
}

func (v *TodoView) OnTodoEditedEvent(ctx context.Context, event *TodoEditedEvent) error {
	todo, ok := v.todos[event.AggregateID]
	if !ok {
		return fmt.Errorf("no such todo")
	}
	todo.Content = event.Content
	todo.Title = event.Title
	return nil
}

func (v *TodoView) OnTodoCompletedEvent(ctx context.Context, event *TodoCompletedEvent) error {
	todo, ok := v.todos[event.AggregateID]
	if !ok {
		return fmt.Errorf("no such todo")
	}
	todo.Done = true
	return nil
}
