package response

// SuccessResponse structure for serializing successful responses.
type SuccessResponse struct {
 // Message contains an informational message about the successful operation.
    Message string `json:"message"`

// Data represents the data returned in the response. This field may be
    // of any type, which allows you to flexibly return different data in responses.
    Data interface{} `json:"data,omitempty"`
}
