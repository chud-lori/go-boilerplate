package ports

import (
	"context"

	"github.com/chud-lori/go-boilerplate/domain/entities"
	"github.com/google/uuid"
)

type PostRepository interface {
	Save(ctx context.Context, tx Transaction, post *entities.Post) (*entities.Post, error)
	Update(ctx context.Context, tx Transaction, post *entities.Post) (*entities.Post, error)
	Delete(ctx context.Context, tx Transaction, id uuid.UUID) error
	CountPost(ctx context.Context, tx Transaction) (uint32, error)
	GetById(ctx context.Context, tx Transaction, id uuid.UUID) (*entities.Post, error)
	GetAll(ctx context.Context, tx Transaction, search string, pagination entities.PaginationParams) ([]entities.Post, error)
}
