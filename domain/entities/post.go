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

// PostAttachment represents a file/image/video attached to a post.
// It is used for async upload processing and status tracking.
type PostAttachment struct {
	ID         uuid.UUID `json:"id"`
	PostID     uuid.UUID `json:"post_id"`
	FileName   string    `json:"file_name"`
	FileType   string    `json:"file_type"`
	FileURL    string    `json:"file_url"`
	Status     UploadStatus `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// UploadStatus represents the state of an async upload process for an attachment.
type UploadStatus string

const (
	UploadStatusPending   UploadStatus = "pending"
	UploadStatusUploading UploadStatus = "uploading"
	UploadStatusSuccess   UploadStatus = "success"
	UploadStatusFailed    UploadStatus = "failed"
)
