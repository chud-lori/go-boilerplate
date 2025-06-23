package mocks

import (
	"github.com/stretchr/testify/mock"
)

// MockTokenManager is a mock implementation of the TokenManager interface
// using testify/mock.
type MockTokenManager struct {
	mock.Mock
}

// GenerateToken mocks the GenerateToken method of the TokenManager interface.
// It records the call and returns the values configured by the expectations.
func (m *MockTokenManager) GenerateToken(userID string) (string, error) {
	args := m.Called(userID)
	var r0 string
	if args.Get(0) != nil {
		r0 = args.Get(0).(string)
	}
	return r0, args.Error(1)
}

// ValidateToken mocks the ValidateToken method of the TokenManager interface.
// It records the call and returns the values configured by the expectations (userID and error).
func (m *MockTokenManager) ValidateToken(tokenStr string) (string, error) {
	args := m.Called(tokenStr)
	var r0 string
	if args.Get(0) != nil {
		r0 = args.Get(0).(string)
	}
	return r0, args.Error(1)
}
