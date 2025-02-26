package response

// BaseResponse is a generic structure for API responses.
type BaseResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`   // omitempty is important
	Errors  interface{} `json:"errors,omitempty"` // Use interface{} for flexibility
	Meta    interface{} `json:"meta,omitempty"`
}

// NewSuccessResponse creates a successful response.
func NewSuccessResponse(data interface{}, message string, meta interface{}) *BaseResponse {
	return &BaseResponse{
		Status:  true,
		Message: message,
		Data:    data,
		Meta:    meta,
	}
}

// NewErrorResponse creates an error response.
func NewErrorResponse(message string, errors interface{}) *BaseResponse {
	return &BaseResponse{
		Status:  false,
		Message: message,
		Errors:  errors,
	}
}
