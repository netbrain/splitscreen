package main

import (
	"context"
	"fmt"
	"github.com/netbrain/splitscreen/cqrs"
	"github.com/netbrain/splitscreen/example/todo/domain/todo"
	"github.com/netbrain/splitscreen/example/todo/view"
)

func main() {
	// bootstrap application
	app := cqrs.New(nil)

	// views
	todos := view.NewTodoView()

	// register aggregates & views
	for _, r := range []cqrs.Registerable{
		&todo.Aggregate{},
		todos,
	} {
		r.Register(app)
	}

	// subscribe to all events and print them out to console
	app.Subscribe(func(ctx context.Context, msg cqrs.Message) error {
		data, err := app.Serialize(msg)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("%s\n", data)
		}
		return nil
	})

	// create a new context
	ctx := app.NewContext(context.Background())

	// dispatch command message
	err := app.DispatchMessage(ctx, todo.NewCreateCommandMessage(ctx, todo.CreateCommand{
		Data: "Hello\nWorld",
	}))
	if err != nil {
		fmt.Println(err)
		return
	}

	// commit changes
	ct := cqrs.ChangeTrackerFromContext(ctx)
	if err := ct.CommitChanges(ctx); err != nil { // in a http setting you would use app.Middleware()
		fmt.Println(err)
	}

	// print view
	for _, t := range todos.All() {
		fmt.Printf("\n-- %s --\n%s\ncompleted=%v", t.Title, t.Content, t.Done)
	}

}
