package errors_test

import (
	"errors"
	"net/http"
	"testing"

	appErr "github.com/chud-lori/go-boilerplate/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewBadRequestError(t *testing.T) {
	orig := errors.New("invalid input")
	err := appErr.NewBadRequestError("bad request", orig)

	assert.Equal(t, "bad request", err.Message)
	assert.Equal(t, http.StatusBadRequest, err.StatusCode)
	assert.Equal(t, orig, err.Err)
}

func TestNewInternalServerError(t *testing.T) {
	orig := errors.New("something broke")
	err := appErr.NewInternalServerError("internal error", orig)

	assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
}

func TestNewNotFoundError(t *testing.T) {
	err := appErr.NewNotFoundError("user not found", appErr.ErrUserNotFound)

	assert.Equal(t, http.StatusNotFound, err.StatusCode)
	assert.Equal(t, "user not found", err.Message)
}
