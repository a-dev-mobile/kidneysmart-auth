package auth

// import (
// 	"auth-api/internal/config"
// 	"auth-api/internal/utils"
// 	pb "auth-api/proto"
// 	"context"

// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"golang.org/x/exp/slog"
// )

// type AuthServiceServer struct {
// 	pb.UnimplementedAuthServiceServer
// 	DB     *mongo.Client
// 	Config *config.Config
// 	Logger *slog.Logger
// }

// func NewAuthServiceServer(db *mongo.Client, cfg *config.Config, lg *slog.Logger) *AuthServiceServer {
// 	return &AuthServiceServer{
// 		DB:     db,
// 		Config: cfg,
// 		Logger: lg,
// 	}
// }

// func (s *AuthServiceServer) RegisterUser(ctx context.Context, req *pb.RegistrationRequest) (*pb.RegistrationResponse, error) {
// 	// Логирование запроса на регистрацию
// 	s.Logger.Info("RegisterUser called", slog.String("email", req.GetEmail()))

// 	collection := s.DB.Database(s.Config.Database.Name).Collection(s.Config.Database.Collections.Users)

// 	// Получение данных регистрации (только email).
// 	email := req.GetEmail()

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

// 	// Создание нового пользователя в базе данных (без пароля).
// 	_, err = collection.InsertOne(ctx, bson.M{"email": email})
// 	if err != nil {
// 		// Ошибка при добавлении пользователя.
// 		s.Logger.Error("Failed to create user", slog.String("email", email), slog.String("error", err.Error()))
// 		return nil, err
// 	}

// 	// Отправка подтверждающего URL на email
// 	subject := "Please confirm your email address"
// 	body := "Here is your confirmation URL: [URL]" // Сформируйте подтверждающий URL
// 	err = utils.SendEmail(email, subject, body)
// 	if err != nil {
// 		// Обработка ошибки отправки email
// 		return nil, err
// 	}

// 	// Отправка ответа об успешной регистрации.
// 	return &pb.RegistrationResponse{
// 		Status:  "success",
// 		Message: "User registered successfully, please complete registration process",
// 	}, nil
// }
