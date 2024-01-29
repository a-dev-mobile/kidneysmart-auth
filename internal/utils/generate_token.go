package utils

import (
	"errors"

	"time"

	"github.com/golang-jwt/jwt/v5"
	"strings"
)

var (
	ErrTokenExpired         = errors.New("token is expired")
	ErrInvalidToken         = errors.New("invalid token")
	ErrInvalidTokenType     = errors.New("invalid token type")
	ErrUserIDNotFound       = errors.New("userID not found in token")
	ErrTokenClaimsInvalid   = errors.New("token claims are invalid")
	ErrTokenSignatureInvalid = errors.New("token signature is invalid")
)

// CalculateAccessTokenExpiryTime возвращает время истечения access токена в UTC.
func CalculateAccessTokenExpiryTime(hours int) time.Time {
	return time.Now().UTC().Add(time.Duration(hours) * time.Hour)
}

// CalculateRefreshTokenExpiryTime возвращает время истечения refresh токена в UTC.
func CalculateRefreshTokenExpiryTime(days int) time.Time {
	return time.Now().UTC().Add(time.Duration(days) * 24 * time.Hour)
}

// GenerateAccessToken создает JWT access токен для верифицированного пользователя.
func GenerateAccessToken(userID string, jwtSecret string, accessTokenExpiryHours int) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		"type":   "access",
		// "exp":    CalculateAccessTokenExpiryTime(accessTokenExpiryHours).Unix(),
		"exp": time.Now().Add(time.Duration(1) * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

// GenerateRefreshToken создает JWT refresh токен для верифицированного пользователя.
func GenerateRefreshToken(userID string, jwtSecret string, refreshTokenExpiryDays int) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		"type":   "refresh",
		"exp":    CalculateRefreshTokenExpiryTime(refreshTokenExpiryDays).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

// ParseToken парсит и валидирует JWT, извлекая userID и проверяя срок действия.
func ParseToken(tokenString, jwtSecret, expectedTokenType string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenSignatureInvalid
		}
		return []byte(jwtSecret), nil
	})
// отлов встроеным методом что токен протух
	if err != nil {
		if strings.Contains(err.Error(), jwt.ErrTokenExpired.Error()) {
			return "", ErrTokenExpired
		}

		return "", ErrInvalidToken
	}

	if claims, ok := token.Claims.(*jwt.MapClaims); ok {
		// Проверка типа токена
		if tokenType, ok := (*claims)["type"].(string); !ok || tokenType != expectedTokenType {
			return "", ErrInvalidTokenType
		}

		if exp, ok := (*claims)["exp"].(float64); !ok || time.Now().UTC().Unix() > int64(exp) {
			return "", ErrTokenExpired
		}

		if userID, ok := (*claims)["userID"].(string); ok {
			return userID, nil
		}
		return "", ErrUserIDNotFound
	}

	return "", ErrTokenClaimsInvalid
}
