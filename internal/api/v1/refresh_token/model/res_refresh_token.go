package model

import "time"

type ResponseRefreshToken struct {
	Message string `json:"message"`

	AccessToken string `json:"accessToken"`

	RefreshToken string `json:"refreshToken"`

	ExpiresIn time.Time `json:"expiresIn"`
}