package ports

import (
	"context"

	"github.com/chud-lori/go-boilerplate/domain/entities"
)

type AuthService interface {
	SignIn(ctx context.Context, user *entities.User) (*entities.User, string, error)
	SignUp(ctx context.Context, user *entities.User) (*entities.User, string, error)
}
