package login

import (
	"context"
	"fmt"
	"net/http"

	"github.com/a-dev-mobile/kidneysmart-auth/internal/api/v1/login/model"
	"github.com/a-dev-mobile/kidneysmart-auth/internal/config"
	"github.com/a-dev-mobile/kidneysmart-auth/internal/model/db"
	"github.com/a-dev-mobile/kidneysmart-auth/internal/utils"
	"github.com/a-dev-mobile/kidneysmart-auth/pkg/emailclient"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

type LoginServiceContext struct {
	DB          *mongo.Client
	Logger      *slog.Logger
	Config      *config.Config
	EmailClient *emailclient.EmailClient
}

func NewLoginServiceContext(db *mongo.Client, lg *slog.Logger, cfg *config.Config, emailClient *emailclient.EmailClient) *LoginServiceContext {
	return &LoginServiceContext{
		DB:          db,
		Config:      cfg,
		Logger:      lg,
		EmailClient: emailClient,
	}
}
// LoginUserHandler handles the login of a user.
// @Summary Login a new user
// @Description This endpoint logs in a new user by their email address.
// If the email is not in the database, it registers the user and sends a verification code.
// If the email is in the database but not verified, it prompts for email verification.
// If the email is verified but no password is set, it prompts to set a password.
// If the email is verified and password is set, it prompts to enter the password.
// @Tags user
// @Accept json
// @Produce json
// @Param RequestLogin body model.RequestLogin true "Login Info"
// @Success 200 {object} model.ResponseLogin "User registered successfully, verification code sent"
// @Success 401 {object} model.ResponseLogin "Unauthorized - Email not verified or Password entry required"
// @Failure 400 {object} model.ResponseLogin "Invalid request body or parameters"
// @Failure 500 {object} model.ResponseLogin "Internal server error"
// @Router /login [post]
func (s *LoginServiceContext) LoginUserHandler(c *gin.Context) {
	var reqLogin model.RequestLogin

	if err := c.ShouldBindJSON(&reqLogin); err != nil {
		s.Logger.Error("Failed to bind JSON", "error", err.Error())
		c.JSON(http.StatusBadRequest, model.ResponseLogin{
			Message: "Invalid request body",
			Status:  "INVALID_REQUEST_BODY",
		})
		return
	}

	if err := reqLogin.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, model.ResponseLogin{
			Message: "Invalid request parameters",
			Status:  "INVALID_PARAMETERS",
		})
		return
	}

	if !utils.ValidateEmail(reqLogin.Email) {
		c.JSON(http.StatusBadRequest, model.ResponseLogin{
			Message: "Invalid email format",
			Status:  "INVALID_EMAIL_FORMAT",
		})
		return
	}

	authUserCollection := s.Config.Database.Collections[string(config.AuthUserCollection)]
	collection := s.DB.Database(s.Config.Database.Name).Collection(string(authUserCollection))

	ctx := c.Request.Context()

	userDetails, err := getUserDetails(ctx, collection, reqLogin.Email)
	if err != nil {
		s.Logger.Error("Failed to get user details", "email", reqLogin.Email, "error", err.Error())
		c.JSON(http.StatusInternalServerError, model.ResponseLogin{
			Message: "Failed to get user details",
			Status:  "INTERNAL_ERROR",
		})
		return
	}

	if userDetails != nil {
		if !userDetails.EmailVerified {
			c.JSON(http.StatusUnauthorized, model.ResponseLogin{
				Message: "Email not verified. Please verify your email.",
				Status:  "EMAIL_VERIFICATION_REQUIRED",
			})
			return
		} else if userDetails.Password == "" {
			c.JSON(http.StatusUnauthorized, model.ResponseLogin{
				Message: "Password not set. Please set your password.",
				Status:  "PASSWORD_SET_REQUIRED",
			})
			return
		} else {
			c.JSON(http.StatusUnauthorized, model.ResponseLogin{
				Message: "Please enter your password.",
				Status:  "PASSWORD_ENTRY_REQUIRED",
			})
			return
		}
	}

	code := utils.GenerateRandomCode()
	if err := createUser(ctx, collection, reqLogin.Email, code); err != nil {
		s.Logger.Error("Failed to create user", "email", reqLogin.Email, "error", err.Error())
		c.JSON(http.StatusInternalServerError, model.ResponseLogin{
			Message: "Failed to create user",
			Status:  "USER_CREATION_FAILED",
		})
		return
	}

	if err := sendConfirmationEmail(s.EmailClient, reqLogin.Email, code); err != nil {
		s.Logger.Warn("Failed to send email", "email", reqLogin.Email, "error", err.Error())
		c.JSON(http.StatusInternalServerError, model.ResponseLogin{
			Message: "User registered but failed to send confirmation email",
			Status:  "EMAIL_SEND_FAILED",
		})
		return
	}

	c.JSON(http.StatusOK, model.ResponseLogin{
		Message: "User registered successfully, verification code sent",
		Status:  "REGISTRATION_SUCCESSFUL",
	})
}

func getUserDetails(ctx context.Context, collection *mongo.Collection, email string) (*db.AuthUser, error) {
	filter := bson.M{"email": email}
	var existingUser db.AuthUser
	err := collection.FindOne(ctx, filter).Decode(&existingUser)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &existingUser, nil
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
	return client.SendEmail(email, subject, "KidneySmart", "hello@wayofdt.com", body)
}
