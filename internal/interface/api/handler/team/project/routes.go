package project

import "github.com/gin-gonic/gin"

type ProjectHandler struct {
	listProjectsHandler  *ListProjectsHandler
	createProjectHandler *CreateProjectHandler
}

func NewProjectHandler(
	listProjectsHandler *ListProjectsHandler,
	createProjectHandler *CreateProjectHandler,
) *ProjectHandler {
	return &ProjectHandler{
		listProjectsHandler:  listProjectsHandler,
		createProjectHandler: createProjectHandler,
	}
}

func (h *ProjectHandler) RegisterRoutes(r *gin.RouterGroup) {
	projects := r.Group("/projects")
	{
		projects.GET("", h.listProjectsHandler.Handle)
		projects.POST("", h.createProjectHandler.Handle)
	}
}
