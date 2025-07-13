package session

import "github.com/gin-gonic/gin"

type SessionHandler struct {
	loginUserHandler *LoginUserHandler
}

func NewSessionHandler(loginUserHandler *LoginUserHandler) *SessionHandler {
	return &SessionHandler{
		loginUserHandler: loginUserHandler,
	}
}

func (h *SessionHandler) RegisterRoutes(r *gin.RouterGroup) {
	sessions := r.Group("/sessions")
	{
		sessions.POST("", h.loginUserHandler.Handle)
	}
}
