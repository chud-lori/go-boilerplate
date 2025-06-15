package testutils

import (
	"context"
	"fmt"
	"testing"

	"github.com/chud-lori/go-boilerplate/domain/ports"
	"github.com/chud-lori/go-boilerplate/infrastructure/datastore"
	"github.com/docker/go-connections/nat"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// func SetupTestDB(t *testing.T) ports.Database {
// 	err := godotenv.Load("../../.env.test")
// 	require.NoError(t, err, "failed to load .env.test")

// 	dbURL := os.Getenv("DB_URL")
// 	require.NotEmpty(t, dbURL, "DB_URL must be set in .env.test")

// 	baseLogger := logrus.New()
// 	db, err := datastore.NewDatabase(dbURL, baseLogger)
// 	require.NoError(t, err, "failed to connect to test database")

// 	return db
// }

// func SetupTestTx(t *testing.T, db ports.Database) ports.Transaction {
// 	ctx := context.Background()
// 	tx, err := db.BeginTx(ctx)
// 	require.NoError(t, err, "failed to begin transaction")
// 	return tx
// }

func SetupTestDBWithTestcontainers(t *testing.T) (ports.Database, func()) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:16",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_USER":     "test",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	// Get host/port to construct DSN
	host, err := container.Host(ctx)
	require.NoError(t, err)

	port, err := container.MappedPort(ctx, nat.Port("5432/tcp"))
	require.NoError(t, err)

	dsn := fmt.Sprintf("postgres://test:test@%s:%s/testdb?sslmode=disable", host, port.Port())

	// Create your custom Database instance
	logger := logrus.New()
	db, err := datastore.NewDatabase(dsn, logger)
	require.NoError(t, err)

	// Create schema
	tx, err := db.BeginTx(ctx)
	require.NoError(t, err)
	_, err = tx.ExecContext(ctx, `
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			email VARCHAR(255) UNIQUE NOT NULL,
			passcode VARCHAR(255),
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`)
	require.NoError(t, err)
	require.NoError(t, tx.Commit())

	// Cleanup
	terminate := func() {
		db.Close()
		container.Terminate(ctx)
	}

	return db, terminate
}

// SetupTestTx opens a transaction for testing.
// func SetupTestTx(t *testing.T, db *sql.DB) *sql.Tx {
// 	t.Helper()
// 	tx, err := db.Begin()
// 	require.NoError(t, err, "failed to begin tx")
// 	return tx
// }
