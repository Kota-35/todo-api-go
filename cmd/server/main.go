package main

import (
	"todo-api-go/internal/handler/task"
	"todo-api-go/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	taskService := service.TaskService{}

	taskHandler := task.NewHandler(&taskService)

	v1 := router.Group("/api/v1")
	{
		taskHandler.RegisterRoutes(v1)
	}

	router.Run(":8080")
}
