package model

// ResponseLogin represents the response payload for a login request.
// @Description The response payload returned after a user tries to log in.
// @Object
type ResponseLogin struct {
    // Message provides information about the login outcome.
    // @Example "User registered successfully, verification code sent"
    // @Required This field must be provided in the response.
    Message string `json:"message"`

    // Status indicates the current stage or state of the login process.
    // @Example "REGISTRATION_SUCCESSFUL"
    // @Required This field must be provided in the response.
    // Possible values are:
    // - "INVALID_REQUEST_BODY": The request body is invalid.
    // - "INVALID_PARAMETERS": The request parameters are invalid.
    // - "INVALID_EMAIL_FORMAT": The provided email format is invalid.
    // - "INTERNAL_ERROR": An internal error occurred.
    // - "EMAIL_VERIFICATION_REQUIRED": Email verification is required.
    // - "PASSWORD_SET_REQUIRED": Setting a password is required.
    // - "PASSWORD_ENTRY_REQUIRED": Password entry is required.
    // - "USER_CREATION_FAILED": User creation failed.
    // - "EMAIL_SEND_FAILED": Sending email failed.
    // - "REGISTRATION_SUCCESSFUL": Registration was successful.
    Status string `json:"status"`
}
