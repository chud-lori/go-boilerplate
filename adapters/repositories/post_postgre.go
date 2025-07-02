package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chud-lori/go-boilerplate/domain/entities"
	"github.com/chud-lori/go-boilerplate/domain/ports"
	appErrors "github.com/chud-lori/go-boilerplate/pkg/errors"
	"github.com/chud-lori/go-boilerplate/pkg/logger"
	"github.com/google/uuid"

	"github.com/sirupsen/logrus"
)

type PostRepositoryPostgre struct {
}

func (r *PostRepositoryPostgre) Save(ctx context.Context, tx ports.Transaction, post *entities.Post) (*entities.Post, error) {
	logger, _ := ctx.Value(logger.LoggerContextKey).(logrus.FieldLogger)

	var id uuid.UUID
	var title string

	query := `
            INSERT INTO posts (title, body, author_id)
            VALUES ($1, $2, $3)
            RETURNING id, title`
	err := tx.QueryRowContext(ctx, query, post.Title, post.Body, post.AuthorID).Scan(&id, &title)
	if err != nil {
		logger.Error("Failed to post: ", err)
		return nil, err
	}

	post.ID = id
	post.Title = title

	return post, nil
}

func (r *PostRepositoryPostgre) Update(ctx context.Context, tx ports.Transaction, post *entities.Post) (*entities.Post, error) {
	logger, _ := ctx.Value(logger.LoggerContextKey).(logrus.FieldLogger)

	query := "UPDATE posts SET title = $1, body = $2 WHERE id = $3"
	result, err := tx.ExecContext(ctx, query, post.Title, post.Body, post.ID)

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
		logger.Errorf("Post ID %s not found", post.ID)
		return nil, appErrors.ErrDataNotFound
	}

	return post, nil
}

func (r *PostRepositoryPostgre) Delete(ctx context.Context, tx ports.Transaction, id uuid.UUID) error {
	logger, _ := ctx.Value(logger.LoggerContextKey).(logrus.FieldLogger)

	query := "DELETE FROM posts WHERE id = $1"
	result, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete post: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.WithError(err).Error("Failed get row affected")
		return err
	}

	if rowsAffected == 0 {
		return appErrors.ErrDataNotFound
	}

	return nil
}

func (r *PostRepositoryPostgre) GetById(ctx context.Context, tx ports.Transaction, id uuid.UUID) (*entities.Post, error) {
	post := &entities.Post{}
	query := "SELECT id, title, body, author_id, created_at FROM posts WHERE id = $1"
	err := tx.QueryRowContext(ctx, query, id).Scan(&post.ID, &post.Title, &post.Body, &post.AuthorID, &post.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, appErrors.ErrDataNotFound
		}
		return nil, err
	}

	return post, nil
}

func (r *PostRepositoryPostgre) GetAll(ctx context.Context, tx ports.Transaction, search string, pagination entities.PaginationParams) ([]entities.Post, error) {
	query := "SELECT id, title, body, author_id, created_at FROM posts WHERE 1=1"
	args := []interface{}{}
	argCounter := 1

	if search != "" {
		query += fmt.Sprintf(" AND title ILIKE $%d", argCounter)
		args = append(args, "%"+search+"%")
		argCounter++
	}
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argCounter, argCounter+1)
	args = append(args, pagination.Limit, (pagination.Page-1)*pagination.Limit)

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	var posts []entities.Post
	for rows.Next() {
		var post entities.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Body, &post.AuthorID, &post.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("Failed to scan post row")
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("errors during rows iteration")
	}

	return posts, nil
}
