package todos

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func LoadFromFile(path string) (*TodoRepository, error) {
	if _, err := os.Stat(path); err != nil {
		return NewTodoRepository(), nil
	}

	handle, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer handle.Close()

	content, err := io.ReadAll(handle)
	if err != nil {
		return nil, err
	}

	if len(content) < 2 {
		return NewTodoRepository(), nil
	}

	var store []Todo
	err = json.Unmarshal(content, &store)
	if err != nil {
		return nil, err
	}
	return &TodoRepository{store}, nil
}

func PersistToFile(repo *TodoRepository, path string) error {
	if repo.store == nil {
		return fmt.Errorf("store is nil. Not saving")
	}
	handle, err := os.Create(path)
	if err != nil {
		return err
	}
	defer handle.Close()

	content, err := json.Marshal(repo.store)
	if err != nil {
		return err
	}
	handle.Write(content)
	return nil
}
