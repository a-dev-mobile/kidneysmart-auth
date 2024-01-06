package login
/* 
Вход в Систему:

Обработчик для входа пользователя в систему.
Принимает email и пароль.
Проверяет учетные данные и создает сессию для пользователя.
Пример названия файла: login.go.
 */




import (
	"context"
	
	"net/http"

	"github.com/a-dev-mobile/kidneysmart-auth/internal/api/v1/login/model"
	"github.com/a-dev-mobile/kidneysmart-auth/internal/config"
	"github.com/a-dev-mobile/kidneysmart-auth/internal/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

type LoginServiceContext struct {
	DB          *mongo.Client
	Logger      *slog.Logger
	Config      *config.Config

}

func NewLoginServiceContext(db *mongo.Client, lg *slog.Logger, cfg *config.Config, ) *LoginServiceContext {
	return &LoginServiceContext{
		DB:          db,
		Config:      cfg,
		Logger:      lg,
	
	}
}


func (s *LoginServiceContext) LoginUserHandler(c *gin.Context) {
	var reqLogin model.RequestLogin

	if err := c.ShouldBindJSON(&reqLogin); err != nil {
		s.Logger.Error("Failed to bind JSON", "error", err.Error())
		c.JSON(http.StatusBadRequest, model.ResponseLogin{Message: "Invalid request body"})
		return
	}

	if err := reqLogin.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, model.ResponseLogin{Message: "Invalid request parameters"})
		return
	}

	if !utils.ValidateEmail(reqLogin.Email) {
		c.JSON(http.StatusBadRequest, model.ResponseLogin{Message: "Invalid email format"})
		return
	}

	authUserCollection := s.Config.Database.Collections[string(config.AuthUserCollection)]
	collection := s.DB.Database(s.Config.Database.Name).Collection(string(authUserCollection))

	ctx := c.Request.Context()
	// Generate a random 4-digit code
	
	// Check if user already exists and create new user with the code
	if userExists(ctx, collection, reqLogin.Email) {
		c.JSON(http.StatusConflict, model.ResponseLogin{Message: "User already exists"})
		return
	}


	c.JSON(http.StatusOK, model.ResponseLogin{Message: "User registered successfully"})
}

func userExists(ctx context.Context, collection *mongo.Collection, email string) bool {
	filter := bson.M{"email": email}
	var existingUser struct{}
	err := collection.FindOne(ctx, filter).Decode(&existingUser)
	return err == nil || err != mongo.ErrNoDocuments
}
