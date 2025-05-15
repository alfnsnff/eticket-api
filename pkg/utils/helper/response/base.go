package response

// BaseResponse is a generic structure for API responses.
type BaseResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`   // omitempty is important
	Errors  interface{} `json:"errors,omitempty"` // Use interface{} for flexibility
	Meta    interface{} `json:"meta,omitempty"`
}

type PaginationMeta struct {
	Total       int `json:"total"`
	PerPage     int `json:"per_page"`
	CurrentPage int `json:"current_page"`
	TotalPages  int `json:"total_pages"`
	NextPage    int `json:"next_page,omitempty"`
	PrevPage    int `json:"prev_page,omitempty"`
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

func NewPaginatedResponse(data interface{}, message string, total, perPage, currentPage int) *BaseResponse {
	totalPages := (total + perPage - 1) / perPage // ceil
	nextPage := 0
	if currentPage < totalPages {
		nextPage = currentPage + 1
	}
	prevPage := 0
	if currentPage > 1 {
		prevPage = currentPage - 1
	}

	meta := PaginationMeta{
		Total:       total,
		PerPage:     perPage,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
		NextPage:    nextPage,
		PrevPage:    prevPage,
	}

	return NewSuccessResponse(data, message, meta)
}
