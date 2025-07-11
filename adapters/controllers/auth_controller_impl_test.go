package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
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

func TestAuthController_SignIn_Success(t *testing.T) {
	mockService := new(mocks.MockAuthService)
	controller := &controllers.AuthController{
		AuthService: mockService,
	}

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	reqBody := &dto.AuthSignInRequest{
		Email:    "user@mail.com",
		Password: "password1234",
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/signin", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	req = req.WithContext(ctx)

	expectedUserArgument := &entities.User{
		Email:    "user@mail.com",
		Password: "password1234",
	}

	// This is the user object you want to be returned by the service
	returnedUser := &entities.User{
		ID:        uuid.New(),
		Email:     "user@mail.com",
		Password:  "password1234",
		CreatedAt: time.Date(2023, time.January, 15, 10, 0, 0, 0, time.UTC),
	}
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTA0MTMxMDAsInVzZXJfaWQiOiJhMzA4MzkwMS02OTQ2LTRlZWQtYmEyMi0zN2EzZWRjYzU5NzkifQ.pOcBPTxP0ZDCp_Kv-wdpVdC0XoEVh5_Pt-mf_V6G7KY"

	mockService.On("SignIn", mock.Anything, expectedUserArgument).Return(returnedUser, token, nil)

	controller.SignIn(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response dto.WebResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err, "Failed to unmarshal response body")

	// Assert top-level WebResponse fields
	assert.Equal(t, "Successfully signed in", response.Message)
	assert.Equal(t, 1, response.Status)

	expectedAuthResponse := dto.AuthResponse{
		Token: token,
		User: dto.UserResponse{
			Id:        returnedUser.ID.String(),
			Email:     returnedUser.Email,
			CreatedAt: returnedUser.CreatedAt,
		},
	}

	expectedDataBytes, _ := json.Marshal(expectedAuthResponse)
	var expectedDataInterface interface{}
	json.Unmarshal(expectedDataBytes, &expectedDataInterface)

	assert.NotNil(t, response.Data, "Response Data should not be nil")
	assert.IsType(t, map[string]interface{}{}, response.Data, "Response Data should be of type map[string]interface{}")

	actualAuthResponseMap, ok := response.Data.(map[string]interface{})
	assert.True(t, ok, "Failed to cast response.Data to map[string]interface{}")

	// Convert the map back to AuthResponse DTO to compare fields easily
	var actualAuthResponse dto.AuthResponse
	actualAuthResponseBytes, _ := json.Marshal(actualAuthResponseMap)
	err = json.Unmarshal(actualAuthResponseBytes, &actualAuthResponse)
	assert.NoError(t, err, "Failed to unmarshal actualAuthResponseMap to AuthResponse")

	assert.Equal(t, expectedAuthResponse.Token, actualAuthResponse.Token, "Token mismatch")
	assert.Equal(t, expectedAuthResponse.User.Id, actualAuthResponse.User.Id, "User ID mismatch")
	assert.Equal(t, expectedAuthResponse.User.Email, actualAuthResponse.User.Email, "User Email mismatch")

	mockService.AssertExpectations(t)
}

func TestAuthController_SignIn_Unauthorized(t *testing.T) {
	mockService := new(mocks.MockAuthService)
	controller := &controllers.AuthController{
		AuthService: mockService,
	}

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	reqBody := &dto.AuthSignInRequest{
		Email:    "user@mail.com",
		Password: "password1234",
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/signin", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	req = req.WithContext(ctx)

	expectedUserArgument := &entities.User{
		Email:    "user@mail.com",
		Password: "password1234",
	}

	token := ""

	mockService.On("SignIn", mock.Anything, expectedUserArgument).Return(nil, token, appErrors.NewNotFoundError("Unauthorized", nil))

	controller.SignIn(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response dto.WebResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err, "Failed to unmarshal response body")

	assert.Equal(t, "Invalid email or password", response.Message)
	assert.Equal(t, 0, response.Status)
}

func TestAuthController_SignUp_Success(t *testing.T) {
	mockService := new(mocks.MockAuthService)
	controller := &controllers.AuthController{
		AuthService: mockService,
	}

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	reqBody := &dto.AuthSignUpRequest{
		Email:           "user@mail.com",
		Password:        "password1234",
		ConfirmPassword: "password1234",
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/signup", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	req = req.WithContext(ctx)

	expectedUserArgument := &entities.User{
		Email:    "user@mail.com",
		Password: "password1234",
	}

	// This is the user object you want to be returned by the service
	returnedUser := &entities.User{
		ID:        uuid.New(),
		Email:     "user@mail.com",
		Password:  "password1234",
		CreatedAt: time.Date(2023, time.January, 15, 10, 0, 0, 0, time.UTC),
	}
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTA0MTMxMDAsInVzZXJfaWQiOiJhMzA4MzkwMS02OTQ2LTRlZWQtYmEyMi0zN2EzZWRjYzU5NzkifQ.pOcBPTxP0ZDCp_Kv-wdpVdC0XoEVh5_Pt-mf_V6G7KY"

	mockService.On("SignUp", mock.Anything, expectedUserArgument).Return(returnedUser, token, nil)

	controller.SignUp(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response dto.WebResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err, "Failed to unmarshal response body")

	// Assert top-level WebResponse fields
	assert.Equal(t, "Successfully signed up", response.Message)
	assert.Equal(t, 1, response.Status)

	expectedAuthResponse := dto.AuthResponse{
		Token: token,
		User: dto.UserResponse{
			Id:        returnedUser.ID.String(),
			Email:     returnedUser.Email,
			CreatedAt: returnedUser.CreatedAt,
		},
	}

	expectedDataBytes, _ := json.Marshal(expectedAuthResponse)
	var expectedDataInterface interface{}
	json.Unmarshal(expectedDataBytes, &expectedDataInterface)

	assert.NotNil(t, response.Data, "Response Data should not be nil")
	assert.IsType(t, map[string]interface{}{}, response.Data, "Response Data should be of type map[string]interface{}")

	actualAuthResponseMap, ok := response.Data.(map[string]interface{})
	assert.True(t, ok, "Failed to cast response.Data to map[string]interface{}")

	// Convert the map back to AuthResponse DTO to compare fields easily
	var actualAuthResponse dto.AuthResponse
	actualAuthResponseBytes, _ := json.Marshal(actualAuthResponseMap)
	err = json.Unmarshal(actualAuthResponseBytes, &actualAuthResponse)
	assert.NoError(t, err, "Failed to unmarshal actualAuthResponseMap to AuthResponse")

	assert.Equal(t, expectedAuthResponse.Token, actualAuthResponse.Token, "Token mismatch")
	assert.Equal(t, expectedAuthResponse.User.Id, actualAuthResponse.User.Id, "User ID mismatch")
	assert.Equal(t, expectedAuthResponse.User.Email, actualAuthResponse.User.Email, "User Email mismatch")

	mockService.AssertExpectations(t)
}

func TestAuthController_SignUp_InvalidPassword(t *testing.T) {
	mockService := new(mocks.MockAuthService)
	controller := &controllers.AuthController{
		AuthService: mockService,
	}

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	reqBody := &dto.AuthSignUpRequest{
		Email:           "user@mail.com",
		Password:        "password1234",
		ConfirmPassword: "password1224",
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/signup", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	req = req.WithContext(ctx)

	controller.SignUp(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response dto.WebResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err, "Failed to unmarshal response body")

	assert.Equal(t, "Password and confirm password do not match", response.Message)
	assert.Equal(t, 0, response.Status)

	mockService.AssertNotCalled(t, "SignUp")
}

func TestAuthController_SignUp_InvalidEmail(t *testing.T) {
	mockService := new(mocks.MockAuthService)
	controller := &controllers.AuthController{
		AuthService: mockService,
	}

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	reqBody := &dto.AuthSignUpRequest{
		Email:           "invalidemail",
		Password:        "password1234",
		ConfirmPassword: "password1234",
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/signup", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	req = req.WithContext(ctx)

	controller.SignUp(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response dto.WebResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err, "Failed to unmarshal response body")

	assert.Equal(t, "Invalid email format for email", response.Message)
	assert.Equal(t, 0, response.Status)

	mockService.AssertNotCalled(t, "SignUp")
}

func TestAuthController_Refresh_Success(t *testing.T) {
    mockService := new(mocks.MockAuthService)
    controller := &controllers.AuthController{
        AuthService: mockService,
    }

    ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

    reqBody := &dto.RefreshTokenRequest{
        RefreshToken: "refresh-token-123",
    }
    bodyBytes, _ := json.Marshal(reqBody)
    req := httptest.NewRequest(http.MethodPost, "/api/refresh", bytes.NewReader(bodyBytes))
    req.Header.Set("Content-Type", "application/json")
    rec := httptest.NewRecorder()
    req = req.WithContext(ctx)

    accessToken := "new-access-token"
    newRefreshToken := "new-refresh-token"
    mockService.On("RefreshToken", mock.Anything, reqBody.RefreshToken).Return(accessToken, newRefreshToken, nil)

    controller.Refresh(rec, req)
    assert.Equal(t, http.StatusOK, rec.Code)

    var response dto.WebResponse
    err := json.Unmarshal(rec.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, "Successfully refreshed tokens", response.Message)
    assert.Equal(t, 1, response.Status)
    data := response.Data.(map[string]interface{})
    assert.Equal(t, accessToken, data["token"])
    assert.Equal(t, newRefreshToken, data["refresh_token"])
}

func TestAuthController_Refresh_InvalidToken(t *testing.T) {
    mockService := new(mocks.MockAuthService)
    controller := &controllers.AuthController{
        AuthService: mockService,
    }

    ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

    reqBody := &dto.RefreshTokenRequest{
        RefreshToken: "invalid-refresh-token",
    }
    bodyBytes, _ := json.Marshal(reqBody)
    req := httptest.NewRequest(http.MethodPost, "/api/refresh", bytes.NewReader(bodyBytes))
    req.Header.Set("Content-Type", "application/json")
    rec := httptest.NewRecorder()
    req = req.WithContext(ctx)

    mockService.On("RefreshToken", mock.Anything, reqBody.RefreshToken).Return("", "", appErrors.NewUnauthorizedError("Invalid or expired refresh token", nil))

    controller.Refresh(rec, req)
    assert.Equal(t, http.StatusBadRequest, rec.Code)
    var response dto.WebResponse
    err := json.Unmarshal(rec.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, "Invalid or expired refresh token", response.Message)
    assert.Equal(t, 0, response.Status)
}

func TestAuthController_Logout_Success(t *testing.T) {
    mockService := new(mocks.MockAuthService)
    controller := &controllers.AuthController{
        AuthService: mockService,
    }

    ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

    reqBody := &dto.LogoutRequest{
        RefreshToken: "refresh-token-logout",
    }
    bodyBytes, _ := json.Marshal(reqBody)
    req := httptest.NewRequest(http.MethodPost, "/api/logout", bytes.NewReader(bodyBytes))
    req.Header.Set("Content-Type", "application/json")
    rec := httptest.NewRecorder()
    req = req.WithContext(ctx)

    mockService.On("Logout", mock.Anything, reqBody.RefreshToken).Return(nil)

    controller.Logout(rec, req)
    assert.Equal(t, http.StatusOK, rec.Code)
    var response dto.WebResponse
    err := json.Unmarshal(rec.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, "Successfully logged out", response.Message)
    assert.Equal(t, 1, response.Status)
}

func TestAuthController_Logout_InvalidToken(t *testing.T) {
    mockService := new(mocks.MockAuthService)
    controller := &controllers.AuthController{
        AuthService: mockService,
    }

    ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

    reqBody := &dto.LogoutRequest{
        RefreshToken: "invalid-refresh-token",
    }
    bodyBytes, _ := json.Marshal(reqBody)
    req := httptest.NewRequest(http.MethodPost, "/api/logout", bytes.NewReader(bodyBytes))
    req.Header.Set("Content-Type", "application/json")
    rec := httptest.NewRecorder()
    req = req.WithContext(ctx)

    mockService.On("Logout", mock.Anything, reqBody.RefreshToken).Return(appErrors.NewUnauthorizedError("Invalid or expired refresh token", nil))

    controller.Logout(rec, req)
    assert.Equal(t, http.StatusBadRequest, rec.Code)
    var response dto.WebResponse
    err := json.Unmarshal(rec.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, "Invalid or expired refresh token", response.Message)
    assert.Equal(t, 0, response.Status)
}
