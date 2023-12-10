package utils

import (
    "time"

    "github.com/golang-jwt/jwt/v5"
)

// GenerateAccessToken creates a JWT access token for a verified user.
func GenerateAccessToken(email, jwtSecret string, accessTokenExpiryHours int) (string, error) {
    claims := jwt.MapClaims{
        "email": email,
        "exp":   time.Now().Add(time.Duration(accessTokenExpiryHours) * time.Hour).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte(jwtSecret))
    return tokenString, err
}

// GenerateRefreshToken creates a JWT refresh token for a verified user.
func GenerateRefreshToken(email, jwtSecret string, refreshTokenExpiryDays int) (string, error) {
    claims := jwt.MapClaims{
        "email": email,
        "exp":   time.Now().Add(time.Duration(refreshTokenExpiryDays) * 24 * time.Hour).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    refreshTokenString, err := token.SignedString([]byte(jwtSecret))
    return refreshTokenString, err
}
