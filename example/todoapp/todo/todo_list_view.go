package todo

//go:generate ss-listener
type TodoListView struct {
	Todos []*Todo
	index map[string]int
}

type Todo struct {
	ID string `json:"id"`
	Text string `json:"text"`
	Version int `json:"version"`
}

func NewTodoListView() *TodoListView{
	return &TodoListView{
		index: make(map[string]int),
	}
}

func (t *TodoListView) OnTodoItemCreatedEvent(event *TodoItemCreatedEvent){
	t.Todos = append(t.Todos,&Todo{
		ID:   event.Aggregate.ID,
		Text: event.Text,
		Version:event.Aggregate.Version,
	})
	t.index[event.Aggregate.ID] = len(t.Todos)-1
}

func (t *TodoListView) OnTodoItemArchivedEvent(event *TodoItemArchivedEvent){
	todo := t.Todos[t.index[event.Aggregate.ID]]
	todo.Text += " (archived)"
	todo.Version = event.Aggregate.Version
}

func (t *TodoListView) OnTodoItemDeletedEvent(event *TodoItemDeletedEvent){
	t.Todos = append(t.Todos[:t.index[event.Aggregate.ID]],t.Todos[t.index[event.Aggregate.ID]+1:]...)
	delete(t.index,event.Aggregate.ID)
}
