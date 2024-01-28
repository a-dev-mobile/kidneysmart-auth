package middleware

import (
	"net/http"
	"strings"

	"github.com/a-dev-mobile/kidneysmart-auth/internal/utils" // Убедитесь, что путь к пакету корректен
	"github.com/gin-gonic/gin"
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
			switch err {
			case utils.ErrTokenExpired:
				status = "TOKEN_EXPIRED"
				message = "Your token has expired. Please log in again."
			case utils.ErrInvalidToken:
				status = "INVALID_TOKEN"
				message = "The provided token is invalid. Check the token and try again."	
				case  utils.ErrInvalidTokenType:
				status = "INVALID_TOKEN_TYPE"
				message = "The token provided is of a different type. Check the token and try again."
			case utils.ErrUserIDNotFound:
				status = "USER_ID_NOT_FOUND"
				message = "UserID not found in token."
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
