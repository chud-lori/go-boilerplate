package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/chud-lori/go-boilerplate/domain/entities"
	"github.com/chud-lori/go-boilerplate/domain/ports"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type UserRepositoryPostgre struct {
	db     ports.Database
	logger *logrus.Entry
}

func NewUserRepositoryPostgre(db ports.Database) (*UserRepositoryPostgre, error) {
	return &UserRepositoryPostgre{
		db: db,
	}, nil
}

func (repository *UserRepositoryPostgre) Save(ctx context.Context, user *entities.User) (*entities.User, error) {
	tx, err := repository.db.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	var id string
	var createdAt time.Time
	query := `
            INSERT INTO users (email, passcode)
            VALUES ($1, $2)
            RETURNING id, created_at`
	err = tx.QueryRowContext(ctx, query, user.Email, user.Passcode).Scan(&id, &createdAt)
	if err != nil {
		return nil, err
	}

	user.Id = id
	user.Created_at = createdAt

	if err = tx.Commit(); err!= nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return user, nil
}

func (repository *UserRepositoryPostgre) Update(ctx context.Context, user *entities.User) (*entities.User, error) {
    tx, err := repository.db.BeginTx(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to begin transaction: %v", err)
    }
    defer tx.Rollback()

    query := "UPDATE users SET email = $1, passcode = $2 WHERE id = $3"
    _, err = tx.ExecContext(ctx, query, user.Email, user.Passcode, user.Id)
    if err != nil {
        return nil, err
    }

    if err = tx.Commit(); err != nil {
        return nil, fmt.Errorf("failed to commit transaction: %v", err)
    }

    return user, nil
}

func (repository *UserRepositoryPostgre) Delete(ctx context.Context, id string) error {
    tx, err := repository.db.BeginTx(ctx)
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %v", err)
    }
    defer tx.Rollback()

    query := "DELETE FROM users WHERE id = $1"
    result, err := tx.ExecContext(ctx, query, id)
    if err != nil {
        return fmt.Errorf("failed to delete user: %v", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %v", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("user not found")
    }

    if err = tx.Commit(); err != nil {
        return fmt.Errorf("failed to commit transaction: %v", err)
    }

    return nil
}

func (r *UserRepositoryPostgre) FindById(ctx context.Context, id string) (*entities.User, error) {
    tx, err := r.db.BeginTx(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to begin transaction: %v", err)
    }
    defer tx.Rollback()

    if _, err := uuid.Parse(id); err != nil {
		//r.logger.Info("Invalid UUID Format: ", id)
        return nil, fmt.Errorf("Invalid UUID Format")
    }
    
    user := &entities.User{}
    query := "SELECT id, email, created_at FROM users WHERE id = $1"
    err = tx.QueryRowContext(ctx, query, id).Scan(&user.Id, &user.Email, &user.Created_at)

    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("user not found")
        }
        return nil, err
    }

    if err = tx.Commit(); err != nil {
        return nil, fmt.Errorf("failed to commit transaction: %v", err)
    }

    return user, nil
}

func (repository *UserRepositoryPostgre) FindAll(ctx context.Context) ([]*entities.User, error) {
    tx, err := repository.db.BeginTx(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to begin transaction: %v", err)
    }
    defer tx.Rollback()

    query := "SELECT id, email, created_at FROM users"
    rows, err := tx.QueryContext(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []*entities.User
    for rows.Next() {
        var user entities.User
        err := rows.Scan(&user.Id, &user.Email, &user.Created_at)
        if err != nil {
            return nil, err
        }
        users = append(users, &user)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    if err = tx.Commit(); err != nil {
        return nil, fmt.Errorf("failed to commit transaction: %v", err)
    }

    return users, nil
}
