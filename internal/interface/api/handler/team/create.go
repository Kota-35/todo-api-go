package team

import (
	"fmt"

	"github.com/gin-gonic/gin"

	teamDTO "todo-api-go/internal/application/dto/team"
	"todo-api-go/internal/application/usecase/team"
	"todo-api-go/internal/interface/api/middleware"
	"todo-api-go/internal/interface/api/response"
)

type CreateTeamHandler struct {
	createTeamUseCase *team.CreateTeamUseCase
}

func NewCreateTeamHandler(
	createTeamUserCase *team.CreateTeamUseCase,
) CreateTeamHandler {
	return CreateTeamHandler{
		createTeamUseCase: createTeamUserCase,
	}
}

func (h *CreateTeamHandler) Handle(c *gin.Context) {
	// リクエストの解析
	var input teamDTO.CreateTeamInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "リクエストが正しくありません", err)
		return
	}
	userId, err := middleware.GetUserID(c)
	if err != nil {
		fmt.Println(err.Error())
		response.BadRequest(c, "リクエストが正しくありません", err)
		return
	}

	output, err := h.createTeamUseCase.Execute(input, *userId)
	if err != nil {
		fmt.Println(err.Error())
		response.InternalServerError(c, "処理中にエラーが発生しました", err)
		return
	}

	response.Created(c, "チームを作成しました", output)
}
