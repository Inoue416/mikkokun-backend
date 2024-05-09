package main

import (
	"backend/handler"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	_ "backend/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type GreetingResponse struct {
	Message string `json:"Message"`
}

// greeting godoc
// @Summary      return a greeting
// @Description  return Hello World
// @Accept       json
// @Produce      json
// @Success      200 {object} GreetingResponse
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /greeting [get]
func greeting(c *gin.Context) {
	c.JSON(http.StatusOK, GreetingResponse{
		Message: "Hello World",
	})
}

// @title           MikkokuApp Backend API
// @version         1.0.0
// @description     密っこくんのバックエンドAPI

// @host      localhost:8080
// @BasePath  /

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	fmt.Println("⭐️⭐️⭐️  Start Server ⭐️⭐️⭐️ ")
	fmt.Println("🚀🚀🚀 Swagger docs 🚀🚀🚀")
	fmt.Println("URL  :  http://localhost:8080/swagger/index.html\n\n")

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

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// サーバーの起動
	r.Run(":8080")
}
