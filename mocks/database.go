package mocks

import (
	"context"
	"database/sql"

	"github.com/chud-lori/go-boilerplate/domain/ports"
	"github.com/stretchr/testify/mock"
)

type MockRow struct {
	scanFunc func(dest ...interface{}) error
}

func (m *MockRow) Scan(dest ...interface{}) error {
	if m.scanFunc != nil {
		return m.scanFunc(dest...)
	}
	return nil
}

func NewMockRowWithScan(scanFunc func(dest ...interface{}) error) *sql.Row {
	return &sql.Row{}
}

type MockTransaction struct {
	mock.Mock
}

func (m *MockTransaction) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	callArgs := make([]interface{}, 0, len(args)+2)
	callArgs = append(callArgs, ctx, query)
	callArgs = append(callArgs, args...)

	argsM := m.Called(callArgs...)
	return argsM.Get(0).(sql.Result), argsM.Error(1)
}

func (m *MockTransaction) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	callArgs := make([]interface{}, 0, len(args)+2)
	callArgs = append(callArgs, ctx, query)
	callArgs = append(callArgs, args...)

	argsM := m.Called(callArgs...)
	return argsM.Get(0).(*sql.Row)
}

func (m *MockTransaction) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	callArgs := make([]interface{}, 0, len(args)+2)
	callArgs = append(callArgs, ctx, query)
	callArgs = append(callArgs, args...)

	argsM := m.Called(callArgs...)
	return argsM.Get(0).(*sql.Rows), argsM.Error(1)
}

func (m *MockTransaction) Commit() error {
	return m.Called().Error(0)
}

func (m *MockTransaction) Rollback() error {
	return m.Called().Error(0)
}

type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) BeginTx(ctx context.Context) (ports.Transaction, error) {
	args := m.Called(ctx)
	return args.Get(0).(ports.Transaction), args.Error(1)
}

func (m *MockDatabase) Close() error {
	return m.Called().Error(0)
}
