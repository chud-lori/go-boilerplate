package services_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/chud-lori/go-boilerplate/mocks"
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
	mockTx := new(mocks.MockTransaction)

	service := &services.UserServiceImpl{
		DB:             mockDB,
		UserRepository: mockRepo,
		CtxTimeout:     2 * time.Second,
	}

	user := &entities.User{Id: "", Email: "user@mail.com", Passcode: ""}

	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil)
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
	mockTx := new(mocks.MockTransaction)

	service := &services.UserServiceImpl{
		DB:             mockDB,
		UserRepository: mockRepo,
		CtxTimeout:     2 * time.Second,
	}

	user := &entities.User{Id: "", Email: "user@mail.com", Passcode: ""}

	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil)
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
	mockTx := new(mocks.MockTransaction)

	service := &services.UserServiceImpl{
		DB:             mockDB,
		UserRepository: mockRepo,
		CtxTimeout:     2 * time.Second,
	}

	user := &entities.User{Id: "", Email: "user@mail.com", Passcode: ""}

	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil)
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
	mockTx := new(mocks.MockTransaction)

	service := &services.UserServiceImpl{
		DB:             mockDB,
		UserRepository: mockRepo,
		CtxTimeout:     2 * time.Second,
	}

	user := &entities.User{Id: "", Email: "user@mail.com", Passcode: ""}

	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil)
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
	mockTx := new(mocks.MockTransaction)

	service := &services.UserServiceImpl{
		DB:             mockDB,
		UserRepository: mockRepo,
		CtxTimeout:     2 * time.Second,
	}

	user := &entities.User{Id: "", Email: "user@mail.com", Passcode: ""}

	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil)
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
}

func TestUserService_Delete_Success(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	mockDB := new(mocks.MockDatabase)
	mockRepo := new(mocks.MockUserRepository)
	mockTx := new(mocks.MockTransaction)

	service := &services.UserServiceImpl{
		DB:             mockDB,
		UserRepository: mockRepo,
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
	mockTx := new(mocks.MockTransaction)

	service := &services.UserServiceImpl{
		DB:             mockDB,
		UserRepository: mockRepo,
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
	mockTx := new(mocks.MockTransaction)

	service := &services.UserServiceImpl{
		DB:             mockDB,
		UserRepository: mockRepo,
		CtxTimeout:     2 * time.Second,
	}

	userId := "ad24a17d-2925-4aa8-b077-d358a0788df7"
	user := &entities.User{Id: userId, Email: "user@mail.com", Passcode: ""}

	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockRepo.On("FindById", mock.Anything, mockTx, userId).Return(user, nil)
	mockTx.On("Commit").Return(nil)

	result, err := service.FindById(ctx, userId)

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
	mockTx := new(mocks.MockTransaction)

	service := &services.UserServiceImpl{
		DB:             mockDB,
		UserRepository: mockRepo,
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
	mockTx := new(mocks.MockTransaction)
	mockCache := new(mocks.MockCache)

	service := &services.UserServiceImpl{
		DB:             mockDB,
		UserRepository: mockRepo,
		Cache:          mockCache,
		CtxTimeout:     2 * time.Second,
	}

	listUsers := []*entities.User{
		{
			Id:        "a234f98c-3239-4c34-8ad8-f63e41bb20c8", // Define userId directly here
			Email:     "user1@mail.com",
			Passcode:  "pass1",
			CreatedAt: time.Date(2023, time.January, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			Id:        "b567g89d-4321-5d67-9fg0-g76h54ij32k1",
			Email:     "user2@mail.com",
			Passcode:  "pass2",
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
	mockTx := new(mocks.MockTransaction)
	mockCache := new(mocks.MockCache)

	service := &services.UserServiceImpl{
		DB:             mockDB,
		UserRepository: mockRepo,
		Cache:          mockCache,
		CtxTimeout:     2 * time.Second,
	}

	listUsers := []*entities.User{
		{
			Id:        "a234f98c-3239-4c34-8ad8-f63e41bb20c8", // Define userId directly here
			Email:     "user1@mail.com",
			Passcode:  "pass1",
			CreatedAt: time.Date(2023, time.January, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			Id:        "b567g89d-4321-5d67-9fg0-g76h54ij32k1",
			Email:     "user2@mail.com",
			Passcode:  "pass2",
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
