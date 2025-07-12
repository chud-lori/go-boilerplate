package dto

type WebResponse struct {
	Message string      `json:"message"`
	Status  int         `json:"status"`
	Data    interface{} `json:"data"`
}

// PaginatedWebResponse represents a paginated API response with data and pagination metadata
type PaginatedWebResponse struct {
	Message    string             `json:"message"`
	Status     int                `json:"status"`
	Data       interface{}        `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
}
