package controllers

import (
	"errors"
	"net/http"

	"github.com/chud-lori/go-boilerplate/adapters/web/dto"
	"github.com/chud-lori/go-boilerplate/adapters/web/helper"
	"github.com/chud-lori/go-boilerplate/domain/entities"
	"github.com/chud-lori/go-boilerplate/domain/ports"
	appErrors "github.com/chud-lori/go-boilerplate/pkg/errors"
	"github.com/google/uuid"

	"github.com/sirupsen/logrus"
)

type UserController struct {
	ports.UserService
}

func (controller *UserController) Create(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value("logger").(*logrus.Entry)

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

	result, err := controller.UserService.Save(r.Context(), userPayload)

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
			CreatedAt: result.Created_at,
		},
	}
	helper.WriteResponse(w, &response, http.StatusCreated)
}

func (controller *UserController) Update(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value("logger").(*logrus.Entry)
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

	userResponse, err := controller.UserService.Update(r.Context(), userPayload)

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

func (controller *UserController) Delete(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("userId")

	err := controller.UserService.Delete(r.Context(), userId)

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

func (c *UserController) FindById(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value("logger").(*logrus.Entry)
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

	user, err := c.UserService.FindById(r.Context(), userId)
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

func (controller *UserController) FindAll(w http.ResponseWriter, r *http.Request) {
	logger, _ := r.Context().Value("logger").(*logrus.Entry)

	users, err := controller.UserService.FindAll(r.Context())

	if err != nil {
		logger.Error("Error Find All user: ", err)
		helper.WriteResponse(w, dto.WebResponse{
			Message: "Failed Get All Users",
			Status:  0,
			Data:    nil,
		}, http.StatusNotFound)
		return
	}

	response := dto.WebResponse{
		Message: "success get all users",
		Status:  1,
		Data:    users,
	}

	helper.WriteResponse(w, &response, http.StatusOK)
}
