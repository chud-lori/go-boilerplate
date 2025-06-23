package errors

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationErrors struct {
	Messages []string
}

func (ve *ValidationErrors) Error() string {
	return fmt.Sprintf("validation failed: %s", strings.Join(ve.Messages, "; "))
}

func NewValidationErrors(errs validator.ValidationErrors) *ValidationErrors {
	var messages []string
	for _, err := range errs {
		// Customize error messages as needed
		switch err.Tag() {
		case "required":
			messages = append(messages, fmt.Sprintf("%s is required", strings.ToLower(err.Field())))
		case "min":
			if err.Field() == "Password" { // Specific message for password min length
				messages = append(messages, "Password must be at least 8 characters long")
			} else {
				messages = append(messages, fmt.Sprintf("%s must be at least %s characters long", strings.ToLower(err.Field()), err.Param()))
			}
		case "email":
			messages = append(messages, fmt.Sprintf("Invalid email format for %s", strings.ToLower(err.Field())))
		case "eqfield":
			if err.Field() == "ConfirmPassword" && err.Param() == "Password" { // Specific message for password mismatch
				messages = append(messages, "Password and confirm password do not match")
			} else {
				messages = append(messages, fmt.Sprintf("%s does not match %s", strings.ToLower(err.Field()), strings.ToLower(err.Param())))
			}
		// Add more cases for other validator tags as needed
		default:
			messages = append(messages, fmt.Sprintf("Validation error on %s: %s", strings.ToLower(err.Field()), err.Tag()))
		}
	}
	return &ValidationErrors{Messages: messages}
}

// IsValidationErrors checks if the given error is a ValidationErrors type.
func IsValidationErrors(err error) bool {
	_, ok := err.(*ValidationErrors)
	return ok
}
