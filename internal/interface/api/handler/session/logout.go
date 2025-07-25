package session

import (
	"fmt"
	"strings"
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
	// 1. BearerからJWTを取得
	authHeader := c.GetHeader("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer") {
		fmt.Println("アクセストークンがありません")
		response.AbortWithUnauthorizedError(
			c,
			"アクセストークンがありません",
			fmt.Errorf("アクセストークンがありません"),
		)
		return
	}

	accessToken := strings.TrimPrefix(authHeader, "Bearer ")

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
	isSecure := cfg.ShouldUseSecureCookies()

	// AccessTokenクッキーのクリア
	c.SetCookie(
		"__Host-session",
		"",
		-1,
		"/",
		"",
		isSecure,
		true,
	)

	// RefreshTokenクッキーのクリア
	c.SetCookie(
		"__Host-refresh",
		"",
		-1,
		"/",
		"",
		isSecure,
		true,
	)
}
