package mocks

import (
	"context"

	"github.com/chud-lori/go-boilerplate/domain/entities"
	"github.com/chud-lori/go-boilerplate/domain/ports"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

// Save provides a mock function with given fields: ctx, tx, user
func (m *MockUserRepository) Save(ctx context.Context, tx ports.Transaction, user *entities.User) (*entities.User, error) {
	args := m.Called(ctx, tx, user)
	// The first argument is the returned User entity
	var r0 *entities.User
	if args.Get(0) != nil {
		r0 = args.Get(0).(*entities.User)
	}
	// The second argument is the error
	r1 := args.Error(1)
	return r0, r1
}

// Update provides a mock function with given fields: ctx, tx, user
func (m *MockUserRepository) Update(ctx context.Context, tx ports.Transaction, user *entities.User) (*entities.User, error) {
	args := m.Called(ctx, tx, user)
	var r0 *entities.User
	if args.Get(0) != nil {
		r0 = args.Get(0).(*entities.User)
	}
	r1 := args.Error(1)
	return r0, r1
}

// Delete provides a mock function with given fields: ctx, tx, id
func (m *MockUserRepository) Delete(ctx context.Context, tx ports.Transaction, id string) error {
	args := m.Called(ctx, tx, id)
	return args.Error(0) // Assuming Delete only returns an error, it's the first (and only) error
}

// FindById provides a mock function with given fields: ctx, tx, id
func (m *MockUserRepository) FindById(ctx context.Context, tx ports.Transaction, id string) (*entities.User, error) {
	args := m.Called(ctx, tx, id)
	var r0 *entities.User
	if args.Get(0) != nil {
		r0 = args.Get(0).(*entities.User)
	}
	r1 := args.Error(1)
	return r0, r1
}

// FindAll provides a mock function with given fields: ctx, tx
func (m *MockUserRepository) FindAll(ctx context.Context, tx ports.Transaction) ([]*entities.User, error) {
	args := m.Called(ctx, tx)
	var r0 []*entities.User
	if args.Get(0) != nil {
		r0 = args.Get(0).([]*entities.User)
	}
	r1 := args.Error(1)
	return r0, r1
}
