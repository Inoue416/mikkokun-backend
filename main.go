package main

import (
	"backend/handler"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func greeting(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello World!\n"})
}

func main() {
	fmt.Printf("⭐️⭐️⭐️  Start Server ⭐️⭐️⭐️ \n")
	// Ginルーターの初期化
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
		},
		AllowMethods: []string{
			"POST",
			"GET",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Access-Control-Allow-Credentials",
			"Access-Control-Allow-Headers",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"Authorization",
		},
		AllowCredentials: false,
	}))

	// WebSocketハンドラーの登録
	r.GET("/ws", handler.WebsocketHandler)
	r.GET("/checkSameSeatNumber", handler.CheckSameSeatNumber)
	r.GET("/greeting", greeting)

	// サーバーの起動
	r.Run(":8080")
}
