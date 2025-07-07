package repositories_test

import (
	"context"
	"testing"

	"github.com/chud-lori/go-boilerplate/adapters/repositories"
	"github.com/chud-lori/go-boilerplate/domain/entities"
	"github.com/chud-lori/go-boilerplate/domain/ports"
	"github.com/chud-lori/go-boilerplate/internal/testutils"
	appErrors "github.com/chud-lori/go-boilerplate/pkg/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_Save(t *testing.T) {
	testutils.WithTransactionTest(t,
		func(db ports.Database) (ports.UserRepository, error) {
			return &repositories.UserRepositoryPostgre{}, nil
		},
		func(ctx context.Context, repo ports.UserRepository, tx ports.Transaction) {
			user := &entities.User{
				Email:    "save@example.com",
				Password: "pass123",
			}
			savedUser, err := repo.Save(ctx, tx, user)
			require.NoError(t, err)
			require.NotEmpty(t, savedUser.ID)
		},
	)
}

func TestUserRepository_FindById(t *testing.T) {
	testutils.WithTransactionTest(t,
		func(db ports.Database) (ports.UserRepository, error) {
			return &repositories.UserRepositoryPostgre{}, nil
		},
		func(ctx context.Context, repo ports.UserRepository, tx ports.Transaction) {
			user := &entities.User{
				Email:    "find@example.com",
				Password: "pass123",
			}
			savedUser, _ := repo.Save(ctx, tx, user)
			found, err := repo.FindById(ctx, tx, savedUser.ID.String())
			require.NoError(t, err)
			require.Equal(t, savedUser.ID, found.ID)
			require.Equal(t, savedUser.Email, found.Email)
		},
	)
}

func TestUserRepository_FindById_NotFound(t *testing.T) {
	testutils.WithTransactionTest(t,
		func(db ports.Database) (ports.UserRepository, error) {
			return &repositories.UserRepositoryPostgre{}, nil
		},
		func(ctx context.Context, repo ports.UserRepository, tx ports.Transaction) {
			_, err := repo.FindById(ctx, tx, "ad24a17d-2925-4aa8-b077-d358a0788df7")
			require.ErrorIs(t, err, appErrors.ErrUserNotFound)
		},
	)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	testutils.WithTransactionTest(t,
		func(db ports.Database) (ports.UserRepository, error) {
			return &repositories.UserRepositoryPostgre{}, nil
		},
		func(ctx context.Context, repo ports.UserRepository, tx ports.Transaction) {
			user := &entities.User{
				Email:    "find@example.com",
				Password: "pass123",
			}
			savedUser, _ := repo.Save(ctx, tx, user)
			found, err := repo.FindByEmail(ctx, tx, savedUser.Email)
			require.NoError(t, err)
			require.Equal(t, savedUser.ID, found.ID)
			require.Equal(t, savedUser.Email, found.Email)
		},
	)
}

func TestUserRepository_FindByEmail_NotFound(t *testing.T) {
	testutils.WithTransactionTest(t,
		func(db ports.Database) (ports.UserRepository, error) {
			return &repositories.UserRepositoryPostgre{}, nil
		},
		func(ctx context.Context, repo ports.UserRepository, tx ports.Transaction) {
			_, err := repo.FindByEmail(ctx, tx, "find@example.com")
			require.ErrorIs(t, err, appErrors.ErrUserNotFound)
		},
	)
}

func TestUserRepository_Update(t *testing.T) {
	testutils.WithTransactionTest(t,
		func(db ports.Database) (ports.UserRepository, error) {
			return &repositories.UserRepositoryPostgre{}, nil
		},
		func(ctx context.Context, repo ports.UserRepository, tx ports.Transaction) {
			user := &entities.User{
				Email:    "update@example.com",
				Password: "pass123",
			}
			saved, _ := repo.Save(ctx, tx, user)
			saved.Email = "updated@example.com"
			updated, err := repo.Update(ctx, tx, saved)
			require.NoError(t, err)
			require.Equal(t, "updated@example.com", updated.Email)
		},
	)
}

func TestUserRepository_Update_NotFound(t *testing.T) {
	testutils.WithTransactionTest(t,
		func(db ports.Database) (ports.UserRepository, error) {
			return &repositories.UserRepositoryPostgre{}, nil
		},
		func(ctx context.Context, repo ports.UserRepository, tx ports.Transaction) {
			nonExistent := &entities.User{ID: uuid.New(), Email: "none@example.com", Password: "123"}
			_, err := repo.Update(ctx, tx, nonExistent)
			require.ErrorIs(t, err, appErrors.ErrUserNotFound)
		},
	)
}

func TestUserRepository_Delete(t *testing.T) {
	testutils.WithTransactionTest(t,
		func(db ports.Database) (ports.UserRepository, error) {
			return &repositories.UserRepositoryPostgre{}, nil
		},
		func(ctx context.Context, repo ports.UserRepository, tx ports.Transaction) {
			user := &entities.User{
				Email:    "delete@example.com",
				Password: "pass123",
			}
			saved, _ := repo.Save(ctx, tx, user)
			err := repo.Delete(ctx, tx, saved.ID.String())
			require.NoError(t, err)
		},
	)
}

func TestUserRepository_Delete_NotFound(t *testing.T) {
	testutils.WithTransactionTest(t,
		func(db ports.Database) (ports.UserRepository, error) {
			return &repositories.UserRepositoryPostgre{}, nil
		},
		func(ctx context.Context, repo ports.UserRepository, tx ports.Transaction) {
			err := repo.Delete(ctx, tx, "ad24a17d-2925-4aa8-b077-d358a0788df7")
			require.ErrorIs(t, err, appErrors.ErrUserNotFound)
		},
	)
}

func TestUserRepository_FindAll(t *testing.T) {
	testutils.WithTransactionTest(t,
		func(db ports.Database) (ports.UserRepository, error) {
			return &repositories.UserRepositoryPostgre{}, nil
		},
		func(ctx context.Context, repo ports.UserRepository, tx ports.Transaction) {
			repo.Save(ctx, tx, &entities.User{Email: "all1@example.com", Password: "123"})
			repo.Save(ctx, tx, &entities.User{Email: "all2@example.com", Password: "123"})
			users, err := repo.FindAll(ctx, tx)
			require.NoError(t, err)
			require.GreaterOrEqual(t, len(users), 2)
		},
	)
}
