package ports

import (
	"context"
	"database/sql"
)

// DBTX interface for database operations (shared by *sql.DB and *sql.Tx)
type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

// Transaction interface (wraps DBTX)
type Transaction interface {
	DBTX // Embed the DBTX interface
	Commit() error
	Rollback() error
}

// Database interface
type Database interface {
	BeginTx(ctx context.Context) (Transaction, error)
	Close() error
}
