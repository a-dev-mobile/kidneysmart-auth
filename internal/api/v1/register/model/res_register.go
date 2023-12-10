package model

// ResponseRegister represents the response payload for a registration request.
// @Description The response payload returned after a user registers.
type ResponseRegister struct {
    // Message provides information about the registration outcome.
    // @Required This field must be provided in the response.
    Message string `json:"message"`
}
