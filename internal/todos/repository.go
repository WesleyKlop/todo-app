package todos

type TodoRepository struct {
	store []Todo
}

func NewTodoRepository() *TodoRepository {
	return &TodoRepository{
		store: make([]Todo, 0),
	}
}

func (r *TodoRepository) Create(raw *RawTodo) *Todo {
	todo := raw.Create()

	r.store = append(r.store, todo)
	return &todo
}

func (r *TodoRepository) List() *[]Todo {
	return &r.store
}

func (r *TodoRepository) Get(id string) *Todo {
	for _, todo := range r.store {
		if todo.Id == id {
			return &todo
		}
	}
	return nil
}

func (r *TodoRepository) Clear() {
	r.store = make([]Todo, 0)
}
