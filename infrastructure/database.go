package infrastructure

import (
	"context"
	"database/sql"
	"net/url"
	"fmt"
	"time"
	"github.com/chud-lori/go-boilerplate/domain/ports"
	_ "github.com/lib/pq"
)

// compile-time interface check
var _ ports.Transaction = (*Transaction)(nil)
var _ ports.Database = (*Database)(nil)

// Database implements DB interface
type Database struct {
    db *sql.DB
}

// Transaction implements Transaction interface
type Transaction struct {
    tx *sql.Tx
}

// Only transaction methods for Transaction
func (t *Transaction) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
    return t.tx.ExecContext(ctx, query, args...)
}

func (t *Transaction) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
    return t.tx.QueryRowContext(ctx, query, args...)
}

func (t *Transaction) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
    return t.tx.QueryContext(ctx, query, args...)
}

func (t *Transaction) Commit() error {
    return t.tx.Commit()
}

func (t *Transaction) Rollback() error {
    return t.tx.Rollback()
}

// Only connection methods for Database
func (p *Database) BeginTx(ctx context.Context) (ports.Transaction, error) {
    tx, err := p.db.BeginTx(ctx, &sql.TxOptions{
        Isolation: sql.LevelReadCommitted,
    })
    if err != nil {
        return nil, err
    }
    return &Transaction{tx: tx}, nil
}

func (p *Database) Close() error {
    return p.db.Close()
}

func NewDatabase(dbURL string) (ports.Database, error) {
	parseDBUrl, _ := url.Parse(dbURL)
	db, err := sql.Open(parseDBUrl.Scheme, dbURL)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return &Database{db: db}, nil
}
