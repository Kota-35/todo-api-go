package main

import (
	"fmt"
	"os"

	"todo-api-go/internal/application/usecase/auth"
	"todo-api-go/internal/application/usecase/team"
	"todo-api-go/internal/application/usecase/user"
	"todo-api-go/internal/config"
	"todo-api-go/internal/domain/security"
	"todo-api-go/internal/domain/valueobject"
	"todo-api-go/internal/infrastructure/persistence/repository"

	sessionHandler "todo-api-go/internal/interface/api/handler/session"
	teamHandler "todo-api-go/internal/interface/api/handler/team"
	userHandler "todo-api-go/internal/interface/api/handler/user"
	"todo-api-go/internal/interface/api/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// 環境変数の取得
	cfg := config.LoadEnv()
	jwtSecret := cfg.JwtSecret
	if jwtSecret == "" {
		jwtSecret = "default-secret-key" // 開発用デフォルト値
	}

	// リポジトリの初期化
	userRepo := repository.NewUserRepository()
	sessionRepo := repository.NewSessionRepository()
	jwtGenerator := security.NewJWTGenerator(jwtSecret)
	refreshTokenGenerator := security.NewRefreshTokenGenerator()
	teamRepo := repository.NewTeamRepository()

	// 認証用設定
	pepper := valueobject.Pepper([]byte("default-pepper-key"))

	authMiddleware := middleware.NewAuthMiddleware(
		jwtGenerator,
		sessionRepo,
	)

	// ユースケースの初期化
	// userのユースケース
	registerUserUseCase := user.NewRegisterUserUseCase(userRepo)
	authenticateUserUseCase := auth.NewAuthenticateUserUseCase(
		userRepo,
		sessionRepo,
		jwtGenerator,
		refreshTokenGenerator,
		pepper,
	)

	// authのユースケース
	getCurrentSessionUseCase := auth.NewGetCurrentSessionUseCase(
		jwtGenerator,
		userRepo,
		sessionRepo,
	)
	refreshSessionUseCase := auth.NewRefreshSessionUseCase(
		sessionRepo,
		userRepo,
		jwtGenerator,
	)
	logoutUseCase := auth.NewLogoutUserCase(
		pepper,
		sessionRepo,
	)

	// teamのユースケース
	createTeamUseCase := team.NewCreateTeamUseCase(teamRepo)
	findMyTeamsUseCase := team.NewFindMyTeamsUseCase(teamRepo)

	// ハンドラーの初期化
	// /usersのハンドラー
	registerUserHandler := userHandler.NewRegisterUserHandler(
		registerUserUseCase,
	)
	userHandlerInstance := userHandler.NewUserHandler(registerUserHandler)

	// /sessionsのハンドラー
	loginUserHandler := sessionHandler.NewLoginUserHandler(
		authenticateUserUseCase,
	)
	getCurrentSessionHandler := sessionHandler.NewGetCurrentSessionHandler(
		getCurrentSessionUseCase,
	)
	sessionRefreshHandler := sessionHandler.NewSessionRefreshHandler(
		&pepper,
		refreshSessionUseCase,
	)
	logoutUserHandler := sessionHandler.NewLogoutUserHandler(
		jwtGenerator,
		logoutUseCase,
	)
	sessionHandler := sessionHandler.NewSessionHandler(
		loginUserHandler,
		getCurrentSessionHandler,
		sessionRefreshHandler,
		logoutUserHandler,
		authMiddleware,
	)

	// /teamsのハンドラー

	teamCreateHandler := teamHandler.NewCreateTeamHandler(createTeamUseCase)
	findMyTeamsHandler := teamHandler.NewFindMyTeamsHandler(findMyTeamsUseCase)
	teamHandler := teamHandler.NewTeamHandler(
		teamCreateHandler,
		authMiddleware,
		findMyTeamsHandler,
	)

	// Ginルーターの初期化
	router := gin.Default()

	// ルートの設定
	v1 := router.Group("/api/v1")
	{
		userHandlerInstance.RegisterRoutes(v1)
		sessionHandler.RegisterRoutes(v1)
		teamHandler.RegisterRoutes(v1)
	}

	// サーバーの起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // デフォルトポート
	}

	if cfg.Env == "development" {
		// HTTPS証明書のパスを環境変数から取得
		fmt.Println(cfg)
		router.RunTLS(":"+port, cfg.SslCertPath, cfg.SslKeyPath)
	} else {
		router.Run(":" + port)
	}
}
