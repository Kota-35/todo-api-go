package task

// import (
// 	"todo-api-go/internal/service"

// 	"github.com/gin-gonic/gin"
// )

// type Handler struct {
// 	taskService *service.TaskService
// }

// func NewHandler(taskService *service.TaskService) *Handler {
// 	return &Handler{
// 		taskService: taskService,
// 	}
// }

// func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
// 	tasks := r.Group("/tasks")
// 	{
// 		tasks.GET("", h.List)
// 		tasks.POST("", h.Create)
// 		tasks.GET("/:taskId", h.Get)
// 	}
// }
