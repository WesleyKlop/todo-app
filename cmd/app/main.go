package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/WesleyKlop/todo-api/v2/internal/todos"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
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

	todoList, err := readOrCreateExistingTodos()
	if err != nil {
		log.Fatal("Failed to get or create todo db", zap.Error(err))
	}

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "bas")
	})
	router.GET("/api/todos", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, todoList)
	})
	router.GET("/api/todos/:todo", func(ctx *gin.Context) {
		todoId := ctx.Param("todo")
		for _, todo := range *todoList {
			if todo.Id == todoId {
				ctx.JSON(http.StatusOK, todo)
				return
			}
		}
		ctx.AbortWithStatus(http.StatusNotFound)
	})
	router.POST("/api/todos", func(ctx *gin.Context) {
		var todo todos.RawTodo
		if err := ctx.ShouldBindJSON(&todo); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		newTodo := todo.Create()
		*todoList = append(*todoList, newTodo)

		ctx.Header("location", fmt.Sprintf("/api/todos/%s", newTodo.Id))
		ctx.JSON(http.StatusCreated, gin.H{"status": "todo created"})
	})
	router.DELETE("/api/todos", func(ctx *gin.Context) {
		todoList = &[]todos.Todo{}
		ctx.Status(http.StatusOK)
	})

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
	_ = saveExistingTodos(todoList)

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
