package services_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/chud-lori/go-boilerplate/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/chud-lori/go-boilerplate/domain/entities"
	"github.com/chud-lori/go-boilerplate/domain/services"
	appErrors "github.com/chud-lori/go-boilerplate/pkg/errors"
	"github.com/chud-lori/go-boilerplate/pkg/logger"
	"github.com/sirupsen/logrus"
)

func TestUserService_Save_Success(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	mockDB := new(mocks.MockDatabase)
	mockRepo := new(mocks.MockUserRepository)
	mockEnc := new(mocks.MockEncryptor)
	mockCache := new(mocks.MockCache)
	mockTx := new(mocks.MockTransaction)

	service := &services.UserServiceImpl{
		DB:             mockDB,
		UserRepository: mockRepo,
		Encryptor:      mockEnc,
		Cache:          mockCache,
		CtxTimeout:     2 * time.Second,
	}

	user := &entities.User{ID: uuid.New(), Email: "user@mail.com", Password: "password1234"}

	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockEnc.On("HashPassword", user.Password).Return("hashed", nil)
	mockRepo.On("Save", mock.Anything, mockTx, user).Return(user, nil)
	mockTx.On("Commit").Return(nil)

	result, err := service.Save(ctx, user)

	assert.NoError(t, err)
	assert.Equal(t, "user@mail.com", result.Email)

	mockDB.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestUserService_Save_Failed(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	mockDB := new(mocks.MockDatabase)
	mockRepo := new(mocks.MockUserRepository)
	mockEnc := new(mocks.MockEncryptor)
	mockCache := new(mocks.MockCache)
	mockTx := new(mocks.MockTransaction)

	service := &services.UserServiceImpl{
		DB:             mockDB,
		UserRepository: mockRepo,
		Encryptor:      mockEnc,
		Cache:          mockCache,
		CtxTimeout:     2 * time.Second,
	}

	user := &entities.User{ID: uuid.New(), Email: "user@mail.com", Password: "password1234"}

	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockEnc.On("HashPassword", user.Password).Return("hashed", nil)
	mockRepo.On("Save", mock.Anything, mockTx, user).Return(nil, errors.New("Error"))
	mockTx.On("Rollback").Return(nil)

	result, err := service.Save(ctx, user)

	assert.Error(t, err)
	assert.Nil(t, result)

	mockDB.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestUserService_Save_FailedCommit(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	mockDB := new(mocks.MockDatabase)
	mockRepo := new(mocks.MockUserRepository)
	mockEnc := new(mocks.MockEncryptor)
	mockCache := new(mocks.MockCache)
	mockTx := new(mocks.MockTransaction)

	service := &services.UserServiceImpl{
		DB:             mockDB,
		UserRepository: mockRepo,
		Encryptor:      mockEnc,
		Cache:          mockCache,
		CtxTimeout:     2 * time.Second,
	}

	user := &entities.User{ID: uuid.New(), Email: "user@mail.com", Password: "password1234"}

	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockEnc.On("HashPassword", user.Password).Return("hashed", nil)
	mockRepo.On("Save", mock.Anything, mockTx, user).Return(user, nil)
	mockTx.On("Commit").Return(errors.New("Error"))
	mockTx.On("Rollback").Return(nil)

	result, err := service.Save(ctx, user)

	assert.Error(t, err)
	assert.Nil(t, result)

	mockDB.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestUserService_Update_Success(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	mockDB := new(mocks.MockDatabase)
	mockRepo := new(mocks.MockUserRepository)
	mockEnc := new(mocks.MockEncryptor)
	mockCache := new(mocks.MockCache)
	mockTx := new(mocks.MockTransaction)

	service := &services.UserServiceImpl{
		DB:             mockDB,
		UserRepository: mockRepo,
		Encryptor:      mockEnc,
		Cache:          mockCache,
		CtxTimeout:     2 * time.Second,
	}

	user := &entities.User{ID: uuid.New(), Email: "user@mail.com", Password: "password1234"}

	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockEnc.On("HashPassword", user.Password).Return("hashed", nil)
	mockRepo.On("Update", mock.Anything, mockTx, user).Return(user, nil)
	mockTx.On("Commit").Return(nil)

	result, err := service.Update(ctx, user)

	assert.NoError(t, err)
	assert.Equal(t, "user@mail.com", result.Email)

	mockDB.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestUserService_Update_UserNotFound(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	mockDB := new(mocks.MockDatabase)
	mockRepo := new(mocks.MockUserRepository)
	mockEnc := new(mocks.MockEncryptor)
	mockCache := new(mocks.MockCache)
	mockTx := new(mocks.MockTransaction)

	service := &services.UserServiceImpl{
		DB:             mockDB,
		UserRepository: mockRepo,
		Encryptor:      mockEnc,
		Cache:          mockCache,
		CtxTimeout:     2 * time.Second,
	}

	user := &entities.User{ID: uuid.New(), Email: "user@mail.com", Password: "password1234"}

	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockEnc.On("HashPassword", user.Password).Return("hashed", nil)
	mockRepo.On("Update", mock.Anything, mockTx, user).Return(nil, appErrors.ErrUserNotFound)
	mockTx.On("Rollback").Return(nil)

	result, err := service.Update(ctx, user)

	assert.Error(t, err)
	appErr, ok := err.(*appErrors.AppError)
	assert.True(t, ok)
	assert.Equal(t, "User not found", appErr.Message)
	assert.Equal(t, 404, appErr.StatusCode)

	assert.Nil(t, result)

	mockDB.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
	mockEnc.AssertNotCalled(t, "HashPassword")
}

func TestUserService_Delete_Success(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	mockDB := new(mocks.MockDatabase)
	mockRepo := new(mocks.MockUserRepository)
	mockEnc := new(mocks.MockEncryptor)
	mockCache := new(mocks.MockCache)
	mockTx := new(mocks.MockTransaction)

	service := &services.UserServiceImpl{
		DB:             mockDB,
		UserRepository: mockRepo,
		Encryptor:      mockEnc,
		Cache:          mockCache,
		CtxTimeout:     2 * time.Second,
	}

	userId := "ad24a17d-2925-4aa8-b077-d358a0788df7"

	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockRepo.On("Delete", mock.Anything, mockTx, userId).Return(nil)
	mockTx.On("Commit").Return(nil)

	err := service.Delete(ctx, userId)

	assert.NoError(t, err)

	mockDB.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestUserService_Delete_UserNotFound(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	mockDB := new(mocks.MockDatabase)
	mockRepo := new(mocks.MockUserRepository)
	mockEnc := new(mocks.MockEncryptor)
	mockCache := new(mocks.MockCache)
	mockTx := new(mocks.MockTransaction)

	service := &services.UserServiceImpl{
		DB:             mockDB,
		UserRepository: mockRepo,
		Encryptor:      mockEnc,
		Cache:          mockCache,
		CtxTimeout:     2 * time.Second,
	}

	userId := "ad24a17d-2925-4aa8-b077-d358a0788df7"

	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockRepo.On("Delete", mock.Anything, mockTx, userId).Return(appErrors.ErrUserNotFound)
	mockTx.On("Rollback").Return(nil)

	err := service.Delete(ctx, userId)

	assert.Error(t, err)
	appErr, ok := err.(*appErrors.AppError)
	assert.True(t, ok)
	assert.Equal(t, "User not found", appErr.Message)
	assert.Equal(t, 404, appErr.StatusCode)

	mockDB.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestUserService_FindById_Success(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	mockDB := new(mocks.MockDatabase)
	mockRepo := new(mocks.MockUserRepository)
	mockEnc := new(mocks.MockEncryptor)
	mockCache := new(mocks.MockCache)
	mockTx := new(mocks.MockTransaction)

	service := &services.UserServiceImpl{
		DB:             mockDB,
		UserRepository: mockRepo,
		Encryptor:      mockEnc,
		Cache:          mockCache,
		CtxTimeout:     2 * time.Second,
	}

	user := &entities.User{ID: uuid.New(), Email: "user@mail.com", Password: "password1234"}

	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockRepo.On("FindById", mock.Anything, mockTx, user.ID.String()).Return(user, nil)
	mockTx.On("Commit").Return(nil)

	result, err := service.FindById(ctx, user.ID.String())

	assert.NoError(t, err)
	assert.NotNil(t, result)

	mockDB.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestUserService_FindById_UserNotFound(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	mockDB := new(mocks.MockDatabase)
	mockRepo := new(mocks.MockUserRepository)
	mockEnc := new(mocks.MockEncryptor)
	mockCache := new(mocks.MockCache)
	mockTx := new(mocks.MockTransaction)

	service := &services.UserServiceImpl{
		DB:             mockDB,
		UserRepository: mockRepo,
		Encryptor:      mockEnc,
		Cache:          mockCache,
		CtxTimeout:     2 * time.Second,
	}

	userId := "ad24a17d-2925-4aa8-b077-d358a0788df7"

	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockRepo.On("FindById", mock.Anything, mockTx, userId).Return(nil, appErrors.ErrUserNotFound)
	mockTx.On("Rollback").Return(nil)

	result, err := service.FindById(ctx, userId)

	assert.Nil(t, result)
	assert.Error(t, err)
	appErr, ok := err.(*appErrors.AppError)
	assert.True(t, ok)
	assert.Equal(t, "User not found", appErr.Message)
	assert.Equal(t, 404, appErr.StatusCode)

	mockDB.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestUserService_FindAll_Success(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	mockDB := new(mocks.MockDatabase)
	mockRepo := new(mocks.MockUserRepository)
	mockEnc := new(mocks.MockEncryptor)
	mockCache := new(mocks.MockCache)
	mockTx := new(mocks.MockTransaction)

	service := &services.UserServiceImpl{
		DB:             mockDB,
		UserRepository: mockRepo,
		Encryptor:      mockEnc,
		Cache:          mockCache,
		CtxTimeout:     2 * time.Second,
	}

	listUsers := []*entities.User{
		{
			ID:        uuid.New(),
			Email:     "user1@mail.com",
			Password:  "pass1",
			CreatedAt: time.Date(2023, time.January, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.New(),
			Email:     "user2@mail.com",
			Password:  "pass2",
			CreatedAt: time.Date(2023, time.February, 20, 11, 30, 0, 0, time.UTC),
		},
	}

	mockCache.On("Get", mock.Anything, "users").Return("", nil).Once()
	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockRepo.On("FindAll", mock.Anything, mockTx).Return(listUsers, nil)
	mockCache.On("Set", mock.Anything, "users", mock.AnythingOfType("[]uint8"), mock.AnythingOfType("time.Duration")).Return(nil).Once()
	mockTx.On("Commit").Return(nil)

	result, err := service.FindAll(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	mockCache.AssertExpectations(t)
	mockDB.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestUserService_FindAll_SuccessCache(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	mockDB := new(mocks.MockDatabase)
	mockRepo := new(mocks.MockUserRepository)
	mockEnc := new(mocks.MockEncryptor)
	mockCache := new(mocks.MockCache)
	mockTx := new(mocks.MockTransaction)

	service := &services.UserServiceImpl{
		DB:             mockDB,
		UserRepository: mockRepo,
		Encryptor:      mockEnc,
		Cache:          mockCache,
		CtxTimeout:     2 * time.Second,
	}

	listUsers := []*entities.User{
		{
			ID:        uuid.New(),
			Email:     "user1@mail.com",
			Password:  "pass1",
			CreatedAt: time.Date(2023, time.January, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.New(),
			Email:     "user2@mail.com",
			Password:  "pass2",
			CreatedAt: time.Date(2023, time.February, 20, 11, 30, 0, 0, time.UTC),
		},
	}

	listUsersJson, err := json.Marshal(listUsers)
	assert.NoError(t, err, "Failed to marshal listUsers for test setup")

	mockCache.On("Get", mock.Anything, "users").Return(string(listUsersJson), nil).Once()

	result, err := service.FindAll(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	mockCache.AssertExpectations(t)

	mockDB.AssertNotCalled(t, "BeginTx", mock.Anything, mock.Anything)   // Assuming BeginTx is the only DB method
	mockRepo.AssertNotCalled(t, "FindAll", mock.Anything, mock.Anything) // Assuming FindAll is the only Repo method
	mockCache.AssertNotCalled(t, "Set")
	mockTx.AssertNotCalled(t, "Commit") // Assuming Commit is the only Tx method
}
