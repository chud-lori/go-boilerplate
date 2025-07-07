package services_test

import (
	"context"
	"errors"
	"testing"

	"time"

	"github.com/chud-lori/go-boilerplate/domain/entities"
	"github.com/chud-lori/go-boilerplate/domain/services"
	"github.com/chud-lori/go-boilerplate/mocks"
	appErrors "github.com/chud-lori/go-boilerplate/pkg/errors"
	"github.com/chud-lori/go-boilerplate/pkg/logger"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthService_SignIn_Success(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	mockDB := new(mocks.MockDatabase)
	mockRepo := new(mocks.MockUserRepository)
	mockMailSrv := new(mocks.MockMailService)
	mockEnc := new(mocks.MockEncryptor)
	mockToken := new(mocks.MockTokenManager)
	mockTx := new(mocks.MockTransaction)

	service := &services.AuthServiceImpl{
		DB:             mockDB,
		UserRepository: mockRepo,
		MailService:    mockMailSrv,
		Encryptor:      mockEnc,
		TokenManager:   mockToken,
		CtxTimeout:     2 * time.Second,
	}

	userUUID := uuid.New()
	foundUser := &entities.User{ID: userUUID, Email: "user@mail.com", Password: "hashpassword"}
	mockUser := &entities.User{ID: userUUID, Email: "user@mail.com", Password: "password1234"}

	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockRepo.On("FindByEmail", mock.Anything, mockTx, mockUser.Email).Return(foundUser, nil)
	mockEnc.On("CompareHash", foundUser.Password, mockUser.Password).Return(nil)
	mockToken.On("GenerateToken", mockUser.ID.String()).Return("generatedtoken", nil)
	mockTx.On("Commit").Return(nil)
	mockMailSrv.On("SendSignInNotification", mock.Anything, foundUser.Email, "User logged in just now").Return(nil)

	_, token, err := service.SignIn(ctx, mockUser)

	assert.NoError(t, err)
	assert.Equal(t, "generatedtoken", token)

	mockDB.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockMailSrv.AssertExpectations(t)
	mockEnc.AssertExpectations(t)
	mockToken.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestAuthService_SignIn_Failed(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	mockDB := new(mocks.MockDatabase)
	mockRepo := new(mocks.MockUserRepository)
	mockEnc := new(mocks.MockEncryptor)
	mockToken := new(mocks.MockTokenManager)
	mockTx := new(mocks.MockTransaction)

	service := &services.AuthServiceImpl{
		DB:             mockDB,
		UserRepository: mockRepo,
		Encryptor:      mockEnc,
		TokenManager:   mockToken,
		CtxTimeout:     2 * time.Second,
	}

	mockUser := &entities.User{ID: uuid.New(), Email: "user@mail.com", Password: "password1234"}

	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockRepo.On("FindByEmail", mock.Anything, mockTx, mockUser.Email).Return(nil, errors.New("Not found"))
	mockTx.On("Rollback").Return(nil)

	user, token, err := service.SignIn(ctx, mockUser)

	assert.Error(t, err)
	appErr, ok := err.(*appErrors.AppError)
	assert.True(t, ok)
	assert.Equal(t, "Unauthorized", appErr.Message)
	assert.Equal(t, 401, appErr.StatusCode)
	assert.Nil(t, user)
	assert.Equal(t, "", token)

	mockDB.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockEnc.AssertNotCalled(t, "CompareHash")
	mockToken.AssertNotCalled(t, "GenerateToken")
	mockTx.AssertExpectations(t)
}

func TestAuthService_SignUp_Success(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	mockDB := new(mocks.MockDatabase)
	mockRepo := new(mocks.MockUserRepository)
	mockEnc := new(mocks.MockEncryptor)
	mockToken := new(mocks.MockTokenManager)
	mockTx := new(mocks.MockTransaction)

	service := &services.AuthServiceImpl{
		DB:             mockDB,
		UserRepository: mockRepo,
		Encryptor:      mockEnc,
		TokenManager:   mockToken,
		CtxTimeout:     2 * time.Second,
	}

	mockUser := &entities.User{ID: uuid.New(), Email: "user@mail.com", Password: "password1234"}

	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockRepo.On("FindByEmail", mock.Anything, mockTx, mockUser.Email).Return(nil, errors.New("Not found"))
	mockEnc.On("HashPassword", mockUser.Password).Return("hashpassword", nil)
	mockRepo.On("Save", mock.Anything, mockTx, mockUser).Return(mockUser, nil)
	mockToken.On("GenerateToken", mockUser.ID.String()).Return("generatedtoken", nil)
	mockTx.On("Commit").Return(nil)

	_, token, err := service.SignUp(ctx, mockUser)

	assert.NoError(t, err)
	assert.Equal(t, "generatedtoken", token)

	mockDB.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockEnc.AssertExpectations(t)
	mockToken.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestAuthService_SignUp_Failed(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	mockDB := new(mocks.MockDatabase)
	mockRepo := new(mocks.MockUserRepository)
	mockEnc := new(mocks.MockEncryptor)
	mockToken := new(mocks.MockTokenManager)
	mockTx := new(mocks.MockTransaction)

	service := &services.AuthServiceImpl{
		DB:             mockDB,
		UserRepository: mockRepo,
		Encryptor:      mockEnc,
		TokenManager:   mockToken,
		CtxTimeout:     2 * time.Second,
	}

	mockUser := &entities.User{ID: uuid.New(), Email: "user@mail.com", Password: "password1234"}

	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockRepo.On("FindByEmail", mock.Anything, mockTx, mockUser.Email).Return(mockUser, nil)
	mockTx.On("Rollback").Return(nil)

	user, token, err := service.SignUp(ctx, mockUser)

	assert.Error(t, err)
	appErr, ok := err.(*appErrors.AppError)
	assert.True(t, ok)
	assert.Equal(t, "Email already exist", appErr.Message)
	assert.Equal(t, 400, appErr.StatusCode)
	assert.Nil(t, user)
	assert.Equal(t, "", token)

	mockDB.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Save")
	mockEnc.AssertNotCalled(t, "HashPassword")
	mockToken.AssertNotCalled(t, "GenerateToken")
	mockTx.AssertExpectations(t)
}
