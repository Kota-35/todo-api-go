package session

import "github.com/gin-gonic/gin"

type SessionHandler struct {
	loginUserHandler         *LoginUserHandler
	getCurrentSessionHandler *GetCurrentSessionHandler
	sessionRefreshHandler    *SessionRefreshHandler
	logoutUserHandler        *LogoutUserHandler
}

func NewSessionHandler(
	loginUserHandler *LoginUserHandler,
	getCurrentSessionHandler *GetCurrentSessionHandler,
	sessionRefreshHandler *SessionRefreshHandler,
	logoutUserHandler *LogoutUserHandler,
) *SessionHandler {
	return &SessionHandler{
		loginUserHandler:         loginUserHandler,
		getCurrentSessionHandler: getCurrentSessionHandler,
		sessionRefreshHandler:    sessionRefreshHandler,
		logoutUserHandler:        logoutUserHandler,
	}
}

func (h *SessionHandler) RegisterRoutes(r *gin.RouterGroup) {
	sessions := r.Group("/sessions")
	{
		sessions.POST("", h.loginUserHandler.Handle)
		sessions.GET("/me", h.getCurrentSessionHandler.Handle)
		sessions.POST("/refresh", h.sessionRefreshHandler.Handle)
		sessions.DELETE("/", h.logoutUserHandler.Handle)
	}
}
