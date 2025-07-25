package session

import (
	"time"
	"todo-api-go/internal/application/usecase/auth"
	"todo-api-go/internal/config"
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

	cfg := config.LoadEnv()

	if cfg.Env == "development" {
		// RefreshTokenクッキーの設定
		c.SetCookie(
			"__Host-refresh",
			output.RefreshToken,
			int(time.Until(output.RefreshTokenExpiresAt).Seconds()),
			"/",
			"localhost",
			false, // developmentではHTTPSを使わない場合があるためfalse
			true,
		)

	} else if cfg.Env == "production" {

		// RefreshTokenクッキーの設定
		c.SetCookie(
			"__Host-refresh",
			output.RefreshToken,
			int(time.Until(output.RefreshTokenExpiresAt).Seconds()),
			"/",
			"localhost", // productionでは適切なドメインに変更
			true,
			true,
		)
	}

	data := map[string]string{
		"accessToken": output.AccessToken,
		"expiresAt":   output.AccessTokenExpiresAt.Format(time.RFC3339),
	}

	response.OK(c, "ログインしました", data)
}
