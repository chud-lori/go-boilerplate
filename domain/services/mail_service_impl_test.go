package services_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/chud-lori/go-boilerplate/domain/services"
	"github.com/chud-lori/go-boilerplate/mocks"
	"github.com/stretchr/testify/assert"
)

func TestMailService_SendSignInNotification_Success(t *testing.T) {
	mockMailClient := new(mocks.MockMailClient) // Assuming you have a mock for MailClient

	service := &services.MailServiceImpl{
		MailClient: mockMailClient,
	}

	email := "test@example.com"
	notificationText := "User logged in just now" // This is the 'text' argument for SendSignInNotification

	// The message constructed inside SendSignInNotification is "Notif sent"
	expectedMailMessage := fmt.Sprintf("Notif sent")

	// Set up expectations for the mock MailClient
	// We expect SendMail to be called with the email and the *constructed* message
	mockMailClient.On("SendMail", email, expectedMailMessage).Return(nil)

	// Call the service method
	err := service.SendSignInNotification(context.Background(), email, notificationText)

	// Assertions
	assert.NoError(t, err) // Expect no error

	// Assert that all expectations on the mock were met
	mockMailClient.AssertExpectations(t)
}

func TestMailService_SendSignInNotification_MailClientError(t *testing.T) {
	mockMailClient := new(mocks.MockMailClient)
	service := &services.MailServiceImpl{
		MailClient: mockMailClient,
	}

	email := "test@example.com"
	notificationText := "User logged in just now"
	expectedMailMessage := fmt.Sprintf("Notif sent")

	// Simulate an error from the MailClient
	mailClientErr := fmt.Errorf("failed to send mail via client")
	mockMailClient.On("SendMail", email, expectedMailMessage).Return(mailClientErr)

	// Call the service method
	err := service.SendSignInNotification(context.Background(), email, notificationText)

	// Assertions
	assert.Error(t, err)                             // Expect an error
	assert.EqualError(t, err, mailClientErr.Error()) // Check if the error matches
	mockMailClient.AssertExpectations(t)
}
