package main

import (
	"backend/handler"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Printf("Start Server 🚀 \n")
	fmt.Printf(handler.Broadcast)
	// Ginルーターの初期化
	r := gin.Default()

	// WebSocketハンドラーの登録
	r.GET("/ws")

	// サーバーの起動
	r.Run(":8080")
}
