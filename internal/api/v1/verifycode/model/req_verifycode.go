package model

import "github.com/go-playground/validator/v10"


// RequestVerifyCode represents the request payload for code verification.
// @Description The request payload needed for verifying a user's code.
type RequestVerifyCode struct {
	// Email is the user's email address that needs to be verified.
	// @Required This field must be provided for the request to be valid.
	Email string `json:"email" validate:"required,email"`

	// Code is the verification code sent to the user's email.
	// @Required This field must be provided for the request to be valid.
	Code  string `json:"code" validate:"required"`
}

// Validate performs validation on the RequestVerifyCode fields.
func (a *RequestVerifyCode) Validate() error {
	validate := validator.New()
	return validate.Struct(a)
}