package dto

import "github.com/google/uuid"

// CreatePostRequest represents the request body for creating a new post.
type CreatePostRequest struct {
	Title    string    `json:"title"`
	Body     string    `json:"body"`
	AuthorID uuid.UUID `json:"author_id"` // Assuming author_id comes from the client
}

// UpdatePostRequest represents the request body for updating an existing post.
type UpdatePostRequest struct {
	Title *string `json:"title"` // Pointers allow distinguishing missing field from empty string
	Body  *string `json:"body"`
}
