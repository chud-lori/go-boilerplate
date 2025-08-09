package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chud-lori/go-boilerplate/adapters/controllers"
	"github.com/chud-lori/go-boilerplate/adapters/web/dto"
	"github.com/chud-lori/go-boilerplate/domain/entities"
	"github.com/chud-lori/go-boilerplate/mocks"
	appErrors "github.com/chud-lori/go-boilerplate/pkg/errors"
	"github.com/chud-lori/go-boilerplate/pkg/logger"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPostController_Create_Success(t *testing.T) {
	mockService := new(mocks.MockPostService)
	controller := &controllers.PostController{
		PostService: mockService,
	}

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	authorID := uuid.New()
	reqBody := &dto.CreatePostRequest{
		Title:    "Test Post Title",
		Body:     "This is the body of the test post.",
		AuthorID: authorID,
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/post", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	req = req.WithContext(ctx)

	user := &entities.User{
		ID: reqBody.AuthorID,
	}

	expectedPost := &entities.Post{
		ID:        uuid.New(),
		Title:     reqBody.Title,
		Body:      reqBody.Body,
		User:      user,
		CreatedAt: time.Now(),
	}

	mockService.On("Create", mock.Anything, mock.MatchedBy(func(post *entities.Post) bool {
		return post.Title == reqBody.Title && post.Body == reqBody.Body && post.User.ID == reqBody.AuthorID
	})).Return(expectedPost, nil).Once()

	controller.Create(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var response dto.WebResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Successfully Create post", response.Message)
	assert.Equal(t, 1, response.Status)
	assert.NotNil(t, response.Data)

	postResponse, ok := response.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, expectedPost.ID.String(), postResponse["id"])
	assert.Equal(t, expectedPost.Title, postResponse["title"])
	assert.Equal(t, expectedPost.Body, postResponse["body"])
	assert.Equal(t, expectedPost.User.ID.String(), postResponse["author_id"])

	mockService.AssertExpectations(t)
}

func TestPostController_Create_InvalidPayload(t *testing.T) {
	mockService := new(mocks.MockPostService)
	controller := &controllers.PostController{
		PostService: mockService,
	}

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	// Missing required fields for validation error
	reqBody := `{"title": "","body": "","author_id": "invalid-uuid"}`
	req := httptest.NewRequest(http.MethodPost, "/post", bytes.NewReader([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	req = req.WithContext(ctx)

	controller.Create(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response dto.WebResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Contains(t, response.Message, "Invalid payload format") // Assuming GetPayload returns a validation error message
	assert.Equal(t, 0, response.Status)
	assert.Nil(t, response.Data)

	mockService.AssertNotCalled(t, "Create")
}

func TestPostController_Create_ServiceError(t *testing.T) {
	mockService := new(mocks.MockPostService)
	controller := &controllers.PostController{
		PostService: mockService,
	}

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	authorID := uuid.New()
	reqBody := &dto.CreatePostRequest{
		Title:    "Test Post Title",
		Body:     "This is the body of the test post.",
		AuthorID: authorID,
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/post", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	req = req.WithContext(ctx)

	// Simulate an AppError from the service
	mockErr := errors.New("Error")
	mockService.On("Create", mock.Anything, mock.AnythingOfType("*entities.Post")).Return(nil, mockErr).Once()

	controller.Create(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var response dto.WebResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "An unexpected error occurred", response.Message)
	assert.Equal(t, 0, response.Status)
	assert.Nil(t, response.Data)

	mockService.AssertExpectations(t)
}

func TestPostController_Update_Success(t *testing.T) {
	mockService := new(mocks.MockPostService)
	controller := &controllers.PostController{
		PostService: mockService,
	}

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	authorID := uuid.New()
	postID := uuid.New()               // Assuming the update payload includes the ID
	reqBody := &dto.CreatePostRequest{ // Assuming Update uses CreatePostRequest struct
		Title:    "Updated Post Title",
		Body:     "This is the updated body.",
		AuthorID: authorID,
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/post", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	req = req.WithContext(ctx)

	user := &entities.User{
		ID: reqBody.AuthorID,
	}

	updatedPost := &entities.Post{
		ID:        postID,
		Title:     reqBody.Title,
		Body:      reqBody.Body,
		User:      user,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockService.On("Update", mock.Anything, mock.MatchedBy(func(post *entities.Post) bool {
		// Note: Your controller's Update method does not extract ID from URL path
		// and it seems to expect the ID to be part of the 'payload' (req.Title, req.Body, req.AuthorID).
		// For a real update, the payload/request DTO would typically include the ID.
		// If ID is meant to come from the path, the controller needs to retrieve it.
		// For this test, we assume the 'payload' passed to service would contain the ID.
		// A more robust test would ensure the ID is correctly parsed from the URL and passed.
		return post.Title == reqBody.Title && post.Body == reqBody.Body && post.User.ID == reqBody.AuthorID
	})).Return(updatedPost, nil).Once()

	controller.Update(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code) // Controller returns StatusCreated for success
	var response dto.WebResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Successfully Update post", response.Message)
	assert.Equal(t, 1, response.Status)
	assert.NotNil(t, response.Data)

	postResponse, ok := response.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, updatedPost.ID.String(), postResponse["id"])
	assert.Equal(t, updatedPost.Title, postResponse["title"])

	mockService.AssertExpectations(t)
}

func TestPostController_Update_PostNotFound(t *testing.T) {
	mockService := new(mocks.MockPostService)
	controller := &controllers.PostController{
		PostService: mockService,
	}

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	authorID := uuid.New()
	reqBody := &dto.CreatePostRequest{
		Title:    "Updated Post Title",
		Body:     "This is the updated body.",
		AuthorID: authorID,
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/post", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	req = req.WithContext(ctx)

	mockError := appErrors.NewNotFoundError("Post not found", nil)
	mockService.On("Update", mock.Anything, mock.AnythingOfType("*entities.Post")).Return(nil, mockError).Once()

	controller.Update(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	var response dto.WebResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, mockError.Message, response.Message)
	assert.Equal(t, 0, response.Status)
	assert.Nil(t, response.Data)

	mockService.AssertExpectations(t)
}

func TestPostController_Delete_Success(t *testing.T) {
	mockService := new(mocks.MockPostService)
	controller := &controllers.PostController{
		PostService: mockService,
	}

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	postID := uuid.New()

	req := httptest.NewRequest(http.MethodDelete, "/post/"+postID.String(), nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	req.SetPathValue("postId", postID.String())
	req = req.WithContext(ctx)

	mockService.On("Delete", mock.Anything, postID).Return(nil).Once()

	controller.Delete(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code) // Controller returns StatusCreated for success
	var response dto.WebResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Successfully Delete post", response.Message)
	assert.Equal(t, 1, response.Status)
	assert.Nil(t, response.Data)

	mockService.AssertExpectations(t)
}

func TestPostController_Delete_InvalidUUID(t *testing.T) {
	mockService := new(mocks.MockPostService)
	controller := &controllers.PostController{
		PostService: mockService,
	}

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	invalidPostID := "not-a-uuid"

	req := httptest.NewRequest(http.MethodDelete, "/post/"+invalidPostID, nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	req.SetPathValue("postId", invalidPostID)
	req = req.WithContext(ctx)

	controller.Delete(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var response dto.WebResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Invalid postId format", response.Message)
	assert.Equal(t, 0, response.Status)
	assert.Nil(t, response.Data)

	mockService.AssertNotCalled(t, "Delete")
}

func TestPostController_Delete_PostNotFound(t *testing.T) {
	mockService := new(mocks.MockPostService)
	controller := &controllers.PostController{
		PostService: mockService,
	}

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	postID := uuid.New()

	req := httptest.NewRequest(http.MethodDelete, "/post/"+postID.String(), nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	req.SetPathValue("postId", postID.String())
	req = req.WithContext(ctx)

	mockError := appErrors.NewNotFoundError("Post not found", nil)
	mockService.On("Delete", mock.Anything, postID).Return(mockError).Once()

	controller.Delete(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	var response dto.WebResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, mockError.Message, response.Message)
	assert.Equal(t, 0, response.Status)
	assert.Nil(t, response.Data)

	mockService.AssertExpectations(t)
}

func TestPostController_GetById_Success(t *testing.T) {
	mockService := new(mocks.MockPostService)
	controller := &controllers.PostController{
		PostService: mockService,
	}

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	postID := uuid.New()

	user := &entities.User{
		ID: uuid.New(),
	}
	expectedPost := &entities.Post{
		ID:        postID,
		Title:     "Found Post",
		Body:      "This is the content of the found post.",
		User:      user,
		CreatedAt: time.Now(),
	}

	req := httptest.NewRequest(http.MethodGet, "/post/"+postID.String(), nil)
	req = req.WithContext(ctx)
	req.SetPathValue("postId", postID.String())
	rec := httptest.NewRecorder()

	mockService.On("GetById", mock.Anything, postID).Return(expectedPost, nil).Once()

	controller.GetById(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code) // Controller returns StatusCreated for success
	var response dto.WebResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Successfully Get post", response.Message)
	assert.Equal(t, 1, response.Status)
	assert.NotNil(t, response.Data)

	// Since entities.Post has mixed case fields, JSON unmarshaling might convert them.
	// It's safer to unmarshal to a map[string]interface{} and then assert.
	postData, ok := response.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, expectedPost.ID.String(), postData["id"])
	assert.Equal(t, expectedPost.Title, postData["title"])
	assert.Equal(t, expectedPost.Body, postData["body"])

	mockService.AssertExpectations(t)
}

func TestPostController_GetById_InvalidUUID(t *testing.T) {
	mockService := new(mocks.MockPostService)
	controller := &controllers.PostController{
		PostService: mockService,
	}

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	invalidPostID := "not-a-uuid"

	req := httptest.NewRequest(http.MethodGet, "/post/"+invalidPostID, nil)
	req = req.WithContext(ctx)
	req.SetPathValue("postId", invalidPostID)
	rec := httptest.NewRecorder()

	controller.GetById(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var response dto.WebResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Invalid postId format", response.Message)
	assert.Equal(t, 0, response.Status)
	assert.Nil(t, response.Data)

	mockService.AssertNotCalled(t, "GetById")
}

func TestPostController_GetById_PostNotFound(t *testing.T) {
	mockService := new(mocks.MockPostService)
	controller := &controllers.PostController{
		PostService: mockService,
	}

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	postID := uuid.New()

	req := httptest.NewRequest(http.MethodGet, "/post/"+postID.String(), nil)
	req = req.WithContext(ctx)
	req.SetPathValue("postId", postID.String())
	rec := httptest.NewRecorder()

	mockErr := appErrors.NewNotFoundError("Post not found", nil)
	mockService.On("GetById", mock.Anything, postID).Return(nil, mockErr).Once()

	controller.GetById(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	var response dto.WebResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, mockErr.Message, response.Message)
	assert.Equal(t, 0, response.Status)
	assert.Nil(t, response.Data)

	mockService.AssertExpectations(t)
}

func TestPostController_GetAll_Success(t *testing.T) {
	mockService := new(mocks.MockPostService)
	controller := &controllers.PostController{
		PostService: mockService,
	}

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	user1 := &entities.User{
		ID: uuid.New(),
	}
	user2 := &entities.User{
		ID: uuid.New(),
	}
	expectedPosts := []entities.Post{
		{ID: uuid.New(), Title: "Post 1", Body: "Body 1", User: user1, CreatedAt: time.Now()},
		{ID: uuid.New(), Title: "Post 2", Body: "Body 2", User: user2, CreatedAt: time.Now()},
	}

	req := httptest.NewRequest(http.MethodGet, "/post?search=test&page=2&limit=5", nil)
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	mockService.On("GetAll", mock.Anything, "test", 2, 5).Return(expectedPosts, nil).Once()

	controller.GetAll(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code) // Controller returns StatusCreated for success
	var response dto.WebResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Successfully Get posts", response.Message)
	assert.Equal(t, 1, response.Status)
	assert.NotNil(t, response.Data)

	// Unmarshal Data to []dto.PostResponse or []map[string]interface{} for assertion
	postsData, ok := response.Data.([]interface{})
	assert.True(t, ok)
	assert.Len(t, postsData, len(expectedPosts))

	// You might want to do more granular assertions on individual post fields
	// For example:
	post1 := postsData[0].(map[string]interface{})
	assert.Equal(t, expectedPosts[0].ID.String(), post1["id"])
	assert.Equal(t, expectedPosts[0].Title, post1["title"])

	mockService.AssertExpectations(t)
}

func TestPostController_GetAll_DefaultQueryParams(t *testing.T) {
	mockService := new(mocks.MockPostService)
	controller := &controllers.PostController{
		PostService: mockService,
	}

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	user := &entities.User{
		ID: uuid.New(),
	}
	expectedPosts := []entities.Post{
		{ID: uuid.New(), Title: "Post A", Body: "Body A", User: user, CreatedAt: time.Now()},
	}

	req := httptest.NewRequest(http.MethodGet, "/post", nil) // No query params
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	mockService.On("GetAll", mock.Anything, "", 1, 10).Return(expectedPosts, nil).Once() // Expect defaults

	controller.GetAll(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var response dto.WebResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Successfully Get posts", response.Message)
	assert.Equal(t, 1, response.Status)
	assert.NotNil(t, response.Data)

	mockService.AssertExpectations(t)
}

func TestPostController_GetAll_ServiceError(t *testing.T) {
	mockService := new(mocks.MockPostService)
	controller := &controllers.PostController{
		PostService: mockService,
	}

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	req := httptest.NewRequest(http.MethodGet, "/post", nil)
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	mockErr := errors.New("Error")
	mockService.On("GetAll", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, mockErr).Once()

	controller.GetAll(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	var response dto.WebResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "An unexpected error occurred", response.Message)
	assert.Equal(t, 0, response.Status)
	assert.Nil(t, response.Data)

	mockService.AssertExpectations(t)
}
func TestPostController_UploadAttachment_Success(t *testing.T) {
	mockService := new(mocks.MockPostService)
	controller := &controllers.PostController{
		PostService: mockService,
	}

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	// Prepare multipart form data
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, _ := w.CreateFormFile("file", "file.txt")
	fw.Write([]byte("filedata"))
	w.WriteField("file_name", "file.txt")
	w.WriteField("file_type", "text/plain")
	w.Close()

	postID := uuid.New()
	url := "/post/" + postID.String() + "/upload"
	req := httptest.NewRequest(http.MethodPost, url, &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	mockService.On("StartAsyncUpload", mock.Anything, postID, "file.txt", "text/plain", mock.AnythingOfType("[]uint8")).Return(uuid.New(), nil)

	controller.UploadAttachment(rec, req)
	assert.Equal(t, http.StatusAccepted, rec.Code)
	var response dto.WebResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Upload started", response.Message)
	assert.Equal(t, 1, response.Status)
	assert.NotNil(t, response.Data)
	mockService.AssertExpectations(t)
}

func TestPostController_UploadAttachment_InvalidPostID(t *testing.T) {
	mockService := new(mocks.MockPostService)
	controller := &controllers.PostController{
		PostService: mockService,
	}

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, _ := w.CreateFormFile("file", "file.txt")
	fw.Write([]byte("filedata"))
	w.WriteField("file_name", "file.txt")
	w.WriteField("file_type", "text/plain")
	w.Close()

	url := "/post/not-a-uuid/upload"
	req := httptest.NewRequest(http.MethodPost, url, &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	controller.UploadAttachment(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var response dto.WebResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid post ID", response.Message)
	assert.Equal(t, 0, response.Status)
	assert.Nil(t, response.Data)
}

func TestPostController_UploadAttachment_MissingFile(t *testing.T) {
	mockService := new(mocks.MockPostService)
	controller := &controllers.PostController{
		PostService: mockService,
	}

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	postID := uuid.New()
	url := "/post/" + postID.String() + "/upload"
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	// No file added
	w.WriteField("file_name", "file.txt")
	w.WriteField("file_type", "text/plain")
	w.Close()

	req := httptest.NewRequest(http.MethodPost, url, &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	controller.UploadAttachment(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var response dto.WebResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "File is required", response.Message)
	assert.Equal(t, 0, response.Status)
	assert.Nil(t, response.Data)
}

func TestPostController_UploadAttachment_ServiceError(t *testing.T) {
	mockService := new(mocks.MockPostService)
	controller := &controllers.PostController{
		PostService: mockService,
	}

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, _ := w.CreateFormFile("file", "file.txt")
	fw.Write([]byte("filedata"))
	w.WriteField("file_name", "file.txt")
	w.WriteField("file_type", "text/plain")
	w.Close()

	postID := uuid.New()
	url := "/post/" + postID.String() + "/upload"
	req := httptest.NewRequest(http.MethodPost, url, &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	mockService.On("StartAsyncUpload", mock.Anything, postID, "file.txt", "text/plain", mock.AnythingOfType("[]uint8")).Return(uuid.Nil, errors.New("service error"))

	controller.UploadAttachment(rec, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	var response dto.WebResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Failed to start async upload", response.Message)
	assert.Equal(t, 0, response.Status)
	assert.Nil(t, response.Data)
	mockService.AssertExpectations(t)
}

func TestPostController_UploadStatusSSE_Success(t *testing.T) {
	mockService := new(mocks.MockPostService)
	controller := &controllers.PostController{
		PostService: mockService,
	}

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	uploadID := uuid.New()

	req := httptest.NewRequest(http.MethodGet, "/uploads/"+uploadID.String()+"/events", nil)
	req = req.WithContext(ctx)
	req.SetPathValue("uploadId", uploadID.String())
	rec := httptest.NewRecorder()

	// Simulate status progression: uploading -> success
	mockService.On("GetUploadStatus", mock.Anything, uploadID).Return(entities.UploadStatusUploading, nil).Once()
	mockService.On("GetUploadStatus", mock.Anything, uploadID).Return(entities.UploadStatusSuccess, nil).Once()

	controller.UploadStatusSSE(rec, req)

	// Check SSE output contains both statuses
	body := rec.Body.String()
	assert.Contains(t, body, "data: uploading")
	assert.Contains(t, body, "data: success")

	mockService.AssertExpectations(t)
}

func TestPostController_UploadStatusSSE_InvalidUUID(t *testing.T) {
	mockService := new(mocks.MockPostService)
	controller := &controllers.PostController{
		PostService: mockService,
	}

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	invalidUploadID := "not-a-uuid"

	req := httptest.NewRequest(http.MethodGet, "/uploads/"+invalidUploadID+"/events", nil)
	req = req.WithContext(ctx)
	req.SetPathValue("uploadId", invalidUploadID)
	rec := httptest.NewRecorder()

	controller.UploadStatusSSE(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var response dto.WebResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "invalid uploadId format", response.Message)
	assert.Equal(t, 0, response.Status)
	assert.Nil(t, response.Data)
}

func TestPostController_UploadStatusSSE_MissingUploadID(t *testing.T) {
	mockService := new(mocks.MockPostService)
	controller := &controllers.PostController{
		PostService: mockService,
	}

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	req := httptest.NewRequest(http.MethodGet, "/uploads//events", nil)
	req = req.WithContext(ctx)
	req.SetPathValue("uploadId", "")
	rec := httptest.NewRecorder()

	controller.UploadStatusSSE(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var response dto.WebResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "uploadId required", response.Message)
	assert.Equal(t, 0, response.Status)
	assert.Nil(t, response.Data)
}

func TestPostController_UploadStatusSSE_ServiceError(t *testing.T) {
	mockService := new(mocks.MockPostService)
	controller := &controllers.PostController{
		PostService: mockService,
	}

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))
	uploadID := uuid.New()

	req := httptest.NewRequest(http.MethodGet, "/uploads/"+uploadID.String()+"/events", nil)
	req = req.WithContext(ctx)
	req.SetPathValue("uploadId", uploadID.String())
	rec := httptest.NewRecorder()

	mockService.On("GetUploadStatus", mock.Anything, uploadID).Return(entities.UploadStatusFailed, errors.New("service error")).Once()

	controller.UploadStatusSSE(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	var response dto.WebResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Failed to get upload status", response.Message)
	assert.Equal(t, 0, response.Status)
	assert.Nil(t, response.Data)
}
