package ports

import (
	"context"
	"database/sql"
)

// New interface for database connection management
type DB interface {
    BeginTx(context.Context) (DBTransaction, error)
    Close() error
}

type DBTransaction interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	Commit() error
	Rollback() error
}
