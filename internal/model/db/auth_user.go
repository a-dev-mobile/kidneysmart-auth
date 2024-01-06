package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuthUser represents a user in the authUser collection.
type AuthUser struct {
	ID primitive.ObjectID `bson:"_id"`

	Email           string    `json:"email" bson:"email"`
	Code            string    `json:"code" bson:"code"`
	EmailVerified   bool      `json:"emailVerified" bson:"emailVerified"`
	AttemptCount    int       `json:"attemptCount" bson:"attemptCount"`
	LastAttemptTime time.Time `json:"lastAttemptTime" bson:"lastAttemptTime"`
	RefreshToken    string    `json:"refreshToken" bson:"refreshToken"`
	Password        string    `json:"password" bson:"password"`
}
