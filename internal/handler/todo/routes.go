package todo

import (
	"todo-api-go/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	todoService *service.TodoService
}

func NewHandler(todoService *service.TodoService) *Handler {
	return &Handler{
		todoService: todoService,
	}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	todos := r.Group("/todos")
	{
		todos.GET("", h.List)
		todos.POST("", h.Create)
	}
}
