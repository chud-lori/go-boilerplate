package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chud-lori/go-boilerplate/adapters/transport"
	"github.com/chud-lori/go-boilerplate/adapters/web/dto"
	"github.com/chud-lori/go-boilerplate/domain/ports"

	"github.com/sirupsen/logrus"
)

type UserController struct {
	ports.UserService
}

func GetPayload(request *http.Request, result interface{}) {
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(result)

	if err != nil {
		panic(err)
	}
}

func WriteResponse(writer http.ResponseWriter, response interface{}, httpCode int64) {

	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(int(httpCode))
	encoder := json.NewEncoder(writer)
	err := encoder.Encode(response)

	if err != nil {
		panic(err)
	}
}

type WebResponse struct {
	Message string      `json:"message"`
	Status  int         `json:"status"`
	Data    interface{} `json:"data"`
}

func (controller *UserController) Create(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value("logger").(*logrus.Entry)
	userRequest := dto.UserRequest{}
	GetPayload(r, &userRequest)

	userResponse, err := controller.UserService.Save(r.Context(), &userRequest)
	if err != nil {
		logger.Error("Failed to create user: ", err)
		WriteResponse(w, WebResponse{
			Message: "Failed to create user",
			Status:  0,
			Data:    nil,
		}, http.StatusBadRequest)
		return
	}

	response := WebResponse{
		Message: "success save user",
		Status:  1,
		Data:    userResponse,
	}
	WriteResponse(w, &response, http.StatusCreated)
}

func (controller *UserController) Update(w http.ResponseWriter, r *http.Request) {
	userRequest := transport.UserRequest{}
	GetPayload(r, &userRequest)

	userResponse, err := controller.UserService.Update(r.Context(), &userRequest)

	if err != nil {
		fmt.Println("Error update controller")
		panic(err)
	}

	response := WebResponse{
		Message: "success update user",
		Status:  1,
		Data:    userResponse,
	}

	WriteResponse(w, &response, http.StatusOK)
}

func (controller *UserController) Delete(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value("logger").(*logrus.Entry)
	userId := r.PathValue("userId")

	err := controller.UserService.Delete(r.Context(), userId)
	if err != nil {
		if err.Error() == "user not found" {
			WriteResponse(w, WebResponse{
				Message: "User not found",
				Status:  0,
				Data:    nil,
			}, http.StatusNotFound)
			return
		}
		logger.Error("Failed to delete user: ", err)
		WriteResponse(w, WebResponse{
			Message: "Failed to delete user",
			Status:  0,
			Data:    nil,
		}, http.StatusInternalServerError)
		return
	}

	response := WebResponse{
		Message: "success delete user",
		Status:  1,
		Data:    nil,
	}
	WriteResponse(w, &response, http.StatusOK)
}

func (c *UserController) FindById(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value("logger").(*logrus.Entry)
	userId := r.PathValue("userId")

	user, err := c.UserService.FindById(r.Context(), userId)
	if err != nil {
		if err.Error() == "user not found" {
			WriteResponse(w, WebResponse{
				Message: "User not found",
				Status:  0,
				Data:    nil,
			}, http.StatusNotFound)
			return
		}
		if err.Error() == "Invalid UUID Format" {
			WriteResponse(w, WebResponse{
				Message: "Invalid user ID format",
				Status:  0,
				Data:    nil,
			}, http.StatusBadRequest)
			return
		}
		logger.Error("Failed to find user: ", err)
		WriteResponse(w, WebResponse{
			Message: "Internal server error",
			Status:  0,
			Data:    nil,
		}, http.StatusInternalServerError)
		return
	}

	response := WebResponse{
		Message: "success get user by id",
		Status:  1,
		Data:    user,
	}
	WriteResponse(w, &response, http.StatusOK)
}

func (controller *UserController) FindAll(w http.ResponseWriter, r *http.Request) {
	logger, _ := r.Context().Value("logger").(*logrus.Entry)

	users, err := controller.UserService.FindAll(r.Context())

	if err != nil {
		logger.Info("Error Find All user: ", err)
		WriteResponse(w, WebResponse{
			Message: "Failed Get All Users",
			Status:  0,
			Data:    nil,
		}, http.StatusNotFound)
		return
	}

	response := WebResponse{
		Message: "success get all users",
		Status:  1,
		Data:    users,
	}

	WriteResponse(w, &response, http.StatusOK)
}
