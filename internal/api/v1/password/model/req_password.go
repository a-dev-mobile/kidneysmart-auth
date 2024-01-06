package model

import (
    "github.com/go-playground/validator/v10"
)

type RequestPassword struct {
    // @Required
    Password string `json:"password" validate:"required"`
}

func (a *RequestPassword) Validate() error {
    validate := validator.New()
    return validate.Struct(a)
}
