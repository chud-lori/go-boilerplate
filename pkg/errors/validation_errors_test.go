package errors_test

import (
	"testing"

	appErr "github.com/chud-lori/go-boilerplate/pkg/errors"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type TestInput struct {
	Email           string `validate:"required,email"`
	Username        string `validate:"required,min=5"`
	Password        string `validate:"required,min=8"`
	ConfirmPassword string `validate:"required,eqfield=Password"`
}

func TestNewValidationErrors(t *testing.T) {
	validate := validator.New()

	input := TestInput{
		Email:           "invalid-email", // invalid format
		Username:        "abc",           // too short
		Password:        "123",           // too short
		ConfirmPassword: "notmatch",      // does not match Password
	}

	err := validate.Struct(input)
	assert.Error(t, err)

	ve, ok := err.(validator.ValidationErrors)
	assert.True(t, ok)

	validationErr := appErr.NewValidationErrors(ve)

	assert.Contains(t, validationErr.Messages, "Invalid email format for email")
	assert.Contains(t, validationErr.Messages, "username must be at least 5 characters long")
	assert.Contains(t, validationErr.Messages, "Password must be at least 8 characters long")
	assert.Contains(t, validationErr.Messages, "Password and confirm password do not match")
}

func TestIsValidationErrors(t *testing.T) {
	err := &appErr.ValidationErrors{Messages: []string{"test"}}
	assert.True(t, appErr.IsValidationErrors(err))

	nonValidationErr := assert.AnError
	assert.False(t, appErr.IsValidationErrors(nonValidationErr))
}
