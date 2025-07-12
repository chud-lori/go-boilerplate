package entities

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	User      *User     `json:"author,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PaginationParams struct {
	Page  int
	Limit int
}

// PaginatedPostsCache is used for caching paginated posts and total count
// for the GetAllPaginated service method.
type PaginatedPostsCache struct {
	Posts      []Post `json:"posts"`
	TotalItems int    `json:"total_items"`
}
