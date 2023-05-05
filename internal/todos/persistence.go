package todos

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func LoadFromFile(ctx context.Context, path string) (*TodoRepository, error) {
	tr := otel.Tracer("persistence")
	_, span := tr.Start(ctx, "load-from-file")
	defer span.End()
	if _, err := os.Stat(path); err != nil {
		span.SetAttributes(attribute.Bool("new", true))
		return NewTodoRepository(), nil
	}

	handle, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	span.AddEvent("open")
	defer handle.Close()

	content, err := io.ReadAll(handle)
	if err != nil {
		return nil, err
	}
	span.AddEvent("read")

	if len(content) < 2 {
		span.SetAttributes(attribute.Bool("new", true))
		return NewTodoRepository(), nil
	}

	var store []Todo
	err = json.Unmarshal(content, &store)
	if err != nil {
		return nil, err
	}
	span.AddEvent("unmarshal")
	return &TodoRepository{store}, nil
}

func PersistToFile(ctx context.Context, repo *TodoRepository, path string) error {
	tr := otel.Tracer("persistence")
	_, span := tr.Start(ctx, "save-to-file")
	defer span.End()
	if repo.store == nil {
		return fmt.Errorf("store is nil. Not saving")
	}
	handle, err := os.Create(path)
	if err != nil {
		return err
	}
	span.AddEvent("create")
	defer handle.Close()

	content, err := json.Marshal(repo.store)
	if err != nil {
		return err
	}
	span.AddEvent("marshal")
	handle.Write(content)
	span.AddEvent("write")
	return nil
}
