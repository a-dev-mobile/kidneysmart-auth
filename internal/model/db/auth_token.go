package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthToken struct {
	// ID           primitive.ObjectID `bson:"_id"`
	UserID       primitive.ObjectID `bson:"userId"`       // Ссылка на идентификатор пользователя
	DeviceInfoID primitive.ObjectID `bson:"deviceInfoId"` // Ссылка на запись в таблице deviceInfo
	Token        string             `bson:"token"`        // Сам Refresh Token
	CreatedAt    time.Time          `bson:"createdAt"`    // Время создания токена
	ExpiresAt    time.Time          `bson:"expiresAt"`    // Время истечения срока действия токена
	IsActive     bool               `bson:"isActive"`     // Статус активности токена
}
