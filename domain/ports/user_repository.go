package ports

import (
	"context"

	"github.com/chud-lori/go-boilerplate/domain/entities"
)

type UserRepository interface {
	Save(ctx context.Context, tx Transaction, user *entities.User) (*entities.User, error)
	Update(ctx context.Context, tx Transaction, user *entities.User) (*entities.User, error)
	Delete(ctx context.Context, tx Transaction, id string) error
	FindById(ctx context.Context, tx Transaction, id string) (*entities.User, error)
	FindAll(ctx context.Context, tx Transaction) ([]*entities.User, error)
}
