package team

import (
	"fmt"
	"todo-api-go/internal/application/usecase/team"
	"todo-api-go/internal/interface/api/middleware"
	"todo-api-go/internal/interface/api/response"

	"github.com/gin-gonic/gin"
)

type FindMyTeamsHandler struct {
	findMyTeamUseCase *team.FindMyTeamsUseCase
}

func NewFindMyTeamsHandler(
	findMyTeamUseCase *team.FindMyTeamsUseCase,
) *FindMyTeamsHandler {
	return &FindMyTeamsHandler{
		findMyTeamUseCase: findMyTeamUseCase,
	}
}

func (h *FindMyTeamsHandler) Handle(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.BadRequest(c, "リクエストが正しくありません", err)
		return
	}

	output, err := h.findMyTeamUseCase.Exec(*userID)
	if err != nil {
		fmt.Println(err.Error())
		response.InternalServerError(c, "処理中にエラーが発生しました", err)
	}

	response.OK(c, "ユーザーのチームを取得しました", output)
}
