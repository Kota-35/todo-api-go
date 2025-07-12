package main

import (
	"os"
	"todo-api-go/internal/application/usecase/user"
	"todo-api-go/internal/infrastructure/persistence/repository"
	userHandler "todo-api-go/internal/interface/api/handler/user"

	"github.com/gin-gonic/gin"
)

func main() {

	// リポジトリの初期化
	userRepo := repository.NewUserRepository()

	// ユースケースの初期化
	registerUserUseCase := user.NewRegisterUserUseCase(userRepo)

	// ハンドラーの初期化
	registerUserHandler := userHandler.NewRegisterUserHandler(registerUserUseCase)
	userHandlerInstance := userHandler.NewHandler(registerUserHandler)

	// 環境変数の取得
	jwtSecret := os.Getenv("JWT_SECRET")
	_ = jwtSecret // 将来的に使用予定

	// Ginルーターの初期化
	router := gin.Default()

	// ルートの設定
	v1 := router.Group("/api/v1")
	{
		userHandlerInstance.RegisterRoutes(v1)
	}

	// サーバーの起動
	router.Run(":8080")
}
