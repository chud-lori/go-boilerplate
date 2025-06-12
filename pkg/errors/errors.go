package errors

import "net/http"

type AppError struct {
	Message    string
	StatusCode int
	Err        error
}

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
