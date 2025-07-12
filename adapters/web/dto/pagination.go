package dto

// PaginationRequest represents pagination parameters from client requests
type PaginationRequest struct {
	Page  int `json:"page" form:"page" query:"page"`
	Limit int `json:"limit" form:"limit" query:"limit"`
}

// PaginationResponse represents pagination metadata in API responses
type PaginationResponse struct {
	CurrentPage int `json:"current_page"`
	PageSize    int `json:"page_size"`
	TotalItems  int `json:"total_items"`
	TotalPages  int `json:"total_pages"`
	HasNext     bool `json:"has_next"`
	HasPrev     bool `json:"has_prev"`
}
