//go:generate sh -c "SSPATH=$PWD/../../cmd/splitscreen go run ../../cmd/splitscreen/main.go -generate handler"

// ALL OTHER MORTALS SHOULD USE go:generate splitscreen -generate handler

package todo

import "github.com/netbrain/splitscreen/cqrs"

type CreateTodoCommand struct {
	*cqrs.MessageMeta
	Data string `json:"data"`
}

type TodoCreatedEvent struct {
	*cqrs.MessageMeta
	Title   string `json:"title"`
	Content string `json:"content"`
}

type EditTodoCommand struct {
	*cqrs.MessageMeta
	Data string `json:"data"`
}

type TodoEditedEvent struct {
	*cqrs.MessageMeta
	Title   string `json:"title"`
	Content string `json:"content"`
}

type CompleteTodoCommand struct {
	*cqrs.MessageMeta
}

type TodoCompletedEvent struct {
	*cqrs.MessageMeta
}

type TodoAggregate struct {
	*cqrs.AggregateMeta
	Title   string `json:"title"`
	Content string `json:"content"`
	Done    bool   `json:"done"`
}
