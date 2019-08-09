//go:generate sh -c "SSPATH=$PWD/../../../../cmd/splitscreen go run ../../../../cmd/splitscreen/main.go -generate handler"

// ALL OTHER MORTALS SHOULD USE go:generate splitscreen -generate handler

package todo

import "github.com/netbrain/splitscreen/cqrs"

type CreateCommand struct {
	*cqrs.MessageMeta
	Data string `json:"data"`
}

type CreatedEvent struct {
	*cqrs.MessageMeta
	Title   string `json:"title"`
	Content string `json:"content"`
}

type EditCommand struct {
	*cqrs.MessageMeta
	Data string `json:"data"`
}

type EditedEvent struct {
	*cqrs.MessageMeta
	Title   string `json:"title"`
	Content string `json:"content"`
}

type CompleteCommand struct {
	*cqrs.MessageMeta
}

type CompletedEvent struct {
	*cqrs.MessageMeta
}

type Aggregate struct {
	*cqrs.AggregateMeta
	Title   string `json:"title"`
	Content string `json:"content"`
	Done    bool   `json:"done"`
}
