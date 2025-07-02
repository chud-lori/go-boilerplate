package dto

import (
	"time"

	"github.com/google/uuid"
)

// PostResponse represents the response structure for a single post.
type PostResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	AuthorID  uuid.UUID `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
}
