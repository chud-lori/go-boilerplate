package testutils

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/chud-lori/go-boilerplate/domain/ports"
	"github.com/chud-lori/go-boilerplate/infrastructure/datastore"
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
