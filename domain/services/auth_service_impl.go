package services

import (
	"context"

	"time"

	"github.com/chud-lori/go-boilerplate/domain/entities"
	"github.com/chud-lori/go-boilerplate/domain/ports"
	appErrors "github.com/chud-lori/go-boilerplate/pkg/errors"
	"github.com/chud-lori/go-boilerplate/pkg/logger"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type AuthServiceImpl struct {
	DB ports.Database
	ports.UserRepository
	ports.Encryptor
	ports.TokenManager
	ports.MailService
    ports.Cache // Add cache for refresh tokens
	CtxTimeout time.Duration
}

func (s *AuthServiceImpl) SignIn(c context.Context, user *entities.User) (*entities.User, string, string, error) {
	logger, _ := c.Value(logger.LoggerContextKey).(logrus.FieldLogger)
	ctx, cancel := context.WithTimeout(c, s.CtxTimeout)
	defer cancel()

	tx, err := s.DB.BeginTx(ctx)
	if err != nil {
		logger.WithError(err).Error("Failed to begin transaction")
		return nil, "", "", err
	}

	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		}
	}()

	foundUser, err := s.UserRepository.FindByEmail(ctx, tx, user.Email)
	if err != nil {
		logger.WithError(err).Warn("User not found by email")
		return nil, "", "", appErrors.NewUnauthorizedError("Unauthorized", err)
	}

	if err := s.Encryptor.CompareHash(foundUser.Password, user.Password); err != nil {
		logger.WithError(err).Warn("Invalid password")
		return nil, "", "", appErrors.NewUnauthorizedError("Unauthorized", err)
	}

	token, err := s.TokenManager.GenerateToken(foundUser.ID.String())
	if err != nil {
		logger.WithError(err).Error("Failed to generate token")
		return nil, "", "", err
	}

	// Generate refresh token (opaque random string)
	refreshToken := uuid.NewString()
	refreshKey := "refresh_token:" + refreshToken
	// Store in Redis: key=refresh_token:<token>, value=userID, expiry=7 days
	err = s.Cache.Set(ctx, refreshKey, []byte(foundUser.ID.String()), 7*24*time.Hour)
	if err != nil {
		logger.WithError(err).Error("Failed to store refresh token")
		return nil, "", "", err
	}

	if err = tx.Commit(); err != nil {
		logger.WithError(err).Error("Failed to commit transaction")
		return nil, "", "", err
	}

	// Disable throwing error, as this mail service is optional
	errMail := s.MailService.SendSignInNotification(ctx, foundUser.Email, "User logged in just now")
	if errMail != nil {
		logger.WithError(err).Error("Failed to send mail")
		// return nil, "", err
	}

	return foundUser, token, refreshToken, err
}

func (s *AuthServiceImpl) SignUp(c context.Context, user *entities.User) (*entities.User, string, string, error) {
	logger, _ := c.Value(logger.LoggerContextKey).(logrus.FieldLogger)
	ctx, cancel := context.WithTimeout(c, s.CtxTimeout)
	defer cancel()

	tx, err := s.DB.BeginTx(ctx)
	if err != nil {
		logger.WithError(err).Error("Failed to begin transaction")
		return nil, "", "", err
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
		return nil, "", "", err
	}

	hashedPassword, err := s.Encryptor.HashPassword(user.Password)
	if err != nil {
		logger.WithError(err).Error("Failed to hash password")
		return nil, "", "", err
	}
	user.Password = hashedPassword

	if _, err := s.UserRepository.Save(ctx, tx, user); err != nil {
		logger.WithError(err).Error("Failed to create user")
		return nil, "", "", err
	}

	token, err := s.TokenManager.GenerateToken(user.ID.String())
	if err != nil {
		logger.WithError(err).Error("Failed to generate token")
		return nil, "", "", err
	}

	// Generate refresh token (opaque random string)
	refreshToken := uuid.NewString()
	refreshKey := "refresh_token:" + refreshToken
	// Store in Redis: key=refresh_token:<token>, value=userID, expiry=7 days
	err = s.Cache.Set(ctx, refreshKey, []byte(user.ID.String()), 7*24*time.Hour)
	if err != nil {
		logger.WithError(err).Error("Failed to store refresh token")
		return nil, "", "", err
	}

	if err = tx.Commit(); err != nil {
		logger.WithError(err).Error("Failed to commit transaction")
		return nil, "", "", err
	}

	return user, token, refreshToken, nil
}

// RefreshToken issues a new access and refresh token if the provided refresh token is valid
func (s *AuthServiceImpl) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
   logger, _ := ctx.Value(logger.LoggerContextKey).(logrus.FieldLogger)
   cacheKey := "refresh_token:" + refreshToken
   userID, err := s.Cache.Get(ctx, cacheKey)
   if err != nil || userID == "" {
       logger.WithError(err).Warn("Invalid or expired refresh token")
       return "", "", appErrors.NewUnauthorizedError("Invalid or expired refresh token", err)
   }

   // Rotate refresh token: delete old, create new
   _ = s.Cache.Delete(ctx, cacheKey)
   newRefreshToken := uuid.NewString()
   newCacheKey := "refresh_token:" + newRefreshToken
   err = s.Cache.Set(ctx, newCacheKey, []byte(userID), 7*24*time.Hour)
   if err != nil {
       logger.WithError(err).Error("Failed to store new refresh token")
       return "", "", err
   }

   // Issue new access token
   token, err := s.TokenManager.GenerateToken(userID)
   if err != nil {
       logger.WithError(err).Error("Failed to generate new access token")
       return "", "", err
   }
   return token, newRefreshToken, nil
}

// Logout deletes the refresh token from Redis
func (s *AuthServiceImpl) Logout(ctx context.Context, refreshToken string) error {
    logger, _ := ctx.Value(logger.LoggerContextKey).(logrus.FieldLogger)
    cacheKey := "refresh_token:" + refreshToken
    err := s.Cache.Delete(ctx, cacheKey)
    if err != nil {
        logger.WithError(err).Warn("Failed to delete refresh token")
        return appErrors.NewUnauthorizedError("Invalid or expired refresh token", err)
    }
    return nil
}
