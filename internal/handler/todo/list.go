package todo

import (
	"net/http"

	"todo-api-go/pkg/database"

	"github.com/gin-gonic/gin"
)

func GetAllTodos(c *gin.Context) {
	tasks, err := database.PrismaClient.Task.FindMany().OrderBy().Exec(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}
