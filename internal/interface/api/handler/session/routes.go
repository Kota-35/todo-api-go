package session

import (
	"todo-api-go/internal/interface/api/middleware"

	"github.com/gin-gonic/gin"
)

type SessionHandler struct {
	loginUserHandler         *LoginUserHandler
	getCurrentSessionHandler *GetCurrentSessionHandler
	sessionRefreshHandler    *SessionRefreshHandler
	logoutUserHandler        *LogoutUserHandler
	authMiddleware           *middleware.AuthMiddleware
}

func NewSessionHandler(
	loginUserHandler *LoginUserHandler,
	getCurrentSessionHandler *GetCurrentSessionHandler,
	sessionRefreshHandler *SessionRefreshHandler,
	logoutUserHandler *LogoutUserHandler,
	authMiddleware *middleware.AuthMiddleware,
) *SessionHandler {
	return &SessionHandler{
		loginUserHandler:         loginUserHandler,
		getCurrentSessionHandler: getCurrentSessionHandler,
		sessionRefreshHandler:    sessionRefreshHandler,
		logoutUserHandler:        logoutUserHandler,
		authMiddleware:           authMiddleware,
	}
}

func (h *SessionHandler) RegisterRoutes(r *gin.RouterGroup) {
	sessions := r.Group("/sessions")
	{
		sessions.POST("", h.loginUserHandler.Handle)
		sessions.GET(
			"/me",
			h.getCurrentSessionHandler.Handle,
		)
		sessions.POST(
			"/refresh",
			h.sessionRefreshHandler.Handle,
		)
		sessions.DELETE("/", h.logoutUserHandler.Handle)
	}
}
