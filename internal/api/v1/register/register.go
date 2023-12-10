package register

import (
	"context"
	"fmt"
	"net/http"

	"github.com/a-dev-mobile/kidneysmart-auth/internal/api/v1/register/model"
	"github.com/a-dev-mobile/kidneysmart-auth/internal/config"
	"github.com/a-dev-mobile/kidneysmart-auth/internal/utils"
	"github.com/a-dev-mobile/kidneysmart-auth/pkg/emailclient"
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

// RegisterUserHandler registers a new user.
// @Summary Register a new user
// @Description This endpoint registers a new user by their email address and sends a verification code to that email.
// @Tags user
// @Accept json
// @Produce json
// @Param RequestRegister body model.RequestRegister true "Registration Info"
// @Success 200 {object} model.ResponseRegister "User registered successfully, verification code sent"
// @Failure 400 {object} model.ResponseRegister "Invalid request body or parameters, such as incorrect email format"
// @Failure 409 {object} model.ResponseRegister "User already exists with the provided email address"
// @Failure 500 {object} model.ResponseRegister "Internal server error, such as failure in user creation or sending email"
// @Router /register [post]
func (s *RegisterServiceContext) RegisterUserHandler(c *gin.Context) {
	var reqRegister model.RequestRegister

	if err := c.ShouldBindJSON(&reqRegister); err != nil {
		s.Logger.Error("Failed to bind JSON", "error", err.Error())
		c.JSON(http.StatusBadRequest, model.ResponseRegister{Message: "Invalid request body"})
		return
	}

	if err := reqRegister.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, model.ResponseRegister{Message: "Invalid request parameters"})
		return
	}

	if !utils.ValidateEmail(reqRegister.Email) {
		c.JSON(http.StatusBadRequest, model.ResponseRegister{Message: "Invalid email format"})
		return
	}

	authUserCollection := s.Config.Database.Collections[string(config.AuthUserCollection)]
	collection := s.DB.Database(s.Config.Database.Name).Collection(string(authUserCollection))

	ctx := c.Request.Context()
	// Generate a random 4-digit code
	code := utils.GenerateRandomCode()
	// Check if user already exists and create new user with the code
	if userExists(ctx, collection, reqRegister.Email) {
		c.JSON(http.StatusConflict, model.ResponseRegister{Message: "User already exists"})
		return
	}

	if err := createUser(ctx, collection, reqRegister.Email, code); err != nil {
		s.Logger.Error("Failed to create user", "email", reqRegister.Email, "error", err.Error())
		c.JSON(http.StatusInternalServerError, model.ResponseRegister{Message: "Failed to create user"})
		return
	}

	// Send the code via email
	if err := sendConfirmationEmail(s.EmailClient, reqRegister.Email, code); err != nil {
		s.Logger.Warn("Failed to send email", "email", reqRegister.Email, "error", err.Error())
		c.JSON(http.StatusInternalServerError, model.ResponseRegister{Message: "User registered but failed to send confirmation email"})
		return
	}

	c.JSON(http.StatusOK, model.ResponseRegister{Message: "User registered successfully"})
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
	subject := fmt.Sprintf("Your verification code is: %s", code)
	body := fmt.Sprintf("%s \nPlease use this code to complete your registration.", code)
	return client.SendEmail(email, subject, "KidneySmart Team", "hello@wayofdt.com", body)
}
