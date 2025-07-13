package user

import (
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	registerUserHandler *RegisterUserHandler
}

func NewUserHandler(registerUserHandler *RegisterUserHandler) *UserHandler {
	return &UserHandler{
		registerUserHandler: registerUserHandler,
	}
}

func (h *UserHandler) RegisterRoutes(r *gin.RouterGroup) {
	users := r.Group("/users")
	{
		users.POST("", h.registerUserHandler.Handle)
	}
}
