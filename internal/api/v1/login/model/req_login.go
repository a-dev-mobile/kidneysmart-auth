package model

import "github.com/go-playground/validator/v10"

type RequestLogin struct {
	// @Required
	Email string `json:"email" validate:"required,email"`
}
func (a *RequestLogin) Validate() error {
	validate := validator.New()
	return validate.Struct(a)
}
