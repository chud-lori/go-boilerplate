package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/sirupsen/logrus"
)

func TestUserController_Create_Success(t *testing.T) {
	mockService := new(mocks.MockUserService)
	controller := &controllers.UserController{
		UserService: mockService,
	}

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	reqBody := &dto.UserRequest{
		Email: "user@mail.com",
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/user", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	req = req.WithContext(ctx)

	user := &entities.User{
		Id:    "",
		Email: "user@mail.com",
	}

	mockService.On("Save", mock.Anything, user).Return(user, nil)

	controller.Create(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response dto.WebResponse
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, "success save user", response.Message)
	assert.Equal(t, 1, response.Status)

	mockService.AssertExpectations(t)
}

func TestUserController_Create_Failed(t *testing.T) {
	mockService := new(mocks.MockUserService)
	controller := &controllers.UserController{
		UserService: mockService,
	}

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	reqBody := &dto.UserRequest{
		Email: "user@mail.com",
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/user", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	req = req.WithContext(ctx)

	user := &entities.User{
		Id:    "",
		Email: "user@mail.com",
	}

	expectedError := errors.New("something went wrong during save")
	mockService.On("Save", mock.Anything, user).Return(nil, expectedError)

	controller.Create(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response dto.WebResponse
	json.Unmarshal(rec.Body.Bytes(), &response)

	assert.Equal(t, "Failed to create user", response.Message)
	assert.Equal(t, 0, response.Status)
	assert.Nil(t, response.Data)

	mockService.AssertExpectations(t)
}

func TestUserController_Update_Success(t *testing.T) {
	mockService := new(mocks.MockUserService)
	controller := &controllers.UserController{
		UserService: mockService,
	}

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	reqBody := &dto.UserRequest{
		Email:    "user@mail.com",
		Passcode: "IJHD9782",
	}
	userId := "a234f98c-3239-4c34-8ad8-f63e41bb20c8"

	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/api/user/"+userId, bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	req.SetPathValue("userId", userId)
	req = req.WithContext(ctx)

	user := &entities.User{
		Id:       userId,
		Email:    "user@mail.com",
		Passcode: "IJHD9782",
	}

	mockService.On("Update", mock.Anything, user).Return(user, nil)

	controller.Update(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response dto.WebResponse
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, "success update user", response.Message)
	assert.Equal(t, 1, response.Status)

	mockService.AssertExpectations(t)
}

func TestUserController_Update_UserNotFound(t *testing.T) {
	mockService := new(mocks.MockUserService)
	controller := &controllers.UserController{
		UserService: mockService,
	}

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	reqBody := &dto.UserRequest{
		Email:    "user@mail.com",
		Passcode: "IJHD9782",
	}
	userId := "a234f98c-3239-4c34-8ad8-f63e41bb20c8"

	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/api/user/"+userId, bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	req.SetPathValue("userId", userId)
	req = req.WithContext(ctx)

	user := &entities.User{
		Id:       userId,
		Email:    "user@mail.com",
		Passcode: "IJHD9782",
	}

	mockError := appErrors.NewNotFoundError("User not found", nil)
	mockService.On("Update", mock.Anything, user).Return(nil, mockError)

	controller.Update(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response dto.WebResponse
	json.Unmarshal(rec.Body.Bytes(), &response)

	assert.Equal(t, "User not found", response.Message)
	assert.Equal(t, 0, response.Status)
	assert.Nil(t, response.Data)

	mockService.AssertExpectations(t)
}

func TestUserController_Update_Failed(t *testing.T) {
	mockService := new(mocks.MockUserService)
	controller := &controllers.UserController{
		UserService: mockService,
	}

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	reqBody := &dto.UserRequest{
		Email:    "user@mail.com",
		Passcode: "IJHD9782",
	}
	userId := "a234f98c-3239-4c34-8ad8-f63e41bb20c8"

	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/api/user/"+userId, bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	req.SetPathValue("userId", userId)
	req = req.WithContext(ctx)

	user := &entities.User{
		Id:       userId,
		Email:    "user@mail.com",
		Passcode: "IJHD9782",
	}

	mockService.On("Update", mock.Anything, user).Return(nil, errors.New("Errors"))

	controller.Update(rec, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var response dto.WebResponse
	json.Unmarshal(rec.Body.Bytes(), &response)

	assert.Equal(t, "An unexpected error occurred", response.Message)
	assert.Equal(t, 0, response.Status)
	assert.Nil(t, response.Data)

	mockService.AssertExpectations(t)
}

func TestUserController_Delete_Success(t *testing.T) {
	mockService := new(mocks.MockUserService)
	controller := &controllers.UserController{
		UserService: mockService,
	}

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	userId := "a234f98c-3239-4c34-8ad8-f63e41bb20c8"

	req := httptest.NewRequest(http.MethodDelete, "/api/user/"+userId, nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	req.SetPathValue("userId", userId)
	req = req.WithContext(ctx)

	mockService.On("Delete", mock.Anything, userId).Return(nil)

	controller.Delete(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response dto.WebResponse
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, "success delete user", response.Message)
	assert.Equal(t, 1, response.Status)

	mockService.AssertExpectations(t)
}

func TestUserController_Delete_UserNotFound(t *testing.T) {
	mockService := new(mocks.MockUserService)
	controller := &controllers.UserController{
		UserService: mockService,
	}

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	userId := "a234f98c-3239-4c34-8ad8-f63e41bb20c8"

	req := httptest.NewRequest(http.MethodDelete, "/api/user/"+userId, nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	req.SetPathValue("userId", userId)
	req = req.WithContext(ctx)

	mockError := appErrors.NewBadRequestError("User not found", nil)
	mockService.On("Delete", mock.Anything, userId).Return(mockError)

	controller.Delete(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response dto.WebResponse
	json.Unmarshal(rec.Body.Bytes(), &response)

	assert.Equal(t, "User not found", response.Message)
	assert.Equal(t, 0, response.Status)
	assert.Nil(t, response.Data)

	mockService.AssertExpectations(t)
}

func TestUserController_Delete_Failed(t *testing.T) {
	mockService := new(mocks.MockUserService)
	controller := &controllers.UserController{
		UserService: mockService,
	}

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	userId := "a234f98c-3239-4c34-8ad8-f63e41bb20c8"

	req := httptest.NewRequest(http.MethodDelete, "/api/user/"+userId, nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	req.SetPathValue("userId", userId)
	req = req.WithContext(ctx)

	mockService.On("Delete", mock.Anything, userId).Return(errors.New("Errors"))

	controller.Delete(rec, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var response dto.WebResponse
	json.Unmarshal(rec.Body.Bytes(), &response)

	assert.Equal(t, "An unexpected error occurred", response.Message)
	assert.Equal(t, 0, response.Status)
	assert.Nil(t, response.Data)

	mockService.AssertExpectations(t)
}

func TestUserController_FindById_Success(t *testing.T) {
	mockService := new(mocks.MockUserService)
	controller := &controllers.UserController{
		UserService: mockService,
	}

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	userId := "a234f98c-3239-4c34-8ad8-f63e41bb20c8"

	req := httptest.NewRequest(http.MethodGet, "/api/user/"+userId, nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	req.SetPathValue("userId", userId)
	req = req.WithContext(ctx)

	user := &entities.User{
		Id:        userId,
		Email:     "user@mail.com",
		CreatedAt: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

	mockService.On("FindById", mock.Anything, userId).Return(user, nil)

	controller.FindById(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response dto.WebResponse
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, "success get user by id", response.Message)
	assert.Equal(t, 1, response.Status)

	assert.NotNil(t, response.Data, "Expected Data field to not be nil")

	// 2. Type assert Data to map[string]interface{}
	dataMap, ok := response.Data.(map[string]interface{})
	assert.True(t, ok, "Expected Data to be a map[string]interface{}")

	// 3. Extract the "id" field from the map
	actualId, ok := dataMap["id"].(string)
	assert.True(t, ok, "Expected 'id' field in Data to be a string")

	// 4. Assert the extracted ID
	assert.Equal(t, userId, actualId, "Expected user ID to match")

	mockService.AssertExpectations(t)
}

func TestUserController_FindById_UserNotFound(t *testing.T) {
	mockService := new(mocks.MockUserService)
	controller := &controllers.UserController{
		UserService: mockService,
	}

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	userId := "a234f98c-3239-4c34-8ad8-f63e41bb20c8"

	req := httptest.NewRequest(http.MethodGet, "/api/user/"+userId, nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	req.SetPathValue("userId", userId)
	req = req.WithContext(ctx)

	mockError := appErrors.NewBadRequestError("User not found", nil)
	mockService.On("FindById", mock.Anything, userId).Return(nil, mockError)

	controller.FindById(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response dto.WebResponse
	json.Unmarshal(rec.Body.Bytes(), &response)

	assert.Equal(t, "User not found", response.Message)
	assert.Equal(t, 0, response.Status)
	assert.Nil(t, response.Data)

	mockService.AssertExpectations(t)
}

func TestUserController_FindById_Failed(t *testing.T) {
	mockService := new(mocks.MockUserService)
	controller := &controllers.UserController{
		UserService: mockService,
	}

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	userId := "a234f98c-3239-4c34-8ad8-f63e41bb20c8"

	req := httptest.NewRequest(http.MethodGet, "/api/user/"+userId, nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	req.SetPathValue("userId", userId)
	req = req.WithContext(ctx)

	mockService.On("FindById", mock.Anything, userId).Return(nil, errors.New("Errors"))

	controller.FindById(rec, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var response dto.WebResponse
	json.Unmarshal(rec.Body.Bytes(), &response)

	assert.Equal(t, "An unexpected error occurred", response.Message)
	assert.Equal(t, 0, response.Status)
	assert.Nil(t, response.Data)

	mockService.AssertExpectations(t)
}

func TestUserController_FindAll_Success(t *testing.T) {
	mockService := new(mocks.MockUserService)
	controller := &controllers.UserController{
		UserService: mockService,
	}

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	req := httptest.NewRequest(http.MethodGet, "/api/user", nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	req = req.WithContext(ctx)

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

	mockService.On("FindAll", mock.Anything).Return(listUsers, nil)

	controller.FindAll(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var webResponse dto.WebResponse
	json.Unmarshal(rec.Body.Bytes(), &webResponse)

	// Adjust the success message for FindAll
	assert.Equal(t, "success get all users", webResponse.Message)
	assert.Equal(t, 1, webResponse.Status)

	assert.NotNil(t, webResponse.Data, "Expected Data field to not be nil")

	// --- Asserting the array content ---
	// 1. Type assert Data to []interface{} because it's an array of JSON objects
	dataArray, ok := webResponse.Data.([]interface{})
	assert.True(t, ok, "Expected Data to be a slice of interfaces")
	assert.Len(t, dataArray, len(listUsers), "Expected number of users to match")

	// 2. Iterate through the array and assert each element (or just the first one if preferred)
	// For this example, let's assert the first user's ID and Email
	if assert.NotEmpty(t, dataArray, "Expected data array not to be empty") {
		firstUserMap, ok := dataArray[0].(map[string]interface{})
		assert.True(t, ok, "Expected first element in Data to be a map[string]interface{}")

		// Assert ID of the first user
		actualId, ok := firstUserMap["id"].(string)
		assert.True(t, ok, "Expected 'id' field in first user data to be a string")
		assert.Equal(t, listUsers[0].Id, actualId, "Expected first user ID to match")

		// Assert Email of the first user
		actualEmail, ok := firstUserMap["email"].(string)
		assert.True(t, ok, "Expected 'email' field in first user data to be a string")
		assert.Equal(t, listUsers[0].Email, actualEmail, "Expected first user email to match")

		// If you want to assert the Created_at, remember it's a string in JSON
		actualCreatedAt, ok := firstUserMap["created_at"].(string)
		assert.True(t, ok, "Expected 'created_at' field in first user data to be a string")
		// You might need to parse it back to time.Time or compare string formats
		expectedCreatedAtStr := listUsers[0].CreatedAt.Format(time.RFC3339Nano)
		assert.Equal(t, expectedCreatedAtStr, actualCreatedAt, "Expected first user created_at to match")
	}

	mockService.AssertExpectations(t)
}

func TestUserController_FindAll_Failed(t *testing.T) {
	mockService := new(mocks.MockUserService)
	controller := &controllers.UserController{
		UserService: mockService,
	}

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	req := httptest.NewRequest(http.MethodGet, "/api/user", nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	req = req.WithContext(ctx)

	mockService.On("FindAll", mock.Anything).Return(nil, errors.New("Errors"))

	controller.FindAll(rec, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var response dto.WebResponse
	json.Unmarshal(rec.Body.Bytes(), &response)

	assert.Equal(t, "An unexpected error occurred", response.Message)
	assert.Equal(t, 0, response.Status)
	assert.Nil(t, response.Data)

	mockService.AssertExpectations(t)
}
