package verifycode

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/a-dev-mobile/kidneysmart-auth/internal/api/v1/verifycode/model"
	"github.com/a-dev-mobile/kidneysmart-auth/internal/config"
	"github.com/a-dev-mobile/kidneysmart-auth/internal/model/db"
	"github.com/a-dev-mobile/kidneysmart-auth/internal/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/exp/slog"
)

type VerifyCodeServiceContext struct {
	DB     *mongo.Client
	Logger *slog.Logger
	Config *config.Config
}

func NewVerifyCodeServiceContext(db *mongo.Client, lg *slog.Logger, cfg *config.Config) *VerifyCodeServiceContext {
	return &VerifyCodeServiceContext{
		DB:     db,
		Config: cfg,
		Logger: lg,
	}
}

// VerifyCodeHandler handles the verification of the code sent by the user.
// @Summary Verify User Code
// @Description Verifies the verification code sent by the user for account verification.
// @Tags verification
// @Accept json
// @Produce json
// @Param email query string true "Email address of the user"
// @Param code query string true "Verification code sent to the user's email"
// @Success 200 {object} model.ResponseSuccessVerifyCode "Verification successful, includes access and refresh tokens"
// @Success 208 {object} model.ResponseStatusVerifyCode "Email is already verified"
// @Failure 400 {object} model.ResponseStatusVerifyCode "Invalid request body or parameters"
// @Failure 401 {object} model.ResponseStatusVerifyCode "Invalid verification code"
// @Failure 404 {object} model.ResponseStatusVerifyCode "User not found"
// @Failure 429 {object} model.ResponseStatusVerifyCode "Too many attempts, please try again later"
// @Failure 500 {object} model.ResponseStatusVerifyCode "Internal server error"
// @Router /verifycode [post]
func (s *VerifyCodeServiceContext) VerifyCodeHandler(c *gin.Context) {
	var req model.RequestVerifyCode

	if err := c.ShouldBindJSON(&req); err != nil {
		s.Logger.Error("Failed to bind JSON", "error", err.Error())
		c.JSON(http.StatusBadRequest, model.ResponseVerifyCode{
			Message: "Invalid request body",
			Status:  "INVALID_REQUEST_BODY",
		})
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, model.ResponseVerifyCode{
			Message: "Invalid request parameters",
			Status:  "INVALID_PARAMETERS",
		})
		return
	}

	if err := validateRequest(req); err != nil {
		s.Logger.Info("Validation failed", "error", err.Error())
		c.JSON(http.StatusBadRequest, model.ResponseVerifyCode{
			Message: err.Error(),
			Status:  "VALIDATION_FAILED",
		})
		return
	}

	dbAuthUser, err := s.fetchUser(c.Request.Context(), req.Email)

	if errors.Is(err, mongo.ErrNoDocuments) {
		c.JSON(http.StatusNotFound, model.ResponseVerifyCode{
			Message: "User not found",
			Status:  "USER_NOT_FOUND",
		})

		return
	} else if err != nil {
		s.Logger.Error("Failed to retrieve user", "email", req.Email, "error", err.Error())
		c.JSON(http.StatusInternalServerError, model.ResponseVerifyCode{Message: "Error retrieving user"})
		return
	}
	// Check if the email is already verified
	if dbAuthUser.EmailVerified {
		var statusMessage, statusCode string

		// Check if the password is set for the user
		if dbAuthUser.Password == "" {
			// Email is verified but password is not set
			statusMessage = "Email is verified but password is not set"
			statusCode = "EMAIL_VERIFIED_PASSWORD_NOT_SET"
		} else {
			// Email is verified and password is set
			statusMessage = "Email and password are both verified"
			statusCode = "EMAIL_AND_PASSWORD_VERIFIED"
		}

		c.JSON(http.StatusAlreadyReported, model.ResponseVerifyCode{
			Message: statusMessage,
			Status:  statusCode,
		})
		return
	}
	// Check if the user has exceeded the maximum number of attempts
	const MaxAttempts = 5
	if dbAuthUser.AttemptCount >= MaxAttempts && time.Since(dbAuthUser.LastAttemptTime).Minutes() < 15 {

		c.JSON(http.StatusTooManyRequests, model.ResponseVerifyCode{
			Message: "Too many attempts, please try again later",
			Status:  "TOO_MANY_ATTEMPTS"})
		return
	}

	if dbAuthUser.Code != req.Code {
		s.incrementAttemptCount(c.Request.Context(), req.Email, dbAuthUser.AttemptCount)

		c.JSON(http.StatusUnauthorized, model.ResponseVerifyCode{
			Message: "Invalid code",
			Status:  "INVALID_CODE",
		})
		return
	}
	// Update the user's email verification status in the database
	if err := s.UpdateEmailVerificationStatus(c.Request.Context(), req.Email); err != nil {
		s.Logger.Error("Failed to update user email verification status", "email", req.Email, "error", err.Error())
		c.JSON(http.StatusInternalServerError, model.ResponseVerifyCode{
			Message: "Error updating user verification status",
			Status:  "UPDATE_VERIFICATION_STATUS_FAILED",
		})
		return
	}

	// Reset the attempt count on successful verification
	s.resetAttemptCount(c.Request.Context(), req.Email)

	// Generate a token for the verified user
	accessToken, err := utils.GenerateAccessToken(dbAuthUser.ID.Hex(), s.Config.Authentication.JWTSecret, s.Config.Authentication.AccessTokenExpiryHours)
	if err != nil {
		s.Logger.Error("Failed to generate access token", "error", err.Error())
		c.JSON(http.StatusInternalServerError, model.ResponseVerifyCode{
			Message: "Failed to generate access token",
			Status:  "ACCESS_TOKEN_GENERATION_FAILED",
		})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(dbAuthUser.ID.Hex(), s.Config.Authentication.JWTSecret, s.Config.Authentication.RefreshTokenExpiryDays)
	if err != nil {
		s.Logger.Error("Failed to generate refresh token", "error", err.Error())
		c.JSON(http.StatusInternalServerError, model.ResponseVerifyCode{
			Message: "Failed to generate refresh token",
			Status:  "REFRESH_TOKEN_GENERATION_FAILED",
		})
		return
	}

	// Save the refresh token in the database
	if err := s.SaveRefreshToken(c.Request.Context(), dbAuthUser.ID, refreshToken); err != nil {
		s.Logger.Error("Failed to save refresh token", "userID", dbAuthUser.ID.Hex(), "error", err.Error())
		c.JSON(http.StatusInternalServerError, model.ResponseVerifyCode{
			Message: "Error saving refresh token",
			Status:  "REFRESH_TOKEN_SAVING_FAILED",
		})
		return
	}
	// Generate the success response
	expiresIn := utils.CalculateAccessTokenExpiryTime(s.Config.Authentication.AccessTokenExpiryHours)
	successResponse := model.ResponseVerifyCode{
		Message:      "Verification successful",
		Status:       "VERIFICATION_SUCCESSFUL",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    &expiresIn,
	}

	c.JSON(http.StatusOK, successResponse)
}

func validateRequest(req model.RequestVerifyCode) error {
	if !utils.ValidateEmail(req.Email) {
		return errors.New("invalid email format")
	}
	if !utils.ValidateCode(req.Code) {
		return errors.New("invalid code format")
	}
	return nil
}

func (s *VerifyCodeServiceContext) fetchUser(ctx context.Context, email string) (*db.AuthUser, error) {
	authUserCollection := s.Config.Database.Collections.AuthUser
	collection := s.DB.Database(s.Config.Database.Name).Collection(authUserCollection)
	var dbAuthUser db.AuthUser
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&dbAuthUser)
	return &dbAuthUser, err
}
func (s *VerifyCodeServiceContext) UpdateEmailVerificationStatus(ctx context.Context, email string) error {
	authUserCollection := s.Config.Database.Collections.AuthUser
	collection := s.DB.Database(s.Config.Database.Name).Collection(authUserCollection)
	update := bson.M{"$set": bson.M{"emailVerified": true}}
	_, err := collection.UpdateOne(ctx, bson.M{"email": email}, update)
	return err
}

func (s *VerifyCodeServiceContext) incrementAttemptCount(ctx context.Context, email string, currentCount int) {
	authUserCollection := s.Config.Database.Collections.AuthUser
	collection := s.DB.Database(s.Config.Database.Name).Collection(authUserCollection)
	update := bson.M{"$set": bson.M{"attemptCount": currentCount + 1, "lastAttemptTime": time.Now()}}
	_, _ = collection.UpdateOne(ctx, bson.M{"email": email}, update)
}

func (s *VerifyCodeServiceContext) resetAttemptCount(ctx context.Context, email string) {
	authUserCollection := s.Config.Database.Collections.AuthUser
	collection := s.DB.Database(s.Config.Database.Name).Collection(authUserCollection)
	update := bson.M{"$set": bson.M{"attemptCount": 0, "lastAttemptTime": time.Time{}}}
	_, _ = collection.UpdateOne(ctx, bson.M{"email": email}, update)
}

// SaveRefreshToken сохраняет refresh токен в отдельной коллекции AuthToken.
func (s *VerifyCodeServiceContext) SaveRefreshToken(ctx context.Context, userID primitive.ObjectID, refreshToken string) error {
	// Название коллекции токенов
	tokenCollectionName := s.Config.Database.Collections.AuthToken
	collection := s.DB.Database(s.Config.Database.Name).Collection(tokenCollectionName)

	// Создание объекта AuthToken
	authToken := db.AuthToken{
		UserID:       userID,
		DeviceInfoID: primitive.NilObjectID,
		Token:        refreshToken,
		CreatedAt:    time.Now(),
		ExpiresAt:    utils.CalculateRefreshTokenExpiryTime(s.Config.Authentication.RefreshTokenExpiryDays),
		IsActive:     true,
	}

	// Вставка объекта AuthToken в коллекцию
	_, err := collection.InsertOne(ctx, authToken)
	return err
}
