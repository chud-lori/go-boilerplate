package services

import (
	"context"
	"database/sql"
	"errors"

	"time"

	"github.com/chud-lori/go-boilerplate/pkg/auth"
	appErrors "github.com/chud-lori/go-boilerplate/pkg/errors"

	"github.com/chud-lori/go-boilerplate/domain/entities"
	"github.com/chud-lori/go-boilerplate/domain/ports"
	"github.com/chud-lori/go-boilerplate/pkg/logger"
	"github.com/sirupsen/logrus"
)

type UserServiceImpl struct {
	DB ports.Database
	ports.UserRepository
	CtxTimeout time.Duration
}

func (s *UserServiceImpl) Save(c context.Context, request *entities.User) (*entities.User, error) {
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

	request.Passcode = auth.GeneratePasscode()
	result, err := s.UserRepository.Save(ctx, tx, request)

	if err != nil {
		logger.WithError(err).Error("Failed to save user")
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		logger.WithError(err).Error("Failed to commit transaction")
		return nil, err
	}

	return result, nil
}

func (s *UserServiceImpl) Update(c context.Context, request *entities.User) (*entities.User, error) {
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

	result, err := s.UserRepository.Update(ctx, tx, request)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Errorf("UserID %d not found", request.Id)
			return nil, appErrors.NewBadRequestError("User not found", err)
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

	_, err = s.UserRepository.FindById(ctx, tx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Errorf("UserID %d not found", id)
			return appErrors.NewBadRequestError("User not found", err)
		}

		logger.WithError(err).Error("Failed to find userId: ", id)
		return err
	}

	err = s.UserRepository.Delete(ctx, tx, id)

	if err != nil {
		logger.WithError(err).Error("Failed to delete user")
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
		if errors.Is(err, sql.ErrNoRows) {
			logger.Errorf("UserID %d not found", id)
			return nil, appErrors.NewBadRequestError("User not found", err)
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

	result, err := s.UserRepository.FindAll(ctx, tx)

	if err != nil {
		logger.WithError(err).Error("Failed to find all users")
		return nil, err
	}

	users := make([]*entities.User, len(result))

	for i, user := range result {
		users[i] = &entities.User{
			Id:         user.Id,
			Email:      user.Email,
			Created_at: user.Created_at,
		}
	}

	if err := tx.Commit(); err != nil {
		logger.WithError(err).Error("Failed to commit transaction")
		return nil, err
	}

	return users, nil
}
