package middleware

import (

	"github.com/a-dev-mobile/kidneysmart-auth/internal/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

func CORSMiddleware(cfg config.Config, lg *slog.Logger) gin.HandlerFunc {
	// Create a map of allowed sources once during middleware initialization
	allowedOrigins := map[string]bool{}
	for _, allowedOrigin := range cfg.ClientConnection.AllowedOrigins {
		allowedOrigins[allowedOrigin] = true
	}
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		// Разрешить запросы без Origin в режиме разработки
		if origin == "" && cfg.ClientConnection.GinMode == "debug" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		} else if _, ok := allowedOrigins[origin]; ok {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			lg.Warn("CORS policy rejection", slog.String("origin", origin))
			c.JSON(http.StatusForbidden, gin.H{"error": "CORS policy: Origin not allowed"})
			c.Abort() // This stops processing the request and prevents subsequent build from executing.
			return
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
