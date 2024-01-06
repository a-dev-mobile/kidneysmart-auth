package middleware

import (
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func LogHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentTime := time.Now().Format("2006/01/02 15:04:05")
		log.Printf("\n[%s] Headers for %s %s:\n", currentTime, c.Request.Method, c.Request.URL.Path)
		for k, v := range c.Request.Header {
			log.Printf("  %s: %s", k, strings.Join(v, ","))
		}
		c.Next()
	}
}