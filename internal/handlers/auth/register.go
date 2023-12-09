package auth

import (

	"net/http"
	//
	"github.com/a-dev-mobile/kidneysmart-auth/internal/config"

	"github.com/gin-gonic/gin"

	"golang.org/x/exp/slog"
)

//
type AuthServiceContext struct {


	Config *config.Config
	Logger *slog.Logger
}
// 
func NewAuthServiceContext( cfg *config.Config, lg *slog.Logger) *AuthServiceContext {




	return &AuthServiceContext{
	
		Config: cfg,
		Logger: lg,
	}
}
// 
func (hctx *AuthServiceContext) RegisterUser(c *gin.Context) {
	hctx.Logger.Debug("Start RegisterUser handler")

	hctx.Logger.Debug("End App handler", "StatusOK", "OK")
	c.JSON(http.StatusOK, "OK")

// // 
// 	collection := s.DB.Database(s.Config.Database.Name).Collection(s.Config.Database.Collections.Users)
// // 
// 	// Получение данных регистрации (только email).
// 	email := req.GetEmail()
// // 
// 	// Валидация email
// 	if !utils.ValidateEmail(email) {
// 		return &pb.RegistrationResponse{
// 			Status:  "invalid_email",
// 			Message: "Invalid email format",
// 		}, nil
// 	}
// 	// Проверка, существует ли уже пользователь с таким email.
// 	filter := bson.M{"email": email}
// 	var existingUser struct{}
// 	err := collection.FindOne(ctx, filter).Decode(&existingUser)
// 	if err != mongo.ErrNoDocuments {
// 		// Пользователь уже существует.
// 		return &pb.RegistrationResponse{
// 			Status:  "user_exists",
// 			Message: "User with this email already exists",
// 		}, nil
// 	}
// // 
// 	// Создание нового пользователя в базе данных (без пароля).
// 	_, err = collection.InsertOne(ctx, bson.M{"email": email})
// 	if err != nil {
// 		// Ошибка при добавлении пользователя.
// 		s.Logger.Error("Failed to create user", slog.String("email", email), slog.String("error", err.Error()))
// 		return nil, err
// 	}
// // 
// 	// Отправка подтверждающего URL на email
// 	subject := "Please confirm your email address"
// 	body := "Here is your confirmation URL: [URL]" // Сформируйте подтверждающий URL
// 	err = utils.SendEmail(email, subject, body)
// 	if err != nil {
// 		// Обработка ошибки отправки email
// 		return nil, err
// 	}
// 

}
// 