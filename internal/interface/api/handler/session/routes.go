package session

import "github.com/gin-gonic/gin"

type SessionHandler struct {
	loginUserHandler         *LoginUserHandler
	getCurrentSessionHandler *GetCurrentSessionHandler
}

func NewSessionHandler(loginUserHandler *LoginUserHandler, getCurrentSessionHandler *GetCurrentSessionHandler) *SessionHandler {
	return &SessionHandler{
		loginUserHandler:         loginUserHandler,
		getCurrentSessionHandler: getCurrentSessionHandler,
	}
}

func (h *SessionHandler) RegisterRoutes(r *gin.RouterGroup) {
	sessions := r.Group("/sessions")
	{
		sessions.POST("", h.loginUserHandler.Handle)
		sessions.GET("/me", h.getCurrentSessionHandler.Handle)
	}
}
