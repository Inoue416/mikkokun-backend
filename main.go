package main

import (
	"backend/handler"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func greeting(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello World!\n"})
}

func main() {
	fmt.Printf("⭐️⭐️⭐️  Start Server ⭐️⭐️⭐️ \n")
	// Ginルーターの初期化
	r := gin.Default()

	// WebSocketハンドラーの登録
	r.GET("/ws", handler.WebsocketHandler)
	r.GET("/greeting", greeting)

	// サーバーの起動
	r.Run(":8080")
}
