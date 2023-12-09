package db

import "time"

// AuthUser represents a user in the authUser collection.
type AuthUser struct {
	Email           string    `json:"email" bson:"email"`
	Code            string    `json:"code" bson:"code"`
	EmailVerified   bool      `json:"emailVerified" bson:"emailVerified"`
	AttemptCount    int       `json:"attemptCount" bson:"attemptCount"`
	LastAttemptTime time.Time `json:"lastAttemptTime" bson:"lastAttemptTime"`
}
