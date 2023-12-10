package register

import (
	"context"
	"fmt"
	"net/http"

	"github.com/a-dev-mobile/kidneysmart-auth/pkg/emailclient"
	"github.com/a-dev-mobile/kidneysmart-auth/internal/config"
	"github.com/a-dev-mobile/kidneysmart-auth/internal/api/v1/register/model"
	"github.com/a-dev-mobile/kidneysmart-auth/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

type RegisterServiceContext struct {
	DB          *mongo.Client
	Logger      *slog.Logger
	Config      *config.Config
	EmailClient *emailclient.EmailClient
}

func NewRegisterServiceContext(db *mongo.Client, lg *slog.Logger, cfg *config.Config, emailClient *emailclient.EmailClient) *RegisterServiceContext {
	return &RegisterServiceContext{
		DB:          db,
		Config:      cfg,
		Logger:      lg,
		EmailClient: emailClient,
	}
}

func (s *RegisterServiceContext) RegisterUserHandler(c *gin.Context) {
	var reqRegister model.RequestRegister

	if err := c.ShouldBindJSON(&reqRegister); err != nil {
		s.Logger.Error("Failed to bind JSON", "error", err.Error())
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := reqRegister.Validate(); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request parameters")
		return
	}

	if !utils.ValidateEmail(reqRegister.Email) {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid email format")
		return
	}

	authUserCollection := s.Config.Database.Collections[string(config.AuthUserCollection)]
	collection := s.DB.Database(s.Config.Database.Name).Collection(string(authUserCollection))

	ctx := c.Request.Context()
	// Generate a random 4-digit code
	code := utils.GenerateRandomCode()
	// Check if user already exists and create new user with the code
	if userExists(ctx, collection, reqRegister.Email) {
		utils.RespondWithError(c, http.StatusConflict, "User already exists")
		return
	}

	if err := createUser(ctx, collection, reqRegister.Email, code); err != nil {
		s.Logger.Error("Failed to create user", "email", reqRegister.Email, "error", err.Error())
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

// Send the code via email
if err := sendConfirmationEmail(s.EmailClient, reqRegister.Email, code); err != nil {
    s.Logger.Warn("Failed to send email", "email", reqRegister.Email, "error", err.Error())
    utils.RespondWithError(c, http.StatusInternalServerError, "User registered but failed to send confirmation email")
    return
}

	utils.RespondWithSuccess(c, http.StatusOK, "User registered successfully", nil)
}

func userExists(ctx context.Context, collection *mongo.Collection, email string) bool {
	filter := bson.M{"email": email}
	var existingUser struct{}
	err := collection.FindOne(ctx, filter).Decode(&existingUser)
	return err == nil || err != mongo.ErrNoDocuments
}

func createUser(ctx context.Context, collection *mongo.Collection, email, code string) error {
	newUser := bson.M{
		"email": email,
		"code":  code,
	}
	_, err := collection.InsertOne(ctx, newUser)
	return err
}

func sendConfirmationEmail(client *emailclient.EmailClient, email string, code string) error {
	subject :=  fmt.Sprintf("Your verification code is: %s", code)
	body := fmt.Sprintf("%s \nPlease use this code to complete your registration.", code)
	return client.SendEmail(email, subject, "KidneySmart Team", "hello@wayofdt.com", body)
}
