package model

import (
    "github.com/go-playground/validator/v10"
)

type RequestRefreshToken struct {
    // @Required
    RefreshToken string `json:"refreshToken" validate:"required"`
}

func (a *RequestRefreshToken) Validate() error {
    validate := validator.New()
    return validate.Struct(a)
}
