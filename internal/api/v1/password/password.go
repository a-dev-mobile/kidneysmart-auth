package password

import (


	"net/http"

	"github.com/a-dev-mobile/kidneysmart-auth/internal/api/v1/password/model"
	"github.com/a-dev-mobile/kidneysmart-auth/internal/config"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

type PasswordServiceContext struct {
	DB          *mongo.Client
	Logger      *slog.Logger
	Config      *config.Config

}

func NewPasswordServiceContext(db *mongo.Client, lg *slog.Logger, cfg *config.Config, ) *PasswordServiceContext {
	return &PasswordServiceContext{
		DB:          db,
		Config:      cfg,
		Logger:      lg,
	
	}
}


func (s *PasswordServiceContext) PasswordHandler(c *gin.Context) {
	var reqPassword model.RequestPassword
    userID, exists := c.Get("userID")
	s.Logger.Debug("userID", "debug", userID)
	
    if !exists {
        // Обработка ошибки, если userID не найден
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

	if err := c.ShouldBindJSON(&reqPassword); err != nil {
		s.Logger.Error("Failed to bind JSON", "error", err.Error())
		c.JSON(http.StatusBadRequest, model.ResponsePassword{Message: "Invalid request body"})
		return
	}

	if err := reqPassword.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, model.ResponsePassword{Message: "Invalid request parameters"})
		return
	}



	c.JSON(http.StatusOK, model.ResponsePassword{Message: "User registered successfully"})
}
