package middleware

import (
	"net/http"
	"strings"

	"github.com/a-dev-mobile/kidneysmart-auth/internal/utils"
	"github.com/gin-gonic/gin"
	// "yourapp/jwt" // Предполагается, что ParseToken находится в этом пакете
)

// ContextKey - тип для ключей контекста
type ContextKey string

// Определение константы для ключа userID
const UserIDKey ContextKey = "userID"

// AuthMiddleware создает middleware для проверки JWT токена.
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Извлечение токена из заголовка Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Authorization header is required"})
			return
		}

		// Проверка формата токена (обычно это 'Bearer [token]')
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token format"})
			return
		}

		// Парсинг и валидация токена
		userID, err := utils.ParseToken(token, jwtSecret, "access")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}

		// Добавление userID в контекст Gin
		c.Set("userID", userID)

		// Переход к следующему обработчику
		c.Next()
	}
}
