package main

import (
	"fmt"
	"github.com/netbrain/splitscreen/cqrs"
	"github.com/netbrain/splitscreen/example/todoapp/todo"
)


func main(){
	cqrs.RegisterEventListener(eventLogger, "")
	cqrs.RegisterCommandListener(commandLogger, "")

	//Register views
	cqrs.RegisterView(todo.TodoListViewType,todo.NewTodoListView())

	id := cqrs.IDFunc()
	todo.CreateTodoCommand{Text: "New todo"}.Dispatch(id,0)

	todo.ArchiveTodoCommand{}.Dispatch(id, 1)

	todo.DeleteTodoCommand{}.Dispatch(id,2)

}

func eventLogger(e cqrs.Event) {
	fmt.Printf("EventID: %s\nCausationID: %s\nEventType: %s\nAggregateID: %s\nAggregateType: %s\nAggregateVersion: %d\nData: %#v\n\n",
		e.ID,
		e.CausationID,
		e.Type,
		e.Aggregate.ID,
		e.Aggregate.Type,
		e.Aggregate.Version,
		e.Impl,
	)
}

func commandLogger(c cqrs.Command) {
	fmt.Printf("CommandID: %s\nCausationID: %s\nCommandType: %s\nAggregateID: %s\nAggregateType: %s\nAggregateVersion: %d\nData: %#v\n\n",
		c.ID,
		c.CausationID,
		c.Type,
		c.Aggregate.ID,
		c.Aggregate.Type,
		c.Aggregate.Version,
		c.Impl,
	)
}