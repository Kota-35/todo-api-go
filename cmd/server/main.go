package main

import (
	"todo-api-go/internal/handler"
	"todo-api-go/prisma/db"

	"github.com/gin-gonic/gin"
)

func main() {
	client := db.NewClient()
	router := gin.Default()

	router.GET("/todos", func(c *gin.Context) {
		handler.GetAllTodos(c, client)
	})

	router.Run(":8080")
}
