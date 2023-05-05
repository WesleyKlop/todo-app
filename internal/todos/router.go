package todos

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewTodoRouter(router *gin.RouterGroup, repo *TodoRepository) *gin.RouterGroup {
	router.GET("/api/todos", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, repo.List())
	})
	router.GET("/api/todos/:todo", func(ctx *gin.Context) {
		todo := repo.Get(ctx.Param("todo"))
		if todo != nil {
			ctx.JSON(http.StatusOK, todo)
		} else {
			ctx.AbortWithStatus(http.StatusNotFound)
		}
	})
	router.POST("/api/todos", func(ctx *gin.Context) {
		var todo *RawTodo
		if err := ctx.ShouldBindJSON(todo); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		newTodo := repo.Create(todo)

		ctx.Header("location", fmt.Sprintf("/api/todos/%s", newTodo.Id))
		ctx.JSON(http.StatusCreated, gin.H{"status": "todo created"})
	})
	router.DELETE("/api/todos", func(ctx *gin.Context) {
		repo.Clear()
		ctx.Status(http.StatusOK)
	})
	return router
}
