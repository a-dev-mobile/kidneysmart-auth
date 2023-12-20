package model

// ResponseStatusVerifyCode represents the response payload for verification.
// This struct can be used for both successful and error responses.

type ResponseStatusVerifyCode struct {
	// Message provides information about the response.
	Message string `json:"message"`
}
