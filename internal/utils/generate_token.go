package utils

import (
    "time"

    "github.com/golang-jwt/jwt/v5"
)

const (
    AccessTokenExpiryDuration  = 15 * time.Minute  
    RefreshTokenExpiryDuration = 72 * time.Hour   
)

// GenerateAccessToken creates a JWT access token for a verified user.
func GenerateAccessToken(email, jwtSecret string) (string, error) {
    claims := jwt.MapClaims{
        "email": email,
        "exp":   time.Now().Add(AccessTokenExpiryDuration).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte(jwtSecret))
    return tokenString, err
}

// GenerateRefreshToken creates a JWT refresh token for a verified user.
func GenerateRefreshToken(email, jwtSecret string) (string, error) {
    claims := jwt.MapClaims{
        "email": email,
        "exp":   time.Now().Add(RefreshTokenExpiryDuration).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    refreshTokenString, err := token.SignedString([]byte(jwtSecret))
    return refreshTokenString, err
}