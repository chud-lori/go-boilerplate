package testutils

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/chud-lori/go-boilerplate/domain/ports"
	"github.com/chud-lori/go-boilerplate/infrastructure/datastore"
	"github.com/chud-lori/go-boilerplate/pkg/logger"
	"github.com/docker/go-connections/nat"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	// Import the migrate library and its drivers
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // PostgreSQL database driver
	_ "github.com/golang-migrate/migrate/v4/source/file"       // File source driver
)

func SetupTestDBWithTestcontainers(t *testing.T) (ports.Database, func()) {
	// Use a context with a timeout for container operations
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// YOUR ORIGINAL CONTAINER SETUP - KEPT AS IS
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

	// --- APPLY MIGRATIONS USING GOLANG-MIGRATE ---
	// Construct the source URL for your migration files.
	// This assumes your migrations are in a directory named 'migrations'
	// relative to where your tests are being executed.
	migrationSourceURL := "file://../../migrations"

	m, err := migrate.New(
		migrationSourceURL,
		dsn,
	)
	require.NoError(t, err, "Failed to create migrate instance")

	// Apply all up migrations
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		// ErrNoChange means no migrations were applied, which is often fine.
		// Any other error should cause the test to fail.
		require.NoError(t, err, "Failed to apply database migrations")
	}
	// --- END APPLY MIGRATIONS ---

	// Cleanup function
	terminate := func() {
		// Close your database connection
		if db != nil {
			db.Close()
		}
		// Terminate the testcontainer
		if container != nil {
			timeoutCtx, cancelTimeout := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancelTimeout()
			container.Terminate(timeoutCtx)
		}
	}

	return db, terminate
}

// RunInTestTransaction sets up a database transaction for a test and handles rollback.
// It accepts a test function that will receive the context and the transaction.
// The `db` parameter is the *sql.DB connection (expected to be from SetupTestDBWithTestcontainers).
func RunInTestTransaction(t *testing.T, db *sql.DB, testFunc func(ctx context.Context, tx *sql.Tx)) {
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.New()) // Ensure logger is available

	tx, err := db.BeginTx(ctx, nil) // Using nil options for default isolation
	require.NoError(t, err, "Failed to begin transaction for test")

	// Defer rollback to ensure the transaction is always rolled back,
	// keeping the test database clean and isolated from other tests.
	defer func() {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			t.Errorf("Error during test transaction rollback: %v", rollbackErr)
		}
	}()

	testFunc(ctx, tx)
}
