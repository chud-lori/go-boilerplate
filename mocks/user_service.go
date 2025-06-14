package mocks

import (
	"context"

	"github.com/chud-lori/go-boilerplate/domain/entities"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Save(ctx context.Context, user *entities.User) (*entities.User, error) {
	args := m.Called(ctx, user)
	if result := args.Get(0); result != nil {
		return result.(*entities.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserService) Update(ctx context.Context, user *entities.User) (*entities.User, error) {
	args := m.Called(ctx, user)
	if result := args.Get(0); result != nil {
		return result.(*entities.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserService) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserService) FindById(ctx context.Context, id string) (*entities.User, error) {
	args := m.Called(ctx, id)
	if result := args.Get(0); result != nil {
		return result.(*entities.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserService) FindAll(ctx context.Context) ([]*entities.User, error) {
	args := m.Called(ctx)
	if result := args.Get(0); result != nil {
		return result.([]*entities.User), args.Error(1)
	}
	return nil, args.Error(1)
}
