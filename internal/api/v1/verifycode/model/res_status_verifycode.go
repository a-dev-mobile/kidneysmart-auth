package model

// ResponseStatusVerifyCode represents the response payload for verification.
// This struct can be used for both successful and error responses.
// @Description Response payload for the verification code process.
// @Success 200 {object} ResponseStatusVerifyCode "Verification successful"
// @Failure 400 {object} ResponseStatusVerifyCode "Invalid request body or parameters"
// @Failure 401 {object} ResponseStatusVerifyCode "Invalid verification code"
// @Failure 404 {object} ResponseStatusVerifyCode "User not found"
// @Failure 429 {object} ResponseStatusVerifyCode "Too many attempts, please try again later"
// @Failure 500 {object} ResponseStatusVerifyCode "Internal server error"
type ResponseStatusVerifyCode struct {
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
	Status string `json:"status"`
}
