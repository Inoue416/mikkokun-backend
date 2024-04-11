package main

import (
	"log"
	// "time"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func loadEnv() (string, string) {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
		log.Printf("Error: %v", err)
		return "", ""
	}

	serverCrt := os.Getenv("SERVER_CRT_PATH")
	serverKey := os.Getenv("SERVER_KEY_PATH")
	return serverCrt, serverKey
}

// https server
func main() {
	router := gin.Default()
	// Official default settings
	// server := &http.Server{
	// 	Addr: ":8080",
	// 	Handler: router,
	// 	ReadTimeout: 10 * time.Second,
	// 	WriteTimeout: 10 * time.Second,
	// 	MaxHeaderBytes: 1 << 20,
	// }

	// server.ListenAndServe()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World!")
	})

	serverCrt, serverKey := loadEnv()
	log.Println(serverCrt, serverKey)
	router.RunTLS("localhost:8080", serverCrt, serverKey)
}

/* 今後の実装に利用するかも
package main

import (
  "log"
  "net/http"

  "github.com/gin-gonic/autotls"
  "github.com/gin-gonic/gin"
  "golang.org/x/crypto/acme/autocert"
)

func main() {
  r := gin.Default()

  // Ping handler
  r.GET("/ping", func(c *gin.Context) {
    c.String(http.StatusOK, "pong")
  })

  m := autocert.Manager{
    Prompt:     autocert.AcceptTOS,
    HostPolicy: autocert.HostWhitelist("example1.com", "example2.com"),
    Cache:      autocert.DirCache("/var/www/.cache"),
  }

  log.Fatal(autotls.RunWithManager(r, &m))
}
*/
