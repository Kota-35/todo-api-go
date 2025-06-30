package handler

import (
	"net/http"
	"todo-api-go/prisma/db"

	"github.com/gin-gonic/gin"
)

func GetAllTodos(c *gin.Context, client *db.PrismaClient) {
	tasks, err := client.Task.FindMany().OrderBy().Exec(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}
