package ports

import (
	"context"
	"database/sql"
)

// New interface for database connection management
type Database interface {
    BeginTx(context.Context) (Transaction, error)
    Close() error
}

type Transaction interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	Commit() error
	Rollback() error
}
