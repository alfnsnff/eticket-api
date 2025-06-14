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

func NewErrorResponse(message string, errors interface{}) *BaseResponse {
	return &BaseResponse{
		Status:  false,
		Message: message,
		Errors:  errors,
	}
}

type Meta struct {
	Total       int    `json:"total"`               // Total items
	PerPage     int    `json:"per_page"`            // Items per page
	CurrentPage int    `json:"current_page"`        // Current page number
	TotalPages  int    `json:"total_pages"`         // Total number of pages
	NextPage    int    `json:"next_page,omitempty"` // Next page number if available
	PrevPage    int    `json:"prev_page,omitempty"` // Previous page number if available
	From        int    `json:"from"`                // Start index of items in this page (1-based)
	To          int    `json:"to"`                  // End index of items in this page
	Sort        string `json:"sort,omitempty"`      // Sort parameter (e.g., "name:asc")
	Search      string `json:"search,omitempty"`    // Search keyword
	Path        string `json:"path,omitempty"`      // Current path for building links
	HasNextPage bool   `json:"has_next_page"`       // Whether there is a next page
	HasPrevPage bool   `json:"has_prev_page"`       // Whether there is a previous page
}

func NewMetaResponse(
	data interface{},
	message string,
	total, perPage, currentPage int,
	sort, search, path string,
) *BaseResponse {
	totalPages := (total + perPage - 1) / perPage // ceil
	nextPage := 0
	if currentPage < totalPages {
		nextPage = currentPage + 1
	}
	prevPage := 0
	if currentPage > 1 {
		prevPage = currentPage - 1
	}

	from := (currentPage-1)*perPage + 1
	to := currentPage * perPage
	if to > total {
		to = total
	}

	meta := Meta{
		Total:       total,
		PerPage:     perPage,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
		NextPage:    nextPage,
		PrevPage:    prevPage,
		From:        from,
		To:          to,
		Sort:        sort,
		Search:      search,
		Path:        path,
		HasNextPage: currentPage < totalPages,
		HasPrevPage: currentPage > 1,
	}

	return NewSuccessResponse(data, message, meta)
}
