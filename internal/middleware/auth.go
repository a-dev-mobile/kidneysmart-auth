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

// AuthErrorResponse структура для ответов об ошибках аутентификации
type AuthErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// AuthMiddleware создает middleware для проверки JWT токена.
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Извлечение токена из заголовка Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, AuthErrorResponse{Status: "AUTHORIZATION_REQUIRED", Message: "Authorization header is required"})
			return
		}

		// Проверка формата токена (обычно это 'Bearer [token]')
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, AuthErrorResponse{Status: "INVALID_TOKEN_FORMAT", Message: "Invalid token format. Format should be 'Bearer [token]'."})
			return
		}

		// Парсинг и валидация токена
		userID, err := utils.ParseToken(token, jwtSecret, "access")
		if err != nil {
			var status string
			var message string
			switch err.Error() {
			case "Token expired":
				status = "TOKEN_EXPIRED"
				message = "Your token has expired. Please log in again."
			case "Invalid token":
				status = "INVALID_TOKEN"
				message = "The provided token is invalid. Check the token and try again."
			default:
				status = "AUTHENTICATION_FAILED"
				message = "Error occurred during token validation. Please try again."
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, AuthErrorResponse{Status: status, Message: message})
			return
		}

		// Добавление userID в контекст Gin
		c.Set(string(UserIDKey), userID)

		// Переход к следующему обработчику
		c.Next()
	}
}
