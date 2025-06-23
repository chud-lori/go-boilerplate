package mocks

import (
	"context"

	"github.com/chud-lori/go-boilerplate/domain/entities"
	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

// SignIn mocks the SignIn method of the AuthService interface.
// It records the call and returns the values configured by the expectations.
func (m *MockAuthService) SignIn(ctx context.Context, user *entities.User) (*entities.User, string, error) {
	args := m.Called(ctx, user)
	var r0 *entities.User
	if args.Get(0) != nil {
		r0 = args.Get(0).(*entities.User)
	}
	var r1 string
	if args.Get(1) != nil {
		r1 = args.Get(1).(string)
	}
	return r0, r1, args.Error(2)
}

// SignUp mocks the SignUp method of the AuthService interface.
// It records the call and returns the values configured by the expectations.
func (m *MockAuthService) SignUp(ctx context.Context, user *entities.User) (*entities.User, string, error) {
	args := m.Called(ctx, user)
	var r0 *entities.User
	if args.Get(0) != nil {
		r0 = args.Get(0).(*entities.User)
	}
	var r1 string
	if args.Get(1) != nil {
		r1 = args.Get(1).(string)
	}
	return r0, r1, args.Error(2)
}
