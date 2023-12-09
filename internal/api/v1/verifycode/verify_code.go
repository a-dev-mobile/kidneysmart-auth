package verifycode

import (
	"context"
	"errors"
	"net/http"
	"time"

	modelDb "github.com/a-dev-mobile/kidneysmart-auth/internal/model/db"
	"github.com/a-dev-mobile/kidneysmart-auth/internal/config"
	modelApi "github.com/a-dev-mobile/kidneysmart-auth/internal/api/v1/verifycode/model"
	"github.com/a-dev-mobile/kidneysmart-auth/internal/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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
func (s *VerifyCodeServiceContext) VerifyCodeHandler(c *gin.Context) {
	var req modelApi.RequestVerifyCode

	if err := c.ShouldBindJSON(&req); err != nil {
		s.Logger.Error("Failed to bind JSON", "error", err.Error())
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := req.Validate(); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request parameters")
		return
	}

	if err := validateRequest(req); err != nil {
		s.Logger.Info("Validation failed", slog.String("error", err.Error()))
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	dbAuthUser, err := s.fetchUser(c.Request.Context(), req.Email)
	if errors.Is(err, mongo.ErrNoDocuments) {
		utils.RespondWithError(c, http.StatusNotFound, "User not found")
		return
	} else if err != nil {
		s.Logger.Error("Failed to retrieve user", "email", req.Email, "error", err.Error())
		utils.RespondWithError(c, http.StatusInternalServerError, "Error retrieving user")
		return
	}

	// Check if the user has exceeded the maximum number of attempts
	const MaxAttempts = 5
	if dbAuthUser.AttemptCount >= MaxAttempts && time.Since(dbAuthUser.LastAttemptTime).Minutes() < 15 {
		utils.RespondWithError(c, http.StatusTooManyRequests, "Too many attempts, please try again later")
		return
	}

	if dbAuthUser.Code != req.Code {
		s.incrementAttemptCount(c.Request.Context(), req.Email, dbAuthUser.AttemptCount)
		utils.RespondWithError(c, http.StatusUnauthorized, "Invalid code")
		return
	}
	// Update the user's email verification status in the database
	if err := s.UpdateEmailVerificationStatus(c.Request.Context(), req.Email); err != nil {
		s.Logger.Error("Failed to update user email verification status", "email", req.Email, "error", err.Error())
		utils.RespondWithError(c, http.StatusInternalServerError, "Error updating user verification status")
		return
	}

	// Reset the attempt count on successful verification
	s.resetAttemptCount(c.Request.Context(), req.Email)

	utils.RespondWithSuccess(c, http.StatusOK, "Code verified successfully", nil)
}

func validateRequest(req modelApi.RequestVerifyCode) error {
	if !utils.ValidateEmail(req.Email) {
		return errors.New("invalid email format")
	}
	if !utils.ValidateCode(req.Code) {
		return errors.New("invalid code format")
	}
	return nil
}

func (s *VerifyCodeServiceContext) fetchUser(ctx context.Context, email string) (*modelDb.AuthUser, error) {
	authUserCollection := s.Config.Database.Collections[string(config.AuthUserCollection)]
	collection := s.DB.Database(s.Config.Database.Name).Collection(string(authUserCollection))
	var dbAuthUser modelDb.AuthUser
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&dbAuthUser)
	return &dbAuthUser, err
}
func (s *VerifyCodeServiceContext) UpdateEmailVerificationStatus(ctx context.Context, email string) error {
	collection := s.DB.Database(s.Config.Database.Name).Collection(string(config.AuthUserCollection))
	update := bson.M{"$set": bson.M{"emailVerified": true}}
	_, err := collection.UpdateOne(ctx, bson.M{"email": email}, update)
	return err
}

func (s *VerifyCodeServiceContext) incrementAttemptCount(ctx context.Context, email string, currentCount int) {
	collection := s.DB.Database(s.Config.Database.Name).Collection(string(config.AuthUserCollection))
	update := bson.M{"$set": bson.M{"attemptCount": currentCount + 1, "lastAttemptTime": time.Now()}}
	_, _ = collection.UpdateOne(ctx, bson.M{"email": email}, update)
}

func (s *VerifyCodeServiceContext) resetAttemptCount(ctx context.Context, email string) {
	collection := s.DB.Database(s.Config.Database.Name).Collection(string(config.AuthUserCollection))
	update := bson.M{"$set": bson.M{"attemptCount": 0, "lastAttemptTime": time.Time{}}}
	_, _ = collection.UpdateOne(ctx, bson.M{"email": email}, update)
}
