package model

import "github.com/go-playground/validator/v10"

type RequestRegister struct {
	// @Required
	Email string `json:"email" validate:"required,email"`
}
func (a *RequestRegister) Validate() error {
	validate := validator.New()
	return validate.Struct(a)
}
