package team

import (
	"todo-api-go/internal/interface/api/middleware"

	"github.com/gin-gonic/gin"
)

type TeamHandler struct {
	createTeamHandler  *CreateTeamHandler
	authMiddleware     *middleware.AuthMiddleware
	findMyTeamsHandler *FindMyTeamsHandler
}

func NewTeamHandler(
	createTeamHandler *CreateTeamHandler,
	authMiddleware *middleware.AuthMiddleware,
	findMyTeamsHandler *FindMyTeamsHandler,
) *TeamHandler {
	return &TeamHandler{
		createTeamHandler: createTeamHandler,
		authMiddleware:    authMiddleware,
	}
}

func (h *TeamHandler) RegisterRoutes(r *gin.RouterGroup) {
	teams := r.Group("/teams")
	teams.Use(h.authMiddleware.RequireAuth())
	{
		// チーム作成
		teams.POST("", h.createTeamHandler.Handle)
		teams.GET("/me", h.findMyTeamsHandler.Handle)
	}
}
