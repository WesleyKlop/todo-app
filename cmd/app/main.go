package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/wesleyklop/todo-api/v2/internal/todos"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

func readOrCreateExistingTodos() (list *[]todos.Todo, err error) {
	const path = "/mnt/data/todos.json"

	if _, err = os.Stat(path); err != nil {
		err = os.WriteFile(path, []byte("[]"), 0644)
		return
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
	if len(content) > 2 {
		err = json.Unmarshal(content, &list)
	} else {
		tmp := make([]todos.Todo, 0)
		list = &tmp
	}

	return
}

func saveExistingTodos(db *[]todos.Todo) error {
	const path = "/mnt/data/todos.json"

	handle, err := os.Create(path)
	if err != nil {
		return err
	}
	defer handle.Close()

	content, err := json.Marshal(*db)
	if err != nil {
		return err
	}
	handle.Write(content)
	return nil
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	router := gin.New()

	log, _ := zap.NewDevelopment()
	defer log.Sync()

	router.Use(ginzap.Ginzap(log, time.RFC3339, false))
	router.Use(ginzap.RecoveryWithZap(log, true))

	repository, err := todos.LoadFromFile("/mnt/data/todos.json")
	if err != nil {
		log.Fatal("Failed to create todo repository", zap.Error(err))
	}

	todos.NewTodoRouter(router.Group("/api/todos"), repository)

	service := http.Server{
		ReadHeaderTimeout: time.Second * 10,
		Addr:              ":8080",
		Handler:           router,
	}

	go func() {
		if err := service.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen: %s\n", zap.Error(err))
		}
	}()

	<-ctx.Done()

	stop()

	log.Info("shutting down gracefully... saving todos...")
	todos.PersistToFile(repository, "/mnt/data/todos.json")

	shutdownTime := 5 * time.Second
	if mode, exists := os.LookupEnv("MODE"); exists && mode == "development" {
		log.Info("skipping graceful shutdown")
		shutdownTime = 1 * time.Second
	}

	// Give Gin 5 seconds to handle inflight requests (Cloud Run gives us 10 before SIGKILL)
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTime)
	defer cancel()
	if err := service.Shutdown(ctx); err != nil {
		log.Fatal("shutdown deadline exceeded.")
	}

	log.Info("exiting")
}
