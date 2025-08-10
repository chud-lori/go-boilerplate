package ports

import (
	"context"

	"github.com/chud-lori/go-boilerplate/domain/entities"
	"github.com/google/uuid"
)

type PostService interface {
	Create(ctx context.Context, post *entities.Post) (*entities.Post, error)
	Update(ctx context.Context, post *entities.Post) (*entities.Post, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetById(ctx context.Context, id uuid.UUID) (*entities.Post, error)
	GetAll(ctx context.Context, search string, page, limit int) ([]entities.Post, error)
	StartAsyncUpload(ctx context.Context, postID uuid.UUID, fileName, fileType string, fileData []byte) (uploadID uuid.UUID, err error)
	GetUploadStatus(ctx context.Context, uploadID uuid.UUID) (entities.UploadStatus, error)
}
