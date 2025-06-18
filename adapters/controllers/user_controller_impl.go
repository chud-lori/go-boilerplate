package controllers

import (
	"errors"
	"net/http"

	"github.com/chud-lori/go-boilerplate/adapters/web/dto"
	"github.com/chud-lori/go-boilerplate/adapters/web/helper"
	"github.com/chud-lori/go-boilerplate/domain/entities"
	"github.com/chud-lori/go-boilerplate/domain/ports"
	appErrors "github.com/chud-lori/go-boilerplate/pkg/errors"
	"github.com/chud-lori/go-boilerplate/pkg/logger"
	"github.com/google/uuid"

	"github.com/sirupsen/logrus"
)

type UserController struct {
	ports.UserService
}

// Create godoc
// @Summary Create a new user
// @Description Creates a new user with the provided email and returns the created user's details.
// @ID create-user
// @Tags Users
// @Accept json
// @Produce json
// @Param user body dto.UserRequest true "User creation request"
// @Success 201 {object} dto.WebResponse{data=dto.UserResponse} "Successfully created user"
// @Failure 400 {object} dto.WebResponse "Invalid request payload or failed to create user"
// @Router /user [post]
// @Security ApiKeyAuth
func (controller *UserController) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := ctx.Value(logger.LoggerContextKey).(*logrus.Entry)

	userRequest := dto.UserRequest{}

	if err := helper.GetPayload(r, &userRequest); err != nil {
		logger.Error("Failed to get Payload: ", err)
		helper.WriteResponse(w, dto.WebResponse{
			Message: "Invalid request payload",
			Status:  0,
			Data:    nil,
		}, http.StatusBadRequest)
		return
	}

	userPayload := &entities.User{
		Id:    "",
		Email: userRequest.Email,
	}

	result, err := controller.UserService.Save(ctx, userPayload)

	if err != nil {
		logger.Error("Failed to create user: ", err)
		helper.WriteResponse(w, dto.WebResponse{
			Message: "Failed to create user",
			Status:  0,
			Data:    nil,
		}, http.StatusBadRequest)
		return
	}

	response := dto.WebResponse{
		Message: "success save user",
		Status:  1,
		Data: dto.UserResponse{
			Id:        result.Id,
			Email:     result.Email,
			CreatedAt: result.CreatedAt,
		},
	}
	helper.WriteResponse(w, &response, http.StatusCreated)
}

// Update godoc
// @Summary Update an existing user
// @Description Updates an existing user's details by their ID.
// @ID update-user
// @Tags Users
// @Accept json
// @Produce json
// @Param userId path string true "User ID (UUID)"
// @Param user body dto.UserRequest true "User update request"
// @Success 200 {object} dto.WebResponse{data=dto.UserResponse} "Successfully updated user"
// @Failure 400 {object} dto.WebResponse "Invalid request payload, invalid user ID, or user not found"
// @Failure 500 {object} dto.WebResponse "Internal server error"
// @Router /user/{userId} [put]
// @Security ApiKeyAuth
func (controller *UserController) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := ctx.Value(logger.LoggerContextKey).(*logrus.Entry)
	userRequest := &dto.UserRequest{}
	userId := r.PathValue("userId")

	if err := helper.GetPayload(r, &userRequest); err != nil {
		logger.Error("Failed to get Payload: ", err)
		helper.WriteResponse(w, dto.WebResponse{
			Message: "Invalid request payload",
			Status:  0,
			Data:    nil,
		}, http.StatusBadRequest)
		return
	}

	userPayload := &entities.User{
		Id:       userId,
		Email:    userRequest.Email,
		Passcode: userRequest.Passcode,
	}

	userResponse, err := controller.UserService.Update(ctx, userPayload)

	if err != nil {
		var appErr *appErrors.AppError
		if errors.As(err, &appErr) {
			helper.WriteResponse(w, dto.WebResponse{
				Message: appErr.Message,
				Status:  0,
				Data:    nil,
			}, http.StatusBadRequest)
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

	response := dto.WebResponse{
		Message: "success update user",
		Status:  1,
		Data:    userResponse,
	}

	helper.WriteResponse(w, &response, http.StatusOK)
}

// Delete godoc
// @Summary Delete a user
// @Description Deletes a user by their ID.
// @ID delete-user
// @Tags Users
// @Produce json
// @Param userId path string true "User ID (UUID)"
// @Success 200 {object} dto.WebResponse "Successfully deleted user"
// @Failure 400 {object} dto.WebResponse "Invalid user ID or user not found"
// @Failure 500 {object} dto.WebResponse "Internal server error"
// @Router /user/{userId} [delete]
// @Security ApiKeyAuth
func (controller *UserController) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := r.PathValue("userId")

	err := controller.UserService.Delete(ctx, userId)

	if err != nil {
		var appErr *appErrors.AppError
		if errors.As(err, &appErr) {
			helper.WriteResponse(w, dto.WebResponse{
				Message: appErr.Message,
				Status:  0,
				Data:    nil,
			}, http.StatusBadRequest)
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

	response := dto.WebResponse{
		Message: "success delete user",
		Status:  1,
		Data:    nil,
	}
	helper.WriteResponse(w, &response, http.StatusOK)
}

// FindById godoc
// @Summary Get user by ID
// @Description Retrieves a single user's details by their ID.
// @ID get-user-by-id
// @Tags Users
// @Produce json
// @Param userId path string true "User ID (UUID)"
// @Success 200 {object} dto.WebResponse{data=dto.UserResponse} "Successfully retrieved user"
// @Failure 400 {object} dto.WebResponse "Invalid user ID format or user not found"
// @Failure 404 {object} dto.WebResponse "User not found (due to invalid UUID format)"
// @Failure 500 {object} dto.WebResponse "Internal server error"
// @Router /user/{userId} [get]
// @Security ApiKeyAuth
func (c *UserController) FindById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := ctx.Value(logger.LoggerContextKey).(*logrus.Entry)
	userId := r.PathValue("userId")

	if _, err := uuid.Parse(userId); err != nil {
		logger.Error("Invalid UUID format: ", err)
		helper.WriteResponse(w, dto.WebResponse{
			Message: "UUID format is invalid",
			Status:  0,
			Data:    nil,
		}, http.StatusNotFound)
		return
	}

	user, err := c.UserService.FindById(ctx, userId)
	if err != nil {
		var appErr *appErrors.AppError
		if errors.As(err, &appErr) {
			helper.WriteResponse(w, dto.WebResponse{
				Message: appErr.Message,
				Status:  0,
				Data:    nil,
			}, http.StatusBadRequest)
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

	response := dto.WebResponse{
		Message: "success get user by id",
		Status:  1,
		Data:    user,
	}
	helper.WriteResponse(w, &response, http.StatusOK)
}

// FindAll godoc
// @Summary Get all users
// @Description Retrieves a list of all users.
// @ID get-all-users
// @Tags Users
// @Produce json
// @Success 200 {object} dto.WebResponse{data=[]dto.UserResponse} "Successfully retrieved all users"
// @Failure 500 {object} dto.WebResponse "Internal server error"
// @Router /user [get]
// @Security ApiKeyAuth
func (controller *UserController) FindAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger, _ := ctx.Value(logger.LoggerContextKey).(*logrus.Entry)

	users, err := controller.UserService.FindAll(ctx)

	if err != nil {
		logger.Error("Error Find All user: ", err)
		helper.WriteResponse(w, dto.WebResponse{
			Message: "An unexpected error occurred",
			Status:  0,
			Data:    nil,
		}, http.StatusInternalServerError)
		return
	}

	response := dto.WebResponse{
		Message: "success get all users",
		Status:  1,
		Data:    users,
	}

	helper.WriteResponse(w, &response, http.StatusOK)
}
