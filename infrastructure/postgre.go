package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"
	"github.com/chud-lori/go-boilerplate/domain/ports"
	_ "github.com/lib/pq"
)

// compile-time interface check
var _ ports.DBTransaction = (*PostgresTransaction)(nil)

// PostgresDB implements DB interface
type PostgresDB struct {
    db *sql.DB
}

// PostgresTransaction implements DBTransaction interface
type PostgresTransaction struct {
    tx *sql.Tx
}

// Only transaction methods for PostgresTransaction
func (t *PostgresTransaction) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
    return t.tx.ExecContext(ctx, query, args...)
}

func (t *PostgresTransaction) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
    return t.tx.QueryRowContext(ctx, query, args...)
}

func (t *PostgresTransaction) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
    return t.tx.QueryContext(ctx, query, args...)
}

func (t *PostgresTransaction) Commit() error {
    return t.tx.Commit()
}

func (t *PostgresTransaction) Rollback() error {
    return t.tx.Rollback()
}

// Only connection methods for PostgresDB
func (p *PostgresDB) BeginTx(ctx context.Context) (ports.DBTransaction, error) {
    tx, err := p.db.BeginTx(ctx, &sql.TxOptions{
        Isolation: sql.LevelReadCommitted,
    })
    if err != nil {
        return nil, err
    }
    return &PostgresTransaction{tx: tx}, nil
}

func (p *PostgresDB) Close() error {
    return p.db.Close()
}

func NewPostgreDB() ports.DB {
	var (
		host     = os.Getenv("PSQL_HOST")
		port     = 5432
		user     = os.Getenv("PSQL_USER")
		password = os.Getenv("PSQL_PASSWORD")
		dbname   = os.Getenv("DB_NAME")
	)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return &PostgresDB{db: db}
}
