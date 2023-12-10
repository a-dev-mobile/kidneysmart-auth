package model

// ResponseErrorVerifyCode represents the response payload for verification.
// This struct can be used for both successful and error responses.

type ResponseErrorVerifyCode struct {
	// Message provides information about the response.
	Message string `json:"message"`
}
