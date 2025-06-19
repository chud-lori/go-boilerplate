package errors

import (
	"errors"
	"net/http"
)

type AppError struct {
	Message    string
	StatusCode int
	Err        error
}

var ErrUserNotFound = errors.New("user not found")

func (e *AppError) Error() string {
	return e.Message
}

func NewBadRequestError(message string, err error) *AppError {
	return &AppError{
		Message:    message,
		StatusCode: http.StatusBadRequest,
		Err:        err,
	}
}

func NewInternalServerError(message string, err error) *AppError {
	return &AppError{
		Message:    message,
		StatusCode: http.StatusInternalServerError,
		Err:        err,
	}
}

func NewNotFoundError(message string, err error) *AppError {
	return &AppError{
		Message:    message,
		StatusCode: http.StatusNotFound,
		Err:        err,
	}
}

func NewUnauthorizedError(message string, err error) *AppError {
	return &AppError{
		Message:    message,
		StatusCode: http.StatusUnauthorized,
		Err:        err,
	}
}
