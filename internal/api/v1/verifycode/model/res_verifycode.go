package model

import (
	"encoding/json"
	"time"
)

// ResponseVerifyCode represents the response payload for verification.
// This struct can be used for both successful and error responses.
// @Description Response payload for the verification code process.

type ResponseVerifyCode struct {
	// Message provides information about the response.
	// It describes the outcome of the verification process.
	Message string `json:"message"`

	// Status is a string that represents the specific status of the verification process.
	// Possible values include:
	// - "INVALID_REQUEST_BODY" for invalid request bodies.
	// - "INVALID_PARAMETERS" for invalid request parameters.
	// - "VALIDATION_FAILED" for failed validation of the request data.
	// - "USER_NOT_FOUND" if the user's email is not found in the database.
	// - "EMAIL_ALREADY_VERIFIED" if the user's email is already verified.
	// - "INVALID_CODE" for incorrect verification codes.
	// - "UPDATE_VERIFICATION_STATUS_FAILED" if there was an error updating the user's verification status.
	// - "ACCESS_TOKEN_GENERATION_FAILED" if there was an error generating the access token.
	// - "REFRESH_TOKEN_GENERATION_FAILED" if there was an error generating the refresh token.
	// - "REFRESH_TOKEN_SAVING_FAILED" if there was an error saving the refresh token.
	// - "TOO_MANY_ATTEMPTS"
	// - "VERIFICATION_SUCCESSFUL"
	Status string `json:"status"`
	// AccessToken is the JWT token for accessing secured endpoints.

	AccessToken string `json:"accessToken,omitempty"`

	// RefreshToken is the JWT token used to refresh the access token.

	RefreshToken string `json:"refreshToken,omitempty"`

	// ExpiresIn indicates the expiration time of the access token.

	ExpiresIn    *time.Time `json:"-"`
}
// Custom MarshalJSON to handle ExpiresIn.
func (r ResponseVerifyCode) MarshalJSON() ([]byte, error) {
    type Alias ResponseVerifyCode

    if r.ExpiresIn != nil && !r.ExpiresIn.IsZero() {
        return json.Marshal(&struct {
            ExpiresIn time.Time `json:"expiresIn"`
            *Alias
        }{
            ExpiresIn: *r.ExpiresIn,
            Alias:     (*Alias)(&r),
        })
    }

    return json.Marshal(&struct {
        *Alias
    }{
        Alias: (*Alias)(&r),
    })
}