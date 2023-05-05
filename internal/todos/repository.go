package todos

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type TodoRepository struct {
	store []Todo
}

func NewTodoRepository() *TodoRepository {
	return &TodoRepository{
		store: make([]Todo, 0),
	}
}

func getSpan(ctx context.Context, action string) trace.Span {
	tr := otel.Tracer("repository")
	_, span := tr.Start(ctx, action)
	return span
}

func (r *TodoRepository) Create(ctx context.Context, raw *RawTodo) *Todo {
	span := getSpan(ctx, "create")
	defer span.End()
	todo := raw.Create()

	r.store = append(r.store, todo)
	return &todo
}

func (r *TodoRepository) List(ctx context.Context) *[]Todo {
	span := getSpan(ctx, "list")
	defer span.End()
	return &r.store
}

func (r *TodoRepository) getIdx(ctx context.Context, id string) int {
	span := getSpan(ctx, "find")
	defer span.End()
	for idx, todo := range r.store {
		if todo.Id == id {
			return idx
		}
	}
	return -1
}

func (r *TodoRepository) Get(ctx context.Context, id string) *Todo {
	span := getSpan(ctx, "get")
	defer span.End()
	idx := r.getIdx(ctx, id)
	if idx >= 0 {
		return &r.store[idx]
	}
	return nil
}

func (r *TodoRepository) Remove(ctx context.Context, id string) {
	span := getSpan(ctx, "remove")
	defer span.End()
	idx := r.getIdx(ctx, id)
	if idx >= 0 {
		r.store = remove(r.store, idx)
	}
}

func (r *TodoRepository) Update(ctx context.Context, todo Todo) bool {
	span := getSpan(ctx, "update")
	defer span.End()
	idx := r.getIdx(ctx, todo.Id)
	if idx >= 0 {
		r.store[idx] = todo
		return true
	}
	return false
}

func (r *TodoRepository) Exists(ctx context.Context, id string) bool {
	span := getSpan(ctx, "exists")
	defer span.End()
	return r.getIdx(ctx, id) >= 0
}

func (r *TodoRepository) Clear(ctx context.Context) {
	span := getSpan(ctx, "clear")
	defer span.End()
	r.store = make([]Todo, 0)
}

func remove[T any](s []T, i int) []T {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
