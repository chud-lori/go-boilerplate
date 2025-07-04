package services_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/chud-lori/go-boilerplate/domain/services"
	"github.com/chud-lori/go-boilerplate/mocks"
	"github.com/chud-lori/go-boilerplate/pkg/logger"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMailService_SendSignInNotification_Success(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	mockMailClient := new(mocks.MockMailClient) // Assuming you have a mock for MailClient

	service := &services.MailServiceImpl{
		MailClient: mockMailClient,
	}

	email := "test@example.com"
	notificationText := "User logged in just now" // This is the 'text' argument for SendSignInNotification

	// Set up expectations for the mock MailClient
	// We expect SendMail to be called with the email and the *constructed* message
	mockMailClient.On("SendMail", mock.Anything, email, notificationText).Return(nil)

	// Call the service method
	err := service.SendSignInNotification(ctx, email, notificationText)

	// Assertions
	assert.NoError(t, err) // Expect no error

	// Assert that all expectations on the mock were met
	mockMailClient.AssertExpectations(t)
}

func TestMailService_SendSignInNotification_MailClientError(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	mockMailClient := new(mocks.MockMailClient)
	service := &services.MailServiceImpl{
		MailClient: mockMailClient,
	}

	email := "test@example.com"
	notificationText := "User logged in just now"

	// Simulate an error from the MailClient
	mailClientErr := fmt.Errorf("failed to send mail via client")
	mockMailClient.On("SendMail", mock.Anything, email, notificationText).Return(mailClientErr)

	// Call the service method
	err := service.SendSignInNotification(ctx, email, notificationText)

	// Assertions
	assert.Error(t, err)                             // Expect an error
	assert.EqualError(t, err, mailClientErr.Error()) // Check if the error matches
	mockMailClient.AssertExpectations(t)
}
