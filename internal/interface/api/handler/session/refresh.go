package session

import (
	"time"
	"todo-api-go/internal/application/usecase/auth"
	"todo-api-go/internal/config"
	"todo-api-go/internal/domain/valueobject"
	"todo-api-go/internal/interface/api/response"

	authDTO "todo-api-go/internal/application/dto/auth"

	"github.com/gin-gonic/gin"
)

type SessionRefreshHandler struct {
	pepper                *valueobject.Pepper
	refreshSessionUseCase *auth.RefreshSessionUseCase
}

func NewSessionRefreshHandler(
	pepper *valueobject.Pepper,
	refreshSessionUseCase *auth.RefreshSessionUseCase,
) *SessionRefreshHandler {
	return &SessionRefreshHandler{
		pepper:                pepper,
		refreshSessionUseCase: refreshSessionUseCase,
	}
}

func (h *SessionRefreshHandler) Handle(c *gin.Context) {
	refreshToken, err := c.Cookie("__Host-refresh")
	if err != nil {
		response.AbortWithUnauthorizedError(c, "リフレッシュトークンがありません", err)
		return
	}

	// リフレッシュトークンの検証 valueobject
	refreshTokenVO, err := valueobject.NewRefreshToken(refreshToken, *h.pepper)
	if err != nil {
		response.AbortWithUnauthorizedError(c, "リフレッシュトークンが正しくありません", err)
		return
	}

	input := authDTO.RefreshSessionInput{
		RefreshTokenVO: *refreshTokenVO,
	}

	output, err := h.refreshSessionUseCase.Execute(&input)
	if err != nil {
		response.AbortWithUnauthorizedError(c, "認証に失敗しました", err)
	}

	cfg := config.LoadEnv()

	if cfg.Env == "development" {
		// AccessTokenクッキーの設定
		c.SetCookie(
			"__Host-session",
			output.AccessToken,
			int(time.Until(output.AccessTokenExpiresAt).Seconds()),
			"/",
			"localhost",
			false, // developmentではHTTPSを使わない場合があるためfalse
			true,
		)

	} else if cfg.Env == "production" {
		// AccessTokenクッキーの設定
		c.SetCookie(
			"__Host-session",
			output.AccessToken,
			int(time.Until(output.AccessTokenExpiresAt).Seconds()),
			"/",
			"localhost", // productionでは適切なドメインに変更
			true,
			true,
		)

	}

	response.OK(c, "アクセストークンを再生成しました", nil)

}
