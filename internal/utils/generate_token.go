package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
		"exp":    CalculateAccessTokenExpiryTime(accessTokenExpiryHours).Unix(),
		// "exp":    time.Now().Add(time.Duration(1) * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	return tokenString, err
}

// GenerateRefreshToken создает JWT refresh токен для верифицированного пользователя.
func GenerateRefreshToken(userID string, jwtSecret string, refreshTokenExpiryDays int) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		"type":   "refresh",
		"exp":    CalculateRefreshTokenExpiryTime(refreshTokenExpiryDays).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshTokenString, err := token.SignedString([]byte(jwtSecret))
	return refreshTokenString, err
}

// ParseToken парсит и валидирует JWT, извлекая userID.
// Тип токена указывается как "access" или "refresh".
func ParseToken(tokenString string, jwtSecret string, expectedTokenType string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверка метода подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Проверка типа токена
		if tokenType, ok := claims["type"].(string); !ok || tokenType != expectedTokenType {
			return "", errors.New("invalid token type")
		}

		userID, ok := claims["userID"].(string)
		if !ok {
			return "", errors.New("userID not found in token")
		}
		return userID, nil
	}

	return "", errors.New("invalid token")
}
