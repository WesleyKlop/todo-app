package todos

import "github.com/google/uuid"

type RawTodo struct {
	Content string `json:"content" binding:"required"`
}

type Todo struct {
	Content string `json:"content" binding:"required"`
	Id      string `json:"id" binding:"required"`
	Done    bool   `json:"done" binding:"required"`
}

type PatchTodo struct {
	Content *string `json:"content"`
	Done    *bool   `json:"done"`
}

func (r *RawTodo) Create() Todo {
	return Todo{
		Id:      uuid.New().String(),
		Content: r.Content,
		Done:    false,
	}
}

func (r *Todo) Patch(patch PatchTodo) {
	if patch.Content != nil && *patch.Content != r.Content {
		r.Content = *patch.Content
	}
	if patch.Done != nil {
		r.Done = *patch.Done
	}
}
