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

func (r *TodoRepository) getIdx(id string) int {
	for idx, todo := range r.store {
		if todo.Id == id {
			return idx
		}
	}
	return -1
}

func (r *TodoRepository) Get(id string) *Todo {
	idx := r.getIdx(id)
	if idx >= 0 {
		return &r.store[idx]
	}
	return nil
}

func (r *TodoRepository) Remove(id string) {
	idx := r.getIdx(id)
	if idx >= 0 {
		r.store = remove(r.store, idx)
	}
}

func (r *TodoRepository) Update(todo Todo) bool {
	idx := r.getIdx(todo.Id)
	if idx >= 0 {
		r.store[idx] = todo
		return true
	}
	return false
}

func (r *TodoRepository) Exists(id string) bool {
	return r.getIdx(id) >= 0
}

func (r *TodoRepository) Clear() {
	r.store = make([]Todo, 0)
}

func remove[T any](s []T, i int) []T {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
