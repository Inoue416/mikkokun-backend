package main

import (
	"log"
	"time"
	"net/http"

	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
)

// https server
func main() {
	router:= gin.Default()
	// Official default settings
	server := &http.Server{
		Addr: ":8080",
		Handler: router,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	server.ListenAndServe()

	router.GET("/greeting", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World!")
	})
	router.RunTLS(":8080", "")
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