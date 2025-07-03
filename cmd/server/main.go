package main

import (
	"todo-api-go/internal/handler/todo"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/todos", func(c *gin.Context) {
		todo.GetAllTodos(c)
	})

	router.Run(":8080")
}
