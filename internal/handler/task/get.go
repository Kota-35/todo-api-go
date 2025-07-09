package task

// import (
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// )

// func (h *Handler) Get(c *gin.Context) {
// 	taskId := c.Param("taskId")

// 	todo, err := h.taskService.GetById(taskId)

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, todo)
// }
