package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"time"

	appErrors "github.com/chud-lori/go-boilerplate/pkg/errors"
	"github.com/google/uuid"

	"github.com/chud-lori/go-boilerplate/domain/entities"
	"github.com/chud-lori/go-boilerplate/domain/ports"
	"github.com/chud-lori/go-boilerplate/pkg/logger"
	"github.com/sirupsen/logrus"
)

type PostServiceImpl struct {
	DB ports.Database
	ports.PostRepository
	ports.UserRepository
	ports.Cache
	CtxTimeout time.Duration
}

func (s *PostServiceImpl) Create(c context.Context, post *entities.Post) (*entities.Post, error) {
	logger, _ := c.Value(logger.LoggerContextKey).(logrus.FieldLogger)
	ctx, cancel := context.WithTimeout(c, s.CtxTimeout)
	defer cancel()

	tx, err := s.DB.BeginTx(ctx)
	if err != nil {
		logger.WithError(err).Error("Failed to begin transaction")
		return nil, err
	}

	// handle panic gracefully
	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		}
	}()

	_, err = s.UserRepository.FindById(ctx, tx, post.User.ID.String())
	if errors.Is(err, appErrors.ErrUserNotFound) {
		return nil, appErrors.NewNotFoundError("Author not found", err)
	}

	result, err := s.PostRepository.Save(ctx, tx, post)

	if err != nil {
		logger.WithError(err).Error("Failed to save post")
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		logger.WithError(err).Error("Failed to commit transaction")
		return nil, err
	}

	return result, nil
}

func (s *PostServiceImpl) Update(c context.Context, post *entities.Post) (*entities.Post, error) {
	logger, _ := c.Value(logger.LoggerContextKey).(logrus.FieldLogger)

	ctx, cancel := context.WithTimeout(c, s.CtxTimeout)
	defer cancel()

	tx, err := s.DB.BeginTx(ctx)
	if err != nil {
		logger.WithError(err).Error("Failed to begin transaction")
		return nil, err
	}

	defer func() {
		if r := recover(); r != nil || err != nil {
			logger.Errorf("Transaction rollback due to error: %v", err)
			logger.Errorf("Transaction rollback due to panic: %v", r)
			tx.Rollback()
		}
	}()

	result, err := s.PostRepository.Update(ctx, tx, post)

	if err != nil {
		if errors.Is(err, appErrors.ErrDataNotFound) {
			logger.Errorf("PostID %d not found", post.ID)
			return nil, appErrors.NewNotFoundError("Post not found", err)
		}

		logger.WithError(err).Error("Database error")
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		logger.WithError(err).Error("Failed to commit transaction")
		return nil, err
	}
	err = s.Cache.InvalidateByPrefix(c, "posts:") // Assuming "posts:" is the prefix for all your post-related keys
	if err != nil {
		logger.WithError(err).Warn("Failed to invalidate 'posts:' cache keys. Stale data might be served.")
	} else {
		logger.Debug("Successfully invalidated 'posts:' cache keys.")
	}

	return result, nil
}

func (s *PostServiceImpl) Delete(c context.Context, id uuid.UUID) error {
	logger, _ := c.Value(logger.LoggerContextKey).(logrus.FieldLogger)

	ctx, cancel := context.WithTimeout(c, s.CtxTimeout)
	defer cancel()

	tx, err := s.DB.BeginTx(ctx)
	if err != nil {
		logger.WithError(err).Error("Failed to begin transaction")
		return err
	}

	defer func() {
		if r := recover(); r != nil || err != nil {
			logger.Errorf("Transaction rollback due to error: %v", err)
			logger.Errorf("Transaction rollback due to panic: %v", r)
			tx.Rollback()
		}
	}()

	err = s.PostRepository.Delete(ctx, tx, id)
	if err != nil {
		if errors.Is(err, appErrors.ErrDataNotFound) {
			logger.Errorf("PostID %d not found", id)
			return appErrors.NewNotFoundError("Post not found", err)
		}

		logger.WithError(err).Error("Database error")
		return err
	}

	if err = tx.Commit(); err != nil {
		logger.WithError(err).Error("Failed to commit transaction")
		return err
	}

	return nil
}

func (s *PostServiceImpl) GetById(c context.Context, id uuid.UUID) (*entities.Post, error) {
	logger, _ := c.Value(logger.LoggerContextKey).(logrus.FieldLogger)

	ctx, cancel := context.WithTimeout(c, s.CtxTimeout)
	defer cancel()

	tx, err := s.DB.BeginTx(ctx)
	if err != nil {
		logger.WithError(err).Error("Failed to begin transaction")
		return nil, err
	}

	defer func() {
		if r := recover(); r != nil || err != nil {
			logger.Errorf("Transaction rollback due to error: %v", err)
			logger.Errorf("Transaction rollback due to panic: %v", r)
			tx.Rollback()
		}
	}()

	result, err := s.PostRepository.GetById(ctx, tx, id)
	if err != nil {
		if errors.Is(err, appErrors.ErrDataNotFound) {
			logger.Errorf("PostId %d not found", id)
			return nil, appErrors.NewNotFoundError("Post not found", err)
		}

		logger.WithError(err).Error("Database error")
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		logger.WithError(err).Error("Failed to commit transaction")
		return nil, err
	}

	return result, nil
}

func (s *PostServiceImpl) GetAll(c context.Context, search string, page, limit int) ([]entities.Post, error) {
	logger, _ := c.Value(logger.LoggerContextKey).(logrus.FieldLogger)
	ctx, cancel := context.WithTimeout(c, s.CtxTimeout)
	defer cancel()

	var posts []entities.Post

	// cache key based on payload
	queryParams := fmt.Sprintf("search=%s:page=%d:limit=%d", search, page, limit)
	hasher := sha256.New()
	hasher.Write([]byte(queryParams))
	cacheKey := "posts:" + hex.EncodeToString(hasher.Sum(nil))

	// if cached err, won't interupt and using db instead
	postsCached, errCache := s.Cache.Get(c, cacheKey)
	if errCache = json.Unmarshal([]byte(postsCached), &posts); errCache == nil {
		return posts, errCache
	}

	tx, err := s.DB.BeginTx(ctx)
	if err != nil {
		logger.WithError(err).Error("Failed to begin transaction")
		return nil, err
	}

	defer func() {
		if r := recover(); r != nil || err != nil {
			logger.Errorf("Transaction rollback due to error: %v", err)
			logger.Errorf("Transaction rollback due to panic: %v", r)
			tx.Rollback()
		}
	}()

	pagination := entities.PaginationParams{
		Page:  page,
		Limit: limit,
	}

	posts, err = s.PostRepository.GetAll(ctx, tx, search, pagination)

	if err != nil {
		logger.WithError(err).Error("Failed to get all posts")
		return nil, err
	}

	// if cached err, won't interupt
	postsString, _ := json.Marshal(&posts)
	errCache = s.Cache.Set(c, cacheKey, postsString, 30*time.Second)
	if errCache != nil {
		logger.WithError(errCache).Warn("Failed set cache")
	}

	if err = tx.Commit(); err != nil {
		logger.WithError(err).Error("Failed to commit transaction")
		return nil, err
	}

	return posts, nil
}
