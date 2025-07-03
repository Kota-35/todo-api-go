package main

import (
	"todo-api-go/internal/handler/todo"
	"todo-api-go/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	todoService := service.TodoService{}

	todoHandler := todo.NewHandler(&todoService)

	v1 := router.Group("/api/v1")
	{
		todoHandler.RegisterRoutes(v1)
	}

	router.Run(":8080")
}
