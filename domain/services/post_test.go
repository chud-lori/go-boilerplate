package services_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"time"

	"github.com/chud-lori/go-boilerplate/mocks"
	appErrors "github.com/chud-lori/go-boilerplate/pkg/errors"
	"github.com/chud-lori/go-boilerplate/pkg/logger"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/chud-lori/go-boilerplate/domain/entities"
	"github.com/chud-lori/go-boilerplate/domain/services"
)

func TestPostService_Create_Success(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	mockDB := new(mocks.MockDatabase)
	mockPostRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockCache := new(mocks.MockCache) // Not used in Create, but good to include in service struct
	mockTx := new(mocks.MockTransaction)

	service := &services.PostServiceImpl{
		DB:             mockDB,
		PostRepository: mockPostRepo,
		UserRepository: mockUserRepo,
		Cache:          mockCache,
		CtxTimeout:     2 * time.Second,
	}

	user := &entities.User{ID: uuid.New(), Email: "author@example.com"}
	post := &entities.Post{
		ID:    uuid.New(),
		User:  user,
		Title: "Test Title",
		Body:  "Test Content",
	}

	// Setup mock expectations
	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil).Once()
	mockTx.On("Commit").Return(nil).Once()
	mockTx.On("Rollback").Return(nil).Maybe() // Rollback might be called on error
	mockUserRepo.On("FindById", mock.Anything, mockTx, user.ID.String()).Return(user, nil).Once()
	mockPostRepo.On("Save", mock.Anything, mockTx, post).Return(post, nil).Once()

	// Call the service method
	result, err := service.Create(ctx, post)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, post.ID, result.ID)
	assert.Equal(t, post.Title, result.Title)

	// Verify mock expectations
	mockDB.AssertExpectations(t)
	mockPostRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestPostService_Create_BeginTxError(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	mockDB := new(mocks.MockDatabase)
	mockPostRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockCache := new(mocks.MockCache)
	mockTx := new(mocks.MockTransaction) // Still need this for BeginTx return type, even if it's nil

	service := &services.PostServiceImpl{
		DB:             mockDB,
		PostRepository: mockPostRepo,
		UserRepository: mockUserRepo,
		Cache:          mockCache,
		CtxTimeout:     2 * time.Second,
	}

	author := &entities.User{
		ID: uuid.New(),
	}
	post := &entities.Post{
		ID:    uuid.New(),
		User:  author,
		Title: "Test Title",
		Body:  "Test Content",
	}

	// Setup mock expectations for BeginTx error
	expectedErr := errors.New("failed to begin transaction")
	mockDB.On("BeginTx", mock.Anything).Return(mockTx, expectedErr).Once()

	// Call the service method
	result, err := service.Create(ctx, post)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedErr, err)

	// Verify mock expectations
	mockDB.AssertExpectations(t)
	mockPostRepo.AssertNotCalled(t, "Save")     // Should not be called
	mockUserRepo.AssertNotCalled(t, "FindById") // Should not be called
}

func TestPostService_Create_AuthorNotFound(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	mockDB := new(mocks.MockDatabase)
	mockPostRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockCache := new(mocks.MockCache)
	mockTx := new(mocks.MockTransaction)

	service := &services.PostServiceImpl{
		DB:             mockDB,
		PostRepository: mockPostRepo,
		UserRepository: mockUserRepo,
		Cache:          mockCache,
		CtxTimeout:     2 * time.Second,
	}

	author := &entities.User{
		ID: uuid.New(),
	}
	post := &entities.Post{
		ID:    uuid.New(),
		User:  author,
		Title: "Test Title",
		Body:  "Test Content",
	}

	// Setup mock expectations
	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil).Once()
	mockTx.On("Rollback").Return(nil).Once() // Rollback expected on error
	mockUserRepo.On("FindById", mock.Anything, mockTx, author.ID.String()).Return(nil, appErrors.ErrUserNotFound).Once()

	// Call the service method
	result, err := service.Create(ctx, post)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)

	appErr, ok := err.(*appErrors.AppError)
	assert.True(t, ok)
	assert.Equal(t, "Author not found", appErr.Message)
	assert.Equal(t, 404, appErr.StatusCode)

	// Verify mock expectations
	mockDB.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockPostRepo.AssertNotCalled(t, "Save") // Should not be called
	mockTx.AssertExpectations(t)
}

func TestPostService_Create_SavePostError(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	mockDB := new(mocks.MockDatabase)
	mockPostRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockCache := new(mocks.MockCache)
	mockTx := new(mocks.MockTransaction)

	service := &services.PostServiceImpl{
		DB:             mockDB,
		PostRepository: mockPostRepo,
		UserRepository: mockUserRepo,
		Cache:          mockCache,
		CtxTimeout:     2 * time.Second,
	}

	author := &entities.User{ID: uuid.New(), Email: "author@example.com"}
	post := &entities.Post{
		ID:    uuid.New(),
		User:  author,
		Title: "Test Title",
		Body:  "Test Content",
	}

	// Setup mock expectations
	expectedSaveErr := errors.New("database save error")
	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil).Once()
	mockTx.On("Rollback").Return(nil).Once() // Rollback expected on error
	mockUserRepo.On("FindById", mock.Anything, mockTx, post.User.ID.String()).Return(author, nil).Once()
	mockPostRepo.On("Save", mock.Anything, mockTx, post).Return(nil, expectedSaveErr).Once()

	// Call the service method
	result, err := service.Create(ctx, post)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedSaveErr, err)

	// Verify mock expectations
	mockDB.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockPostRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestPostService_Create_CommitError(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	mockDB := new(mocks.MockDatabase)
	mockPostRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockCache := new(mocks.MockCache)
	mockTx := new(mocks.MockTransaction)

	service := &services.PostServiceImpl{
		DB:             mockDB,
		PostRepository: mockPostRepo,
		UserRepository: mockUserRepo,
		Cache:          mockCache,
		CtxTimeout:     2 * time.Second,
	}

	author := &entities.User{ID: uuid.New(), Email: "author@example.com"}
	post := &entities.Post{
		ID:    uuid.New(),
		User:  author,
		Title: "Test Title",
		Body:  "Test Content",
	}

	// Setup mock expectations
	expectedCommitErr := errors.New("failed to commit transaction")
	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil).Once()
	mockTx.On("Commit").Return(expectedCommitErr).Once()
	mockTx.On("Rollback").Return(nil).Once() // Rollback expected on commit error
	mockUserRepo.On("FindById", mock.Anything, mockTx, post.User.ID.String()).Return(author, nil).Once()
	mockPostRepo.On("Save", mock.Anything, mockTx, post).Return(post, nil).Once()

	// Call the service method
	result, err := service.Create(ctx, post)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedCommitErr, err)

	// Verify mock expectations
	mockDB.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockPostRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestPostService_Update_Success(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	mockDB := new(mocks.MockDatabase)
	mockPostRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository) // Not used in Update, but part of service
	mockCache := new(mocks.MockCache)
	mockTx := new(mocks.MockTransaction)

	service := &services.PostServiceImpl{
		DB:             mockDB,
		PostRepository: mockPostRepo,
		UserRepository: mockUserRepo,
		Cache:          mockCache,
		CtxTimeout:     2 * time.Second,
	}

	postID := uuid.New()
	post := &entities.Post{
		ID:    postID,
		Title: "Updated Title",
		Body:  "Updated Content",
	}

	// Setup mock expectations
	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil).Once()
	mockTx.On("Commit").Return(nil).Once()
	mockTx.On("Rollback").Return(nil).Maybe()
	mockPostRepo.On("Update", mock.Anything, mockTx, post).Return(post, nil).Once()
	mockCache.On("InvalidateByPrefix", mock.Anything, "posts:").Return(nil).Once()

	// Call the service method
	result, err := service.Update(ctx, post)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, post.ID, result.ID)
	assert.Equal(t, post.Title, result.Title)

	// Verify mock expectations
	mockDB.AssertExpectations(t)
	mockPostRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestPostService_Update_BeginTxError(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	mockDB := new(mocks.MockDatabase)
	mockPostRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockCache := new(mocks.MockCache)
	mockTx := new(mocks.MockTransaction)

	service := &services.PostServiceImpl{
		DB:             mockDB,
		PostRepository: mockPostRepo,
		UserRepository: mockUserRepo,
		Cache:          mockCache,
		CtxTimeout:     2 * time.Second,
	}

	post := &entities.Post{ID: uuid.New()}

	expectedErr := errors.New("failed to begin transaction")
	mockDB.On("BeginTx", mock.Anything).Return(mockTx, expectedErr).Once()

	result, err := service.Update(ctx, post)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedErr, err)

	mockDB.AssertExpectations(t)
	mockPostRepo.AssertNotCalled(t, "Update")
}

func TestPostService_Update_PostNotFound(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	mockDB := new(mocks.MockDatabase)
	mockPostRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockCache := new(mocks.MockCache)
	mockTx := new(mocks.MockTransaction)

	service := &services.PostServiceImpl{
		DB:             mockDB,
		PostRepository: mockPostRepo,
		UserRepository: mockUserRepo,
		Cache:          mockCache,
		CtxTimeout:     2 * time.Second,
	}

	postID := uuid.New()
	post := &entities.Post{ID: postID}

	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil).Once()
	mockTx.On("Rollback").Return(nil).Once()
	mockPostRepo.On("Update", mock.Anything, mockTx, post).Return(nil, appErrors.ErrDataNotFound).Once()

	result, err := service.Update(ctx, post)

	assert.Error(t, err)
	assert.Nil(t, result)
	appErr, ok := err.(*appErrors.AppError)
	assert.True(t, ok)
	assert.Equal(t, "Post not found", appErr.Message)
	assert.Equal(t, 404, appErr.StatusCode)

	mockDB.AssertExpectations(t)
	mockPostRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
	mockCache.AssertNotCalled(t, "InvalidateByPrefix") // Should not be called on failure
}

func TestPostService_Update_GenericUpdateError(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	mockDB := new(mocks.MockDatabase)
	mockPostRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockCache := new(mocks.MockCache)
	mockTx := new(mocks.MockTransaction)

	service := &services.PostServiceImpl{
		DB:             mockDB,
		PostRepository: mockPostRepo,
		UserRepository: mockUserRepo,
		Cache:          mockCache,
		CtxTimeout:     2 * time.Second,
	}

	postID := uuid.New()
	post := &entities.Post{ID: postID}

	expectedErr := errors.New("database update error")
	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil).Once()
	mockTx.On("Rollback").Return(nil).Once()
	mockPostRepo.On("Update", mock.Anything, mockTx, post).Return(nil, expectedErr).Once()

	result, err := service.Update(ctx, post)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedErr, err)

	mockDB.AssertExpectations(t)
	mockPostRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
	mockCache.AssertNotCalled(t, "InvalidateByPrefix")
}

func TestPostService_Update_CommitError(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	mockDB := new(mocks.MockDatabase)
	mockPostRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockCache := new(mocks.MockCache)
	mockTx := new(mocks.MockTransaction)

	service := &services.PostServiceImpl{
		DB:             mockDB,
		PostRepository: mockPostRepo,
		UserRepository: mockUserRepo,
		Cache:          mockCache,
		CtxTimeout:     2 * time.Second,
	}

	postID := uuid.New()
	post := &entities.Post{ID: postID}

	expectedErr := errors.New("failed to commit transaction")
	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil).Once()
	mockTx.On("Commit").Return(expectedErr).Once()
	mockTx.On("Rollback").Return(nil).Once()
	mockPostRepo.On("Update", mock.Anything, mockTx, post).Return(post, nil).Once()

	result, err := service.Update(ctx, post)

	assert.Error(t, err)
	assert.Nil(t, result) // Result should be nil on commit error
	assert.Equal(t, expectedErr, err)

	mockDB.AssertExpectations(t)
	mockPostRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
	mockCache.AssertNotCalled(t, "InvalidateByPrefix") // Not called if commit fails
}

func TestPostService_Update_CacheInvalidationError(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	mockDB := new(mocks.MockDatabase)
	mockPostRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockCache := new(mocks.MockCache)
	mockTx := new(mocks.MockTransaction)

	service := &services.PostServiceImpl{
		DB:             mockDB,
		PostRepository: mockPostRepo,
		UserRepository: mockUserRepo,
		Cache:          mockCache,
		CtxTimeout:     2 * time.Second,
	}

	postID := uuid.New()
	post := &entities.Post{
		ID:    postID,
		Title: "Updated Title",
		Body:  "Updated Content",
	}

	// Setup mock expectations
	expectedCacheErr := errors.New("failed to invalidate cache")
	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil).Once()
	mockTx.On("Commit").Return(nil).Once()
	mockTx.On("Rollback").Return(nil).Maybe()
	mockPostRepo.On("Update", mock.Anything, mockTx, post).Return(post, nil).Once()
	mockCache.On("InvalidateByPrefix", mock.Anything, "posts:").Return(expectedCacheErr).Once()

	// Call the service method
	result, err := service.Update(ctx, post)

	// Assertions: The main operation should still succeed, error is just logged/warned
	assert.NoError(t, err) // This is crucial: cache error does not return an error from Update
	assert.NotNil(t, result)
	assert.Equal(t, post.ID, result.ID)

	// Verify mock expectations
	mockDB.AssertExpectations(t)
	mockPostRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestPostService_Delete_Success(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	mockDB := new(mocks.MockDatabase)
	mockPostRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockCache := new(mocks.MockCache)
	mockTx := new(mocks.MockTransaction)

	service := &services.PostServiceImpl{
		DB:             mockDB,
		PostRepository: mockPostRepo,
		UserRepository: mockUserRepo,
		Cache:          mockCache,
		CtxTimeout:     2 * time.Second,
	}

	postID := uuid.New()

	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil).Once()
	mockTx.On("Commit").Return(nil).Once()
	mockTx.On("Rollback").Return(nil).Maybe()
	mockPostRepo.On("Delete", mock.Anything, mockTx, postID).Return(nil).Once()

	err := service.Delete(ctx, postID)

	assert.NoError(t, err)

	mockDB.AssertExpectations(t)
	mockPostRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestPostService_Delete_BeginTxError(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	mockDB := new(mocks.MockDatabase)
	mockPostRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockCache := new(mocks.MockCache)
	mockTx := new(mocks.MockTransaction)

	service := &services.PostServiceImpl{
		DB:             mockDB,
		PostRepository: mockPostRepo,
		UserRepository: mockUserRepo,
		Cache:          mockCache,
		CtxTimeout:     2 * time.Second,
	}

	postID := uuid.New()
	expectedErr := errors.New("failed to begin transaction")
	mockDB.On("BeginTx", mock.Anything).Return(mockTx, expectedErr).Once()

	err := service.Delete(ctx, postID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)

	mockDB.AssertExpectations(t)
	mockPostRepo.AssertNotCalled(t, "Delete")
}

func TestPostService_Delete_PostNotFound(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	mockDB := new(mocks.MockDatabase)
	mockPostRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockCache := new(mocks.MockCache)
	mockTx := new(mocks.MockTransaction)

	service := &services.PostServiceImpl{
		DB:             mockDB,
		PostRepository: mockPostRepo,
		UserRepository: mockUserRepo,
		Cache:          mockCache,
		CtxTimeout:     2 * time.Second,
	}

	postID := uuid.New()

	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil).Once()
	mockTx.On("Rollback").Return(nil).Once()
	mockPostRepo.On("Delete", mock.Anything, mockTx, postID).Return(appErrors.ErrDataNotFound).Once()

	err := service.Delete(ctx, postID)

	assert.Error(t, err)

	appErr, ok := err.(*appErrors.AppError)
	assert.True(t, ok)
	assert.Equal(t, "Post not found", appErr.Message)
	assert.Equal(t, 404, appErr.StatusCode)

	mockDB.AssertExpectations(t)
	mockPostRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestPostService_Delete_GenericDeleteError(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	mockDB := new(mocks.MockDatabase)
	mockPostRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockCache := new(mocks.MockCache)
	mockTx := new(mocks.MockTransaction)

	service := &services.PostServiceImpl{
		DB:             mockDB,
		PostRepository: mockPostRepo,
		UserRepository: mockUserRepo,
		Cache:          mockCache,
		CtxTimeout:     2 * time.Second,
	}

	postID := uuid.New()
	expectedErr := errors.New("database delete error")

	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil).Once()
	mockTx.On("Rollback").Return(nil).Once()
	mockPostRepo.On("Delete", mock.Anything, mockTx, postID).Return(expectedErr).Once()

	err := service.Delete(ctx, postID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)

	mockDB.AssertExpectations(t)
	mockPostRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestPostService_Delete_CommitError(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	mockDB := new(mocks.MockDatabase)
	mockPostRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockCache := new(mocks.MockCache)
	mockTx := new(mocks.MockTransaction)

	service := &services.PostServiceImpl{
		DB:             mockDB,
		PostRepository: mockPostRepo,
		UserRepository: mockUserRepo,
		Cache:          mockCache,
		CtxTimeout:     2 * time.Second,
	}

	postID := uuid.New()
	expectedErr := errors.New("failed to commit transaction")

	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil).Once()
	mockPostRepo.On("Delete", mock.Anything, mockTx, postID).Return(nil).Once()
	mockTx.On("Commit").Return(expectedErr).Once()
	mockTx.On("Rollback").Return(nil).Once()

	err := service.Delete(ctx, postID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)

	mockDB.AssertExpectations(t)
	mockPostRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestPostService_GetById_Success(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	mockDB := new(mocks.MockDatabase)
	mockPostRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockCache := new(mocks.MockCache)
	mockTx := new(mocks.MockTransaction)

	service := &services.PostServiceImpl{
		DB:             mockDB,
		PostRepository: mockPostRepo,
		UserRepository: mockUserRepo,
		Cache:          mockCache,
		CtxTimeout:     2 * time.Second,
	}

	postID := uuid.New()
	expectedPost := &entities.Post{ID: postID, Title: "Found Post"}

	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil).Once()
	mockTx.On("Commit").Return(nil).Once()
	mockTx.On("Rollback").Return(nil).Maybe()
	mockPostRepo.On("GetById", mock.Anything, mockTx, postID).Return(expectedPost, nil).Once()

	result, err := service.GetById(ctx, postID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedPost, result)

	mockDB.AssertExpectations(t)
	mockPostRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestPostService_GetById_BeginTxError(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	mockDB := new(mocks.MockDatabase)
	mockPostRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockCache := new(mocks.MockCache)
	mockTx := new(mocks.MockTransaction)

	service := &services.PostServiceImpl{
		DB:             mockDB,
		PostRepository: mockPostRepo,
		UserRepository: mockUserRepo,
		Cache:          mockCache,
		CtxTimeout:     2 * time.Second,
	}

	postID := uuid.New()
	expectedErr := errors.New("failed to begin transaction")
	mockDB.On("BeginTx", mock.Anything).Return(mockTx, expectedErr).Once()

	result, err := service.GetById(ctx, postID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedErr, err)

	mockDB.AssertExpectations(t)
	mockPostRepo.AssertNotCalled(t, "GetById")
}

func TestPostService_GetById_PostNotFound(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	mockDB := new(mocks.MockDatabase)
	mockPostRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockCache := new(mocks.MockCache)
	mockTx := new(mocks.MockTransaction)

	service := &services.PostServiceImpl{
		DB:             mockDB,
		PostRepository: mockPostRepo,
		UserRepository: mockUserRepo,
		Cache:          mockCache,
		CtxTimeout:     2 * time.Second,
	}

	postID := uuid.New()

	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil).Once()
	mockTx.On("Rollback").Return(nil).Once()
	mockPostRepo.On("GetById", mock.Anything, mockTx, postID).Return(nil, appErrors.ErrDataNotFound).Once()

	result, err := service.GetById(ctx, postID)

	assert.Error(t, err)
	assert.Nil(t, result)

	appErr, ok := err.(*appErrors.AppError)
	assert.True(t, ok)
	assert.Equal(t, "Post not found", appErr.Message)
	assert.Equal(t, 404, appErr.StatusCode)

	mockDB.AssertExpectations(t)
	mockPostRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestPostService_GetById_GenericGetError(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	mockDB := new(mocks.MockDatabase)
	mockPostRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockCache := new(mocks.MockCache)
	mockTx := new(mocks.MockTransaction)

	service := &services.PostServiceImpl{
		DB:             mockDB,
		PostRepository: mockPostRepo,
		UserRepository: mockUserRepo,
		Cache:          mockCache,
		CtxTimeout:     2 * time.Second,
	}

	postID := uuid.New()
	expectedErr := errors.New("database get error")

	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil).Once()
	mockTx.On("Rollback").Return(nil).Once()
	mockPostRepo.On("GetById", mock.Anything, mockTx, postID).Return(nil, expectedErr).Once()

	result, err := service.GetById(ctx, postID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedErr, err)

	mockDB.AssertExpectations(t)
	mockPostRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestPostService_GetById_CommitError(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	mockDB := new(mocks.MockDatabase)
	mockPostRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockCache := new(mocks.MockCache)
	mockTx := new(mocks.MockTransaction)

	service := &services.PostServiceImpl{
		DB:             mockDB,
		PostRepository: mockPostRepo,
		UserRepository: mockUserRepo,
		Cache:          mockCache,
		CtxTimeout:     2 * time.Second,
	}

	postID := uuid.New()
	expectedPost := &entities.Post{ID: postID, Title: "Found Post"}
	expectedErr := errors.New("failed to commit transaction")

	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil).Once()
	mockTx.On("Commit").Return(expectedErr).Once()
	mockTx.On("Rollback").Return(nil).Once()
	mockPostRepo.On("GetById", mock.Anything, mockTx, postID).Return(expectedPost, nil).Once()

	result, err := service.GetById(ctx, postID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedErr, err)

	mockDB.AssertExpectations(t)
	mockPostRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestPostService_GetAll_CacheHit(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	mockDB := new(mocks.MockDatabase) // Not used in cache hit scenario, but included
	mockPostRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockCache := new(mocks.MockCache)

	service := &services.PostServiceImpl{
		DB:             mockDB,
		PostRepository: mockPostRepo,
		UserRepository: mockUserRepo,
		Cache:          mockCache,
		CtxTimeout:     2 * time.Second,
	}

	search := "keyword"
	page := 1
	limit := 10
	expectedPosts := []entities.Post{
		{ID: uuid.New(), Title: "Post 1"},
		{ID: uuid.New(), Title: "Post 2"},
	}
	expectedPostsJSON, _ := json.Marshal(expectedPosts)

	queryParams := fmt.Sprintf("search=%s:page=%d:limit=%d", search, page, limit)
	hasher := sha256.New()
	hasher.Write([]byte(queryParams))
	cacheKey := "posts:" + hex.EncodeToString(hasher.Sum(nil))

	mockCache.On("Get", mock.Anything, cacheKey).Return(string(expectedPostsJSON), nil).Once()

	result, err := service.GetAll(ctx, search, page, limit)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, len(expectedPosts))
	assert.Equal(t, expectedPosts[0].Title, result[0].Title)

	mockCache.AssertExpectations(t)
	mockDB.AssertNotCalled(t, "BeginTx")      // DB operations should not occur on cache hit
	mockPostRepo.AssertNotCalled(t, "GetAll") // Repo operations should not occur on cache hit
}

func TestPostService_GetAll_CacheMiss_DBSucceeds(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	mockDB := new(mocks.MockDatabase)
	mockPostRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockCache := new(mocks.MockCache)
	mockTx := new(mocks.MockTransaction)

	service := &services.PostServiceImpl{
		DB:             mockDB,
		PostRepository: mockPostRepo,
		UserRepository: mockUserRepo,
		Cache:          mockCache,
		CtxTimeout:     2 * time.Second,
	}

	search := "keyword"
	page := 1
	limit := 10
	expectedPosts := []entities.Post{
		{ID: uuid.New(), Title: "Post 1"},
		{ID: uuid.New(), Title: "Post 2"},
	}
	expectedPostsJSON, _ := json.Marshal(expectedPosts)

	queryParams := fmt.Sprintf("search=%s:page=%d:limit=%d", search, page, limit)
	hasher := sha256.New()
	hasher.Write([]byte(queryParams))
	cacheKey := "posts:" + hex.EncodeToString(hasher.Sum(nil))

	// Simulate cache miss
	mockCache.On("Get", mock.Anything, cacheKey).Return("", errors.New("cache miss")).Once()
	// DB operations
	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil).Once()
	mockTx.On("Commit").Return(nil).Once()
	mockTx.On("Rollback").Return(nil).Maybe()
	mockPostRepo.On("GetAll", mock.Anything, mockTx, search, entities.PaginationParams{Page: page, Limit: limit}).Return(expectedPosts, nil).Once()
	// Cache set after successful DB fetch
	mockCache.On("Set", mock.Anything, cacheKey, expectedPostsJSON, 30*time.Second).Return(nil).Once()

	result, err := service.GetAll(ctx, search, page, limit)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, len(expectedPosts))
	assert.Equal(t, expectedPosts[0].Title, result[0].Title)

	mockCache.AssertExpectations(t)
	mockDB.AssertExpectations(t)
	mockPostRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestPostService_GetAll_BeginTxError(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	mockDB := new(mocks.MockDatabase)
	mockPostRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockCache := new(mocks.MockCache)
	mockTx := new(mocks.MockTransaction)

	service := &services.PostServiceImpl{
		DB:             mockDB,
		PostRepository: mockPostRepo,
		UserRepository: mockUserRepo,
		Cache:          mockCache,
		CtxTimeout:     2 * time.Second,
	}

	search := "keyword"
	page := 1
	limit := 10

	queryParams := fmt.Sprintf("search=%s:page=%d:limit=%d", search, page, limit)
	hasher := sha256.New()
	hasher.Write([]byte(queryParams))
	cacheKey := "posts:" + hex.EncodeToString(hasher.Sum(nil))

	// Simulate cache miss
	mockCache.On("Get", mock.Anything, cacheKey).Return("", errors.New("cache miss")).Once()
	// Simulate BeginTx error
	expectedErr := errors.New("failed to begin transaction")
	mockDB.On("BeginTx", mock.Anything).Return(mockTx, expectedErr).Once()

	result, err := service.GetAll(ctx, search, page, limit)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedErr, err)

	mockCache.AssertExpectations(t)
	mockDB.AssertExpectations(t)
	mockPostRepo.AssertNotCalled(t, "GetAll")
}

func TestPostService_GetAll_RepoError(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	mockDB := new(mocks.MockDatabase)
	mockPostRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockCache := new(mocks.MockCache)
	mockTx := new(mocks.MockTransaction)

	service := &services.PostServiceImpl{
		DB:             mockDB,
		PostRepository: mockPostRepo,
		UserRepository: mockUserRepo,
		Cache:          mockCache,
		CtxTimeout:     2 * time.Second,
	}

	search := "keyword"
	page := 1
	limit := 10

	queryParams := fmt.Sprintf("search=%s:page=%d:limit=%d", search, page, limit)
	hasher := sha256.New()
	hasher.Write([]byte(queryParams))
	cacheKey := "posts:" + hex.EncodeToString(hasher.Sum(nil))

	// Simulate cache miss
	mockCache.On("Get", mock.Anything, cacheKey).Return("", errors.New("cache miss")).Once()
	// DB operations
	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil).Once()
	mockTx.On("Rollback").Return(nil).Once()
	expectedRepoErr := errors.New("database get all error")
	mockPostRepo.On("GetAll", mock.Anything, mockTx, search, entities.PaginationParams{Page: page, Limit: limit}).Return(nil, expectedRepoErr).Once()

	result, err := service.GetAll(ctx, search, page, limit)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedRepoErr, err)

	mockCache.AssertExpectations(t)
	mockDB.AssertExpectations(t)
	mockPostRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
	mockCache.AssertNotCalled(t, "Set") // Cache should not be set on repo error
}

func TestPostService_GetAll_CommitError(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	mockDB := new(mocks.MockDatabase)
	mockPostRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockCache := new(mocks.MockCache)
	mockTx := new(mocks.MockTransaction)

	service := &services.PostServiceImpl{
		DB:             mockDB,
		PostRepository: mockPostRepo,
		UserRepository: mockUserRepo,
		Cache:          mockCache,
		CtxTimeout:     2 * time.Second,
	}

	search := "keyword"
	page := 1
	limit := 10
	expectedPosts := []entities.Post{
		{ID: uuid.New(), Title: "Post 1"},
	}
	expectedPostsJSON, _ := json.Marshal(expectedPosts)

	queryParams := fmt.Sprintf("search=%s:page=%d:limit=%d", search, page, limit)
	hasher := sha256.New()
	hasher.Write([]byte(queryParams))
	cacheKey := "posts:" + hex.EncodeToString(hasher.Sum(nil))

	// Simulate cache miss
	mockCache.On("Get", mock.Anything, cacheKey).Return("", errors.New("cache miss")).Once()
	// DB operations
	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil).Once()
	expectedCommitErr := errors.New("failed to commit transaction")
	mockTx.On("Commit").Return(expectedCommitErr).Once()
	mockTx.On("Rollback").Return(nil).Once()
	mockPostRepo.On("GetAll", mock.Anything, mockTx, search, entities.PaginationParams{Page: page, Limit: limit}).Return(expectedPosts, nil).Once()
	mockCache.On("Set", mock.Anything, cacheKey, expectedPostsJSON, 30*time.Second).Return(nil).Once() // Cache set should still happen before commit

	result, err := service.GetAll(ctx, search, page, limit)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedCommitErr, err)

	mockCache.AssertExpectations(t)
	mockDB.AssertExpectations(t)
	mockPostRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestPostService_GetAll_CacheSetError(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	mockDB := new(mocks.MockDatabase)
	mockPostRepo := new(mocks.MockPostRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockCache := new(mocks.MockCache)
	mockTx := new(mocks.MockTransaction)

	service := &services.PostServiceImpl{
		DB:             mockDB,
		PostRepository: mockPostRepo,
		UserRepository: mockUserRepo,
		Cache:          mockCache,
		CtxTimeout:     2 * time.Second,
	}

	search := "keyword"
	page := 1
	limit := 10
	expectedPosts := []entities.Post{
		{ID: uuid.New(), Title: "Post 1"},
	}
	expectedPostsJSON, _ := json.Marshal(expectedPosts)

	queryParams := fmt.Sprintf("search=%s:page=%d:limit=%d", search, page, limit)
	hasher := sha256.New()
	hasher.Write([]byte(queryParams))
	cacheKey := "posts:" + hex.EncodeToString(hasher.Sum(nil))

	// Simulate cache miss
	mockCache.On("Get", mock.Anything, cacheKey).Return("", errors.New("cache miss")).Once()
	// DB operations
	mockDB.On("BeginTx", mock.Anything).Return(mockTx, nil).Once()
	mockTx.On("Commit").Return(nil).Once()
	mockTx.On("Rollback").Return(nil).Maybe()
	mockPostRepo.On("GetAll", mock.Anything, mockTx, search, entities.PaginationParams{Page: page, Limit: limit}).Return(expectedPosts, nil).Once()
	// Simulate Cache Set error
	expectedCacheSetErr := errors.New("failed to set cache")
	mockCache.On("Set", mock.Anything, cacheKey, expectedPostsJSON, 30*time.Second).Return(expectedCacheSetErr).Once()

	result, err := service.GetAll(ctx, search, page, limit)

	// Assertions: The main operation should still succeed, error is just logged/warned
	assert.NoError(t, err) // Crucial: cache set error does not return an error from GetAll
	assert.NotNil(t, result)
	assert.Len(t, result, len(expectedPosts))
	assert.Equal(t, expectedPosts[0].Title, result[0].Title)

	mockCache.AssertExpectations(t)
	mockDB.AssertExpectations(t)
	mockPostRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}
