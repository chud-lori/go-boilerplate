package ports

import (
	"context"

	"github.com/chud-lori/go-boilerplate/domain/entities"
)

type AuthService interface {
	SignIn(ctx context.Context, user *entities.User) (*entities.User, string, string, error)
	SignUp(ctx context.Context, user *entities.User) (*entities.User, string, string, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
	Logout(ctx context.Context, refreshToken string) error
}
