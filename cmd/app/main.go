package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/wesleyklop/todo-api/v2/internal/todos"
	"github.com/wesleyklop/todo-api/v2/internal/tracing"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	log, _ := zap.NewDevelopment()
	defer log.Sync()

	tracer, err := tracing.NewTracer("http://jaeger-collector.default.svc:14268/api/traces")
	if err != nil {
		log.Fatal("Failed to create tracer", zap.Error(err))
	}
	otel.SetTracerProvider(tracer)

	repository, err := todos.LoadFromFile(ctx, "/mnt/data/todos.json")
	if err != nil {
		log.Fatal("Failed to create todo repository", zap.Error(err))
	}

	router := gin.New()
	router.Use(ginzap.Ginzap(log, time.RFC3339, false))
	router.Use(ginzap.RecoveryWithZap(log, true))
	router.Use(otelgin.Middleware("todo-api"))

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
	todos.PersistToFile(ctx, repository, "/mnt/data/todos.json")

	shutdownTime := 5 * time.Second
	if mode, exists := os.LookupEnv("MODE"); exists && mode == "development" {
		log.Info("skipping graceful shutdown")
		shutdownTime = 0 * time.Second
	}

	// Give Gin 5 seconds to handle inflight requests (Cloud Run gives us 10 before SIGKILL)
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTime)
	defer cancel()
	if err := service.Shutdown(ctx); err != nil {
		log.Fatal("shutdown deadline exceeded for server shutdown")
	}
	if err := tracer.Shutdown(ctx); err != nil {
		log.Fatal("shutdown deadline exceeded for tracer shutdown")
	}

	log.Info("exiting")
}
