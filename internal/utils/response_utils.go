package utils

import (
    "github.com/gin-gonic/gin"

    "github.com/a-dev-mobile/kidneysmart-auth/internal/model/http/response"
)

// RespondWithError sends a standardized error response to the client.
func RespondWithError(c *gin.Context, statusCode int, message string) {
    c.JSON(statusCode, response.ErrorResponse{Message: message})
}

// RespondWithSuccess sends a standardized success response to the client.
func RespondWithSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
  // Create a SuccessResponse structure and pass data and a message to it
    successResponse := response.SuccessResponse{
        Message: message,
        Data:    data,
    }

// Send a response to the client
    c.JSON(statusCode, successResponse)
}