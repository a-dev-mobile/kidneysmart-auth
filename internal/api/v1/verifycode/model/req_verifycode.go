package model

import "github.com/go-playground/validator/v10"

type RequestVerifyCode struct {
	// @Required
	Email string `json:"email" validate:"required"`
	// @Required
	Code string `json:"code" validate:"required"`
}

func (a *RequestVerifyCode) Validate() error {
	validate := validator.New()
	return validate.Struct(a)
}
