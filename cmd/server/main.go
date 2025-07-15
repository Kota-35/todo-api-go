package main

import (
	"os"

	"todo-api-go/internal/application/usecase/auth"
	"todo-api-go/internal/application/usecase/user"
	"todo-api-go/internal/domain/security"
	"todo-api-go/internal/domain/valueobject"
	"todo-api-go/internal/infrastructure/persistence/repository"

	sessionHandler "todo-api-go/internal/interface/api/handler/session"
	userHandler "todo-api-go/internal/interface/api/handler/user"

	"github.com/gin-gonic/gin"
)

func main() {
	// 環境変数の取得
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-key" // 開発用デフォルト値
	}

	// リポジトリの初期化
	userRepo := repository.NewUserRepository()
	sessionRepo := repository.NewSessionRepository()
	jwtGenerator := security.NewJWTGenerator(jwtSecret)
	refreshTokenGenerator := security.NewRefreshTokenGenerator()

	// 認証用設定
	pepper := valueobject.Pepper([]byte("default-pepper-key"))

	// ユースケースの初期化
	registerUserUseCase := user.NewRegisterUserUseCase(userRepo)
	authenticateUserUseCase := auth.NewAuthenticateUserUseCase(
		userRepo,
		sessionRepo,
		jwtGenerator,
		refreshTokenGenerator,
		pepper,
	)
	getCurrentSessionUseCase := auth.NewGetCurrentSessionUseCase(
		jwtGenerator,
		userRepo,
	)
	refreshSessionUseCase := auth.NewRefreshSessionUseCase(
		sessionRepo,
		userRepo,
		jwtGenerator,
	)

	// ハンドラーの初期化
	registerUserHandler := userHandler.NewRegisterUserHandler(
		registerUserUseCase,
	)
	userHandlerInstance := userHandler.NewUserHandler(registerUserHandler)

	loginUserHandler := sessionHandler.NewLoginUserHandler(
		authenticateUserUseCase,
	)
	getCurrentSessionHandler := sessionHandler.NewGetCurrentSessionHandler(
		*getCurrentSessionUseCase,
	)
	sessionRefreshHandler := sessionHandler.NewSessionRefreshHandler(
		&pepper,
		refreshSessionUseCase,
	)
	sessionHandler := sessionHandler.NewSessionHandler(
		loginUserHandler,
		getCurrentSessionHandler,
		sessionRefreshHandler,
	)

	// Ginルーターの初期化
	router := gin.Default()

	// ルートの設定
	v1 := router.Group("/api/v1")
	{
		userHandlerInstance.RegisterRoutes(v1)
		sessionHandler.RegisterRoutes(v1)
	}

	// サーバーの起動
	router.Run(":8080")
}
