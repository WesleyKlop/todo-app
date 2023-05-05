package todos

import "github.com/google/uuid"

type RawTodo struct {
	Content string `json:"content" binding:"required"`
}

type Todo struct {
	Content string `json:"content" binding:"required"`
	Id      string `json:"id" binding:"required"`
}

func (r *RawTodo) Create() Todo {
	return Todo{
		Content: r.Content,
		Id:      uuid.New().String(),
	}
}
