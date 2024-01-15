package refreshtoken

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/a-dev-mobile/kidneysmart-auth/internal/api/v1/refresh_token/model"
	"github.com/a-dev-mobile/kidneysmart-auth/internal/config"
	"github.com/a-dev-mobile/kidneysmart-auth/internal/model/db"
	"github.com/a-dev-mobile/kidneysmart-auth/internal/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

type RefreshTokenServiceContext struct {
	DB     *mongo.Client
	Logger *slog.Logger
	Config *config.Config
}

func NewRefreshTokenServiceContext(db *mongo.Client, lg *slog.Logger, cfg *config.Config) *RefreshTokenServiceContext {
	return &RefreshTokenServiceContext{
		DB:     db,
		Config: cfg,
		Logger: lg,
	}
}

func (s *RefreshTokenServiceContext) RefreshTokenHandler(c *gin.Context) {

	var reqRefreshToken model.RequestRefreshToken
	if err := c.ShouldBindJSON(&reqRefreshToken); err != nil {
		s.Logger.Error("Failed to bind JSON", "error", err.Error(), "request", c.Request)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	if err := reqRefreshToken.Validate(); err != nil {
		s.Logger.Error("Validation error", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request parameters"})
		return
	}
	userID, err := utils.ParseToken(reqRefreshToken.RefreshToken, s.Config.Authentication.JWTSecret, "refresh")
	if err != nil {
		s.Logger.Error("Token validation error", "error", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
		return
	}

	// Найти, проверить и обновить существующий refresh токен
	newRefreshToken, err := s.validateAndUpdateRefreshToken(c.Request.Context(), reqRefreshToken.RefreshToken)
	if err != nil {
		s.Logger.Error("Refresh token validation/update error", "error", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	// Генерация нового access токена
	newAccessToken, err := utils.GenerateAccessToken(userID, s.Config.Authentication.JWTSecret, s.Config.Authentication.AccessTokenExpiryHours)
	if err != nil {
		s.Logger.Error("Failed to generate access token", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate access token"})
		return
	} 
	
	// Отправка новых токенов в ответе
	expiresIn := utils.CalculateAccessTokenExpiryTime(s.Config.Authentication.AccessTokenExpiryHours)
	c.JSON(http.StatusOK, model.ResponseRefreshToken{
		Message:      "Token refreshed successfully",
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    expiresIn,
	})
}

func (s *RefreshTokenServiceContext) validateAndUpdateRefreshToken(ctx context.Context, oldRefreshToken string) (string, error) {

	tokenCollectionName := s.Config.Database.Collections.AuthTokens
	collection := s.DB.Database(s.Config.Database.Name).Collection(tokenCollectionName)
	// Поиск существующего токена
	var existingToken db.AuthToken
	err := collection.FindOne(ctx, bson.M{"token": oldRefreshToken}).Decode(&existingToken)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", errors.New("refresh token not found")
		}
		return "", err
	}

	// Проверка, активен ли токен
	if !existingToken.IsActive {
		return "", errors.New("refresh token is not active")
	}

	// Генерация нового refresh токена
	newRefreshToken, err := utils.GenerateRefreshToken(existingToken.UserID.Hex(), s.Config.Authentication.JWTSecret, s.Config.Authentication.RefreshTokenExpiryDays)
	if err != nil {
		return "", err
	}

	// Обновление токена в базе данных
	update := bson.M{
		"$set": bson.M{
			"token":     newRefreshToken,
			"createdAt": time.Now(),
			"expiresAt": utils.CalculateRefreshTokenExpiryTime(s.Config.Authentication.RefreshTokenExpiryDays),
		},
	}
    _, err = collection.UpdateOne(ctx, bson.M{"token": oldRefreshToken}, update)
    if err != nil {
        return "", err
    }

	return newRefreshToken, nil
}
