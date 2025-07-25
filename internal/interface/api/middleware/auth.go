package middleware

import (
	"fmt"
	"time"
	"todo-api-go/internal/domain/repository"
	"todo-api-go/internal/domain/security"
	"todo-api-go/internal/interface/api/response"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	jwtGenerator security.JWTGenerator
	sessionRepo  repository.SessionRepository
}

func NewAuthMiddleware(
	jwtGenerator security.JWTGenerator,
	sessionRepo repository.SessionRepository,
) *AuthMiddleware {
	return &AuthMiddleware{
		jwtGenerator: jwtGenerator,
		sessionRepo:  sessionRepo,
	}
}

// 認証が必要なエンドポイント用のミドルウェア
func (am *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. クッキーからJWTを取得
		accessToken, err := c.Cookie("__Host-session")
		if err != nil {
			fmt.Println("アクセストークンがありません", err)
			response.AbortWithUnauthorizedError(c, "アクセストークンがありません", err)
			return
		}

		// 2. JWTを検証
		claims, err := am.jwtGenerator.VerifyAccessToken(accessToken)
		if err != nil {
			fmt.Println("無効なアクセストークンです", err)
			response.AbortWithUnauthorizedError(c, "無効なアクセストークンです", err)
			return
		}

		// 3. セッションの有効性を確認
		session, err := am.sessionRepo.FindByID(claims.RefreshTokenID)
		if err != nil {
			fmt.Println("セッションが見つかりません", err)
			response.AbortWithUnauthorizedError(c, "セッションが見つかりません", err)
			return
		}

		if session.IsRevoked() {
			fmt.Println("無効なセッションです", err)
			response.AbortWithUnauthorizedError(c, "無効なセッションです", err)
			return
		}

		if !session.IsActive(time.Now()) {
			fmt.Println("無効なセッションです", err)
			response.AbortWithUnauthorizedError(c, "無効なセッションです", err)
			return
		}

		// 5. ユーザー情報をコンテキストに設定
		c.Set("user_id", claims.UserID)
		c.Set("refresh_token_id", claims.RefreshTokenID)
		c.Set("jwt_claims", claims)

		// 次のハンドラーに処理を移す
		c.Next()
	}
}
