package session

import (
	"todo-api-go/internal/application/usecase/auth"
	"todo-api-go/internal/interface/api/response"

	"github.com/gin-gonic/gin"

	errorHandler "todo-api-go/internal/interface/api/error"

	authDTO "todo-api-go/internal/application/dto/auth"
)

type LoginUserHandler struct {
	authenticateUserUserCase *auth.AuthenticateUserUseCase
}

func NewLoginUserHandler(
	authenticationUserUseCase *auth.AuthenticateUserUseCase,
) *LoginUserHandler {
	return &LoginUserHandler{
		authenticateUserUserCase: authenticationUserUseCase,
	}
}

func (h *LoginUserHandler) Handle(c *gin.Context) {
	// リクエストの解析
	var input authDTO.AuthenticateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "リクエストが正しくありません", err)
		return
	}

	output, err := h.authenticateUserUserCase.Execute(input)
	if err != nil {
		// エラーハンドリング
		switch {
		case errorHandler.IsAuthenticationError(err):
			response.UnauthorizedError(c, "認証に失敗しました", err)
		default:
			response.InternalServerError(c, "ログインに失敗しました", err)
		}

		return
	}

	response.OK(c, "ログインしました", output)
}
