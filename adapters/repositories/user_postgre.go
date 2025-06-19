package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chud-lori/go-boilerplate/domain/entities"
	"github.com/chud-lori/go-boilerplate/domain/ports"
	appErrors "github.com/chud-lori/go-boilerplate/pkg/errors"
	"github.com/chud-lori/go-boilerplate/pkg/logger"

	"github.com/sirupsen/logrus"
)

type UserRepositoryPostgre struct {
	DB ports.Database
}

func (repository *UserRepositoryPostgre) Save(ctx context.Context, tx ports.Transaction, user *entities.User) (*entities.User, error) {
	logger, _ := ctx.Value(logger.LoggerContextKey).(logrus.FieldLogger)

	var id string
	var createdAt time.Time
	query := `
            INSERT INTO users (email, password)
            VALUES ($1, $2)
            RETURNING id, created_at`
	err := tx.QueryRowContext(ctx, query, user.Email, user.Password).Scan(&id, &createdAt)
	if err != nil {
		logger.Error("Failed to insert user: ", err)
		return nil, err
	}

	user.Id = id
	user.CreatedAt = createdAt

	return user, nil
}

func (repository *UserRepositoryPostgre) Update(ctx context.Context, tx ports.Transaction, user *entities.User) (*entities.User, error) {
	logger, _ := ctx.Value(logger.LoggerContextKey).(logrus.FieldLogger)

	query := "UPDATE users SET email = $1, password = $2 WHERE id = $3"
	result, err := tx.ExecContext(ctx, query, user.Email, user.Password, user.Id)

	if err != nil {
		logger.WithError(err).Error("Error Update")
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.WithError(err).Error("Failed get row affected")
		return nil, err
	}

	if rowsAffected == 0 {
		logger.Error("User ID %s not found", user.Id)
		return nil, appErrors.ErrUserNotFound
	}

	return user, nil
}

func (repository *UserRepositoryPostgre) Delete(ctx context.Context, tx ports.Transaction, id string) error {
	logger, _ := ctx.Value(logger.LoggerContextKey).(logrus.FieldLogger)

	query := "DELETE FROM users WHERE id = $1"
	result, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.WithError(err).Error("Failed get row affected")
		return err
	}

	if rowsAffected == 0 {
		return appErrors.ErrUserNotFound
	}

	return nil
}

func (r *UserRepositoryPostgre) FindById(ctx context.Context, tx ports.Transaction, id string) (*entities.User, error) {
	user := &entities.User{}
	query := "SELECT id, email, created_at FROM users WHERE id = $1"
	err := tx.QueryRowContext(ctx, query, id).Scan(&user.Id, &user.Email, &user.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, appErrors.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *UserRepositoryPostgre) FindByEmail(ctx context.Context, tx ports.Transaction, email string) (*entities.User, error) {
	user := &entities.User{}
	query := "SELECT id, email, created_at FROM users WHERE email = $1"
	err := tx.QueryRowContext(ctx, query, email).Scan(&user.Id, &user.Email, &user.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, appErrors.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (repository *UserRepositoryPostgre) FindAll(ctx context.Context, tx ports.Transaction) ([]*entities.User, error) {
	query := "SELECT id, email, created_at FROM users"
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entities.User
	for rows.Next() {
		var user entities.User
		err := rows.Scan(&user.Id, &user.Email, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
