package controllers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/chud-lori/go-boilerplate/adapters/web/dto"
	"github.com/chud-lori/go-boilerplate/adapters/web/helper"
	"github.com/chud-lori/go-boilerplate/domain/entities"
	"github.com/chud-lori/go-boilerplate/domain/ports"
	appErrors "github.com/chud-lori/go-boilerplate/pkg/errors"
	"github.com/chud-lori/go-boilerplate/pkg/logger"
	"github.com/sirupsen/logrus"
)

type AuthController struct {
	ports.AuthService
}

// SignIn godoc
// @Summary Sign in an existing user
// @Description Authenticate a user and return JWT token with user info
// @ID auth-signin
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body dto.AuthSignInRequest true "SignIn request"
// @Success 200 {object} dto.WebResponse{data=dto.AuthResponse} "Successfully signed in user"
// @Failure 400 {object} dto.WebResponse "Invalid request payload or failed to sign in"
// @Router /signin [post]
// @Security ApiKeyAuth
func (c *AuthController) SignIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := ctx.Value(logger.LoggerContextKey).(*logrus.Entry)

	var req dto.AuthSignInRequest
	if err := helper.GetPayload(r, &req); err != nil {
		logger.Error("Failed to get payload:", err)
		helper.WriteResponse(w, dto.WebResponse{
			Message: "Invalid request payload",
			Status:  0,
			Data:    nil,
		}, http.StatusBadRequest)
		return
	}

	payload := &entities.User{
		Email:    req.Email,
		Password: req.Password,
	}

	user, token, err := c.AuthService.SignIn(ctx, payload)
	if err != nil {
		logger.Error("Failed to authenticate user:", err)
		helper.WriteResponse(w, dto.WebResponse{
			Message: "Invalid email or password",
			Status:  0,
			Data:    nil,
		}, http.StatusBadRequest)
		return
	}

	resp := dto.WebResponse{
		Message: "Successfully signed in",
		Status:  1,
		Data: dto.AuthResponse{
			Token: token,
			User: dto.UserResponse{
				Id:        user.ID.String(),
				Email:     user.Email,
				CreatedAt: user.CreatedAt,
			},
		},
	}
	helper.WriteResponse(w, &resp, http.StatusOK)
}

// SignUp godoc
// @Summary Sign up a new user
// @Description Register a new user with email and password, returning token and user info
// @ID auth-signup
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body dto.AuthSignUpRequest true "SignUp request"
// @Success 201 {object} dto.WebResponse{data=dto.AuthResponse} "Successfully signed up user"
// @Failure 400 {object} dto.WebResponse "Invalid request payload or failed to sign up"
// @Router /signup [post]
// @Security ApiKeyAuth
func (c *AuthController) SignUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := ctx.Value(logger.LoggerContextKey).(*logrus.Entry)

	// if err := helper.GetPayload(r, &req); err != nil {
	// 	logger.Error("Failed to get payload:", err)
	// 	m := fmt.Sprintf("Invalid request payload Errrr: %s", err.Error())
	// 	helper.WriteResponse(w, dto.WebResponse{
	// 		Message: m,
	// 		Status:  0,
	// 		Data:    nil,
	// 	}, http.StatusBadRequest)
	// 	return
	// }

	var req dto.AuthSignUpRequest
	err := helper.GetPayload(r, &req)
	if err != nil {
		// Handle specific validation errors
		var validationErr *appErrors.ValidationErrors
		if errors.As(err, &validationErr) {
			logger.Warn("Signup validation failed:", validationErr.Error()) // Log the internal error
			helper.WriteResponse(w, dto.WebResponse{
				Message: strings.Join(validationErr.Messages, ", "), // Join all validation messages
				Status:  0,
				Data:    nil,
			}, http.StatusBadRequest) // Use 400 Bad Request for validation errors
			return
		}

		// Handle other types of errors from GetPayload (e.g., bad request parsing)
		var badRequestErr *appErrors.AppError
		if errors.As(err, &badRequestErr) && badRequestErr.StatusCode == http.StatusBadRequest {
			logger.Error("Failed to get payload due to bad request:", err)
			helper.WriteResponse(w, dto.WebResponse{
				Message: badRequestErr.Message,
				Status:  0,
				Data:    nil,
			}, http.StatusBadRequest) // Explicitly use 400
			return
		}

		// Generic error for GetPayload failures not covered above
		logger.Error("Failed to get payload with unexpected error:", err)
		helper.WriteResponse(w, dto.WebResponse{
			Message: "Failed to process request payload",
			Status:  0,
			Data:    nil,
		}, http.StatusBadRequest) // Default to 400 if it's a payload issue
		return
	}

	payload := &entities.User{
		Email:    req.Email,
		Password: req.Password,
	}

	user, token, err := c.AuthService.SignUp(ctx, payload)
	if err != nil {
		var appErr *appErrors.AppError
		if errors.As(err, &appErr) {
			helper.WriteResponse(w, dto.WebResponse{
				Message: appErr.Message,
				Status:  0,
				Data:    nil,
			}, int64(appErr.StatusCode))
			return
		} else {
			helper.WriteResponse(w, dto.WebResponse{
				Message: "An unexpected error occurred",
				Status:  0,
				Data:    nil,
			}, http.StatusInternalServerError)
			return
		}
	}

	resp := dto.WebResponse{
		Message: "Successfully signed up",
		Status:  1,
		Data: dto.AuthResponse{
			Token: token,
			User: dto.UserResponse{
				Id:        user.ID.String(),
				Email:     user.Email,
				CreatedAt: user.CreatedAt,
			},
		},
	}
	helper.WriteResponse(w, &resp, http.StatusCreated)
}
