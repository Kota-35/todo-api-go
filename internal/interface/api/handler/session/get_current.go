package session

import (
	"fmt"
	"strings"
	"todo-api-go/internal/application/usecase/auth"
	"todo-api-go/internal/interface/api/response"

	"github.com/gin-gonic/gin"

	authDTO "todo-api-go/internal/application/dto/auth"
)

type GetCurrentSessionHandler struct {
	getCurrentSessionUseCase *auth.GetCurrentSessionUseCase
}

func NewGetCurrentSessionHandler(
	getCurrentSessionUseCase *auth.GetCurrentSessionUseCase,
) *GetCurrentSessionHandler {
	return &GetCurrentSessionHandler{
		getCurrentSessionUseCase: getCurrentSessionUseCase,
	}
}

func (h *GetCurrentSessionHandler) Handle(c *gin.Context) {
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

	// 2. ユースケースでトークンを検証
	input := authDTO.GetCurrentSessionInput{
		JWTToken: accessToken,
	}
	output, err := h.getCurrentSessionUseCase.Execute(input)

	if err != nil {
		response.AbortWithUnauthorizedError(c, "認証に失敗しました", err)
		return
	}

	// 3. レスポンス返却

	response.OK(c, "ok", output)
}
