package repositories

import (
	"boilerplate/domain/entities"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type UserRepositoryPostgre struct {
	db     *sql.DB
	logger *logrus.Entry
}

func NewUserRepositoryPostgre(db *sql.DB) (*UserRepositoryPostgre, error) {
	return &UserRepositoryPostgre{
		db: db,
	}, nil
}

func (repository *UserRepositoryPostgre) Save(ctx context.Context, user *entities.User) (*entities.User, error) {
	var id string
	var createdAt time.Time
	query := `
            INSERT INTO users (email, passcode)
            VALUES ($1, $2)
            RETURNING id, created_at`
	err := repository.db.QueryRow(query, user.Email, user.Passcode).Scan(&id, &createdAt)
	if err != nil {
		return nil, err
	}

	user.Id = id
	user.Created_at = createdAt

	return user, nil
}

func (repository *UserRepositoryPostgre) Update(ctx context.Context, user *entities.User) (*entities.User, error) {
	query := "UPDATE users SET email = $1, passcode = $2 WHERE id = $3"
	_, err := repository.db.Exec(query, user.Email, user.Passcode, user.Id)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repository *UserRepositoryPostgre) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM users WHERE id = $1"

	_, err := repository.db.Exec(query, id)

	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepositoryPostgre) FindById(ctx context.Context, id string) (*entities.User, error) {
	if _, err := uuid.Parse(id); err != nil {
		//r.logger.Info("Invalid UUID Format: ", id)
		return nil, fmt.Errorf("Invalid UUID Format")
	}
	user := &entities.User{}

	query := "SELECT id, email, created_at FROM users WHERE id = $1"

	err := r.db.QueryRow(query, id).Scan(&user.Id, &user.Email, &user.Created_at)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No data found")
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (repository *UserRepositoryPostgre) FindAll(ctx context.Context) ([]*entities.User, error) {
	query := "SELECT id, email, created_at FROM users"

	rows, err := repository.db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

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

	return users, nil
}
