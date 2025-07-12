package user

import (
	"net/http"
	"todo-api-go/internal/service"

	"github.com/gin-gonic/gin"
)

type CreateRequest struct {
	Email    string `json:"email" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *UserHandler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.userService.Register(req.Email, req.Username, req.Password)
	if err != nil {
		switch e := err.(type) {
		// NOTE: 重複の場合は409を返す
		case *service.ErrDuplicateEmail:
			c.JSON(http.StatusConflict, gin.H{"error": e.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": e.Error()})
			return
		}

	}

	c.JSON(http.StatusCreated, response)
}
