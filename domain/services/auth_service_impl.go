package services

import (
	"context"

	"time"

	"encoding/json"

	"github.com/chud-lori/go-boilerplate/domain/entities"
	"github.com/chud-lori/go-boilerplate/domain/ports"
	appErrors "github.com/chud-lori/go-boilerplate/pkg/errors"
	"github.com/chud-lori/go-boilerplate/pkg/logger"
	"github.com/sirupsen/logrus"
)

type AuthServiceImpl struct {
	DB ports.Database
	ports.UserRepository
	ports.Encryptor
	ports.TokenManager
	ports.MailService
	ExternalApi ports.ExternalApiClient
	CtxTimeout  time.Duration
}

func (s *AuthServiceImpl) SignIn(c context.Context, user *entities.User) (*entities.User, string, error) {
	logger, _ := c.Value(logger.LoggerContextKey).(logrus.FieldLogger)
	ctx, cancel := context.WithTimeout(c, s.CtxTimeout)
	defer cancel()

	tx, err := s.DB.BeginTx(ctx)
	if err != nil {
		logger.WithError(err).Error("Failed to begin transaction")
		return nil, "", err
	}

	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		}
	}()

	foundUser, err := s.UserRepository.FindByEmail(ctx, tx, user.Email)
	if err != nil {
		logger.WithError(err).Warn("User not found by email")
		return nil, "", appErrors.NewUnauthorizedError("Unauthorized", err)
	}

	if err := s.Encryptor.CompareHash(foundUser.Password, user.Password); err != nil {
		logger.WithError(err).Warn("Invalid password")
		return nil, "", appErrors.NewUnauthorizedError("Unauthorized", err)
	}

	token, err := s.TokenManager.GenerateToken(foundUser.ID.String())
	if err != nil {
		logger.WithError(err).Error("Failed to generate token")
		return nil, "", err
	}

	if err = tx.Commit(); err != nil {
		logger.WithError(err).Error("Failed to commit transaction")
		return nil, "", err
	}

	// Disable throwing error, as this dummy
	// External API call POST dummy
	auditPayload := map[string]interface{}{
		"user_id": foundUser.ID.String(),
		"event":   "sign_in",
		"time":    time.Now().Format(time.RFC3339),
	}
	body, _ := json.Marshal(auditPayload)
	postResp, errApi := s.ExternalApi.DoRequest(ctx, "POST", "https://dummy-external-api.local/audit", map[string]string{"Content-Type": "application/json"}, body)
	if errApi != nil {
		logger.WithError(errApi).Warn("Failed to notify external audit API (POST)")
	} else {
		logger.Infof("POST external API response: %s", string(postResp))
	}

	// External API call GET dummy
	getResp, errApi := s.ExternalApi.DoRequest(ctx, "GET", "https://dummy-external-api.local/userinfo?id="+foundUser.ID.String(), nil, nil)
	if err != nil {
		logger.WithError(errApi).Warn("Failed to get user info from external API (GET)")
	} else {
		logger.Infof("GET external API response: %s", string(getResp))
	}

	// Disable throwing error, as this mail service is optional
	errMail := s.MailService.SendSignInNotification(ctx, foundUser.Email, "User logged in just now")
	if errMail != nil {
		logger.WithError(errMail).Error("Failed to send mail")
	}

	return foundUser, token, err
}

func (s *AuthServiceImpl) SignUp(c context.Context, user *entities.User) (*entities.User, string, error) {
	logger, _ := c.Value(logger.LoggerContextKey).(logrus.FieldLogger)
	ctx, cancel := context.WithTimeout(c, s.CtxTimeout)
	defer cancel()

	tx, err := s.DB.BeginTx(ctx)
	if err != nil {
		logger.WithError(err).Error("Failed to begin transaction")
		return nil, "", err
	}

	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		}
	}()

	existing, _ := s.UserRepository.FindByEmail(ctx, tx, user.Email)
	if existing != nil {
		logger.Warn("Email already in use")
		err = appErrors.NewBadRequestError("Email already exist", nil)
		return nil, "", err
	}

	hashedPassword, err := s.Encryptor.HashPassword(user.Password)
	if err != nil {
		logger.WithError(err).Error("Failed to hash password")
		return nil, "", err
	}
	user.Password = hashedPassword

	if _, err := s.UserRepository.Save(ctx, tx, user); err != nil {
		logger.WithError(err).Error("Failed to create user")
		return nil, "", err
	}

	token, err := s.TokenManager.GenerateToken(user.ID.String())
	if err != nil {
		logger.WithError(err).Error("Failed to generate token")
		return nil, "", err
	}

	if err = tx.Commit(); err != nil {
		logger.WithError(err).Error("Failed to commit transaction")
		return nil, "", err
	}

	return user, token, nil
}
