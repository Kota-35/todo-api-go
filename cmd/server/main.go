package main

import (
	"os"
	"todo-api-go/internal/handler/user"
	"todo-api-go/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// taskService := service.TaskService{}

	// taskHandler := task.NewHandler(&taskService)

	jwtSecret := os.Getenv("JWT_SECRET")
	userService := service.NewUserService(jwtSecret)
	userHandler := user.NewHandler(userService)

	v1 := router.Group("/api/v1")
	{
		// taskHandler.RegisterRoutes(v1)
		userHandler.RegisterRoutes(v1)
	}

	router.Run(":8080")
}
