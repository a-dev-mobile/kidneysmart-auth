package model

import "time"

// ResponseSuccessVerifyCode represents the response payload for verification.
type ResponseSuccessVerifyCode struct {
    // Message provides information about the response.

    Message string `json:"message"`

    // AccessToken is the JWT token for accessing secured endpoints.
   
    AccessToken string `json:"accessToken"`

    // RefreshToken is the JWT token used to refresh the access token.

    RefreshToken string `json:"refreshToken"`

    // ExpiresIn indicates the expiration time of the access token.
  
    ExpiresIn time.Time `json:"expiresIn"`
    
}