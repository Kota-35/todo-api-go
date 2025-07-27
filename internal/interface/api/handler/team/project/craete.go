package project

import (
	"fmt"
	"todo-api-go/internal/application/usecase/project"
	"todo-api-go/internal/interface/api/middleware"
	"todo-api-go/internal/interface/api/response"

	"github.com/gin-gonic/gin"

	projectDTO "todo-api-go/internal/application/dto/project"
)

type CreateProjectHandler struct {
	createProjectUseCase *project.CreateProjectUseCase
}

func NewCraeteProjectHandler(
	createProjectHandler *project.CreateProjectUseCase,
) *CreateProjectHandler {
	return &CreateProjectHandler{
		createProjectUseCase: createProjectHandler,
	}
}

func (h *CreateProjectHandler) Handle(c *gin.Context) {
	var input projectDTO.CreateProjectInput

	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "リクエストが正しくありません", err)
		return
	}

	// teamIdパラメータを取得
	teamID := c.Param("teamId")
	if teamID == "" {
		response.BadRequest(c, "チームIDが指定されていません", nil)
		return
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.BadRequest(c, "リクエストが正しくありません", err)
		return
	}

	output, err := h.createProjectUseCase.Execute(input, userID, teamID)
	if err != nil {
		fmt.Println(err.Error())
		response.InternalServerError(c, "処理中にエラーが発生しました", err)
		return
	}

	response.Created(c, "プロジェクトを作成しました", output)
}
