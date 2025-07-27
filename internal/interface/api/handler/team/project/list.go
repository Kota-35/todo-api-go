package project

import (
	"fmt"
	"todo-api-go/internal/application/usecase/project"
	"todo-api-go/internal/interface/api/response"

	"github.com/gin-gonic/gin"
)

type ListProjectsHandler struct {
	getProjectsUseCase *project.GetProjectsUseCase
}

func NewListProjectsHandler(
	getProjectsUseCase *project.GetProjectsUseCase,
) *ListProjectsHandler {
	return &ListProjectsHandler{
		getProjectsUseCase: getProjectsUseCase,
	}
}

func (h *ListProjectsHandler) Handle(c *gin.Context) {
	// teamIdパラメータを取得
	teamID := c.Param("teamId")
	if teamID == "" {
		response.BadRequest(c, "チームIDが指定されていません", nil)
		return
	}

	output, err := h.getProjectsUseCase.Execute(teamID)
	if err != nil {
		fmt.Println(err.Error())
		response.InternalServerError(c, "処理中にエラーが発生しました", err)
		return
	}

	response.OK(c, "プロジェクトの一覧を取得しました", output.Projects)
}
