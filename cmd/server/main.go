package main

import (
	"fmt"
	"os"

	"todo-api-go/internal/application/bootstrap"

	"github.com/gin-gonic/gin"
)

func main() {
	// アプリケーション全体の依存関係を初期化
	app := bootstrap.NewApplication()

	// Ginルーターの初期化
	router := gin.Default()

	// ルートの設定
	v1 := router.Group("/api/v1")
	{
		app.Handlers.User.RegisterRoutes(v1)
		app.Handlers.Session.RegisterRoutes(v1)
		app.Handlers.Team.RegisterRoutes(v1)
	}

	// サーバーの起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // デフォルトポート
	}

	if app.Config.Env == "development" {
		// HTTPS証明書のパスを環境変数から取得
		fmt.Println(app.Config)
		router.RunTLS(":"+port, app.Config.SslCertPath, app.Config.SslKeyPath)
	} else {
		router.Run(":" + port)
	}
}
