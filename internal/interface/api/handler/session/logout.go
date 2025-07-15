package session

import (
	"fmt"
	"todo-api-go/internal/application/usecase/auth"
	"todo-api-go/internal/config"
	"todo-api-go/internal/domain/security"
	"todo-api-go/internal/interface/api/response"

	authDTO "todo-api-go/internal/application/dto/auth"

	"github.com/gin-gonic/gin"
)

type LogoutUserHandler struct {
	jwtGenerator security.JWTGenerator

	logoutUseCase *auth.LogoutUseCase
}

func NewLogoutUserHandler(
	jwtGenerator security.JWTGenerator,
	logoutUseCase *auth.LogoutUseCase,
) *LogoutUserHandler {
	return &LogoutUserHandler{
		jwtGenerator:  jwtGenerator,
		logoutUseCase: logoutUseCase,
	}
}

func (h *LogoutUserHandler) Handle(c *gin.Context) {
	// 1. クッキーからJWTを取得
	accessToken, err := c.Cookie("__Host-session")
	if err != nil {
		fmt.Print("[NewLogoutUserHandler]クッキーの取得に失敗しました:\n", err)
		h.clearCookies(c)
		response.OK(c, "ログアウトしました", nil)
		return
	}

	// 2. JWTを検証
	claims, err := h.jwtGenerator.VerifyAccessToken(accessToken)
	if err != nil {
		fmt.Print("[NewLogoutUserHandler]JWTの検証に失敗しました:\n", err)
		// JWT検証失敗時もクッキーをクリア
		h.clearCookies(c)
		response.OK(c, "ログアウトしました", nil)
		return
	}

	input := authDTO.LogoutInput{
		RefreshTokenId: claims.RefreshTokenID,
	}

	fmt.Print("Claims", claims)

	if err := h.logoutUseCase.Execute(input); err != nil {
		// セッションの更新に失敗した場合でもクッキーをクリア
		fmt.Print("[NewLogoutUserHandler]セッションの更新に失敗しました:\n", err)
		h.clearCookies(c)
		response.OK(c, "ログアウトしました", nil)
		return
	}

	// 正常終了
	h.clearCookies(c)
	response.OK(c, "ログラウトしました", nil)

}

func (h *LogoutUserHandler) clearCookies(c *gin.Context) {
	cfg := config.LoadEnv()

	if cfg.Env == "development" {
		// AccessTokenクッキーの設定
		c.SetCookie(
			"__Host-session",
			"",
			-1,
			"/",
			"",
			false, // developmentではHTTPSを使わない場合があるためfalse
			true,
		)

		// RefreshTokenクッキーの設定
		c.SetCookie(
			"__Host-refresh",
			"",
			-1,
			"/",
			"",
			false, // developmentではHTTPSを使わない場合があるためfalse
			true,
		)

	} else if cfg.Env == "production" {
		// AccessTokenクッキーの設定
		c.SetCookie(
			"__Host-session",
			"",
			-1,
			"/",
			"", // productionでは適切なドメインに変更
			true,
			true,
		)

		// RefreshTokenクッキーの設定
		c.SetCookie(
			"__Host-refresh",
			"",
			-1,
			"/",
			"", // productionでは適切なドメインに変更
			true,
			true,
		)
	}
}
