package services

import (
	"context"
	"encoding/json"
	"errors"

	"time"

	appErrors "github.com/chud-lori/go-boilerplate/pkg/errors"

	"github.com/chud-lori/go-boilerplate/domain/entities"
	"github.com/chud-lori/go-boilerplate/domain/ports"
	"github.com/chud-lori/go-boilerplate/pkg/logger"
	"github.com/sirupsen/logrus"
)

type UserServiceImpl struct {
	DB ports.Database
	ports.UserRepository
	ports.Encryptor
	ports.Cache
	CtxTimeout time.Duration
}

func (s *UserServiceImpl) Save(c context.Context, user *entities.User) (*entities.User, error) {
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

	password, err := s.Encryptor.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = password
	result, err := s.UserRepository.Save(ctx, tx, user)

	if err != nil {
		logger.WithError(err).Error("Failed to save user")
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		logger.WithError(err).Error("Failed to commit transaction")
		return nil, err
	}

	return result, nil
}

func (s *UserServiceImpl) Update(c context.Context, user *entities.User) (*entities.User, error) {
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

	result, err := s.UserRepository.Update(ctx, tx, user)

	if err != nil {
		if errors.Is(err, appErrors.ErrUserNotFound) {
			logger.Errorf("UserID %d not found", user.Id)
			return nil, appErrors.NewNotFoundError("User not found", err)
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

func (s *UserServiceImpl) Delete(c context.Context, id string) error {
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

	err = s.UserRepository.Delete(ctx, tx, id)
	if err != nil {
		if errors.Is(err, appErrors.ErrUserNotFound) {
			logger.Errorf("UserID %d not found", id)
			return appErrors.NewNotFoundError("User not found", err)
		}

		logger.WithError(err).Error("Database error")
		return err
	}

	if err := tx.Commit(); err != nil {
		logger.WithError(err).Error("Failed to commit transaction")
		return err
	}

	return nil
}

func (s *UserServiceImpl) FindById(c context.Context, id string) (*entities.User, error) {
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

	result, err := s.UserRepository.FindById(ctx, tx, id)
	if err != nil {
		if errors.Is(err, appErrors.ErrUserNotFound) {
			logger.Errorf("UserID %d not found", id)
			return nil, appErrors.NewNotFoundError("User not found", err)
		}

		logger.WithError(err).Error("Database error")
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		logger.WithError(err).Error("Failed to commit transaction")
		return nil, err
	}

	return result, nil
}

func (s *UserServiceImpl) FindAll(c context.Context) ([]*entities.User, error) {
	logger, _ := c.Value(logger.LoggerContextKey).(logrus.FieldLogger)
	ctx, cancel := context.WithTimeout(c, s.CtxTimeout)
	defer cancel()

	var users []*entities.User

	usersCached, err := s.Cache.Get(c, "users")
	if err = json.Unmarshal([]byte(usersCached), &users); err == nil {
		return users, err
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

	users, err = s.UserRepository.FindAll(ctx, tx)

	if err != nil {
		logger.WithError(err).Error("Failed to find all users")
		return nil, err
	}

	usersString, _ := json.Marshal(&users)
	s.Cache.Set(c, "users", usersString, 30*time.Second)

	if err := tx.Commit(); err != nil {
		logger.WithError(err).Error("Failed to commit transaction")
		return nil, err
	}

	return users, nil
}
