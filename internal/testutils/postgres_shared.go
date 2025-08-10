package testutils

import (
    "context"
    "fmt"
    "os"
    "path/filepath"
    "sync"
    "time"

    "github.com/chud-lori/go-boilerplate/domain/ports"
    "github.com/chud-lori/go-boilerplate/infrastructure/datastore"
    "github.com/docker/go-connections/nat"
    "github.com/sirupsen/logrus"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"

    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
    postgresOnce       sync.Once
    postgresContainer  testcontainers.Container
    postgresDSN        string
    postgresInitErr    error
    postgresStopOnce   sync.Once
)

// StartPostgresOnce boots a shared PostgreSQL container and applies migrations once per test process.
func StartPostgresOnce() error {
    postgresOnce.Do(func() {
        ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
        defer cancel()

        req := testcontainers.ContainerRequest{
            Image:        "postgres:16",
            ExposedPorts: []string{"5432/tcp"},
            Env: map[string]string{
                "POSTGRES_PASSWORD": "test",
                "POSTGRES_USER":     "test",
                "POSTGRES_DB":       "testdb",
            },
            // Use log-based readiness which is more robust than mere port listening
            WaitingFor: wait.ForLog("database system is ready to accept connections").WithStartupTimeout(2 * time.Minute),
        }

        c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
            ContainerRequest: req,
            Started:          true,
        })
        if err != nil {
            postgresInitErr = err
            return
        }
        postgresContainer = c

        host, err := c.Host(ctx)
        if err != nil {
            postgresInitErr = err
            return
        }
        port, err := c.MappedPort(ctx, nat.Port("5432/tcp"))
        if err != nil {
            postgresInitErr = err
            return
        }
        // Normalize localhost to IPv4 to avoid ::1 issues on some CI runners
        if host == "localhost" || host == "::1" {
            host = "127.0.0.1"
        }
        postgresDSN = fmt.Sprintf("postgres://test:test@%s:%s/testdb?sslmode=disable", host, port.Port())

        // Warm up connection with retries to avoid race with postgres readiness
        if err := pingPostgresWithRetry(postgresDSN, 10, 500*time.Millisecond); err != nil {
            postgresInitErr = err
            return
        }

        // Apply migrations once, using best-effort path discovery
        migrationsPath := discoverMigrationsDir()
        if migrationsPath == "" {
            postgresInitErr = fmt.Errorf("could not find migrations directory from working dir")
            return
        }

        src := fmt.Sprintf("file://%s", migrationsPath)
        m, err := migrate.New(src, postgresDSN)
        if err != nil {
            postgresInitErr = err
            return
        }
        if err := m.Up(); err != nil && err != migrate.ErrNoChange {
            postgresInitErr = err
            return
        }
    })
    return postgresInitErr
}

// OpenSharedPostgres creates a new logical DB connection to the shared container.
// Caller should close the returned db when done. Container is not stopped here.
func OpenSharedPostgres() (ports.Database, error) {
    if err := StartPostgresOnce(); err != nil {
        return nil, err
    }
    logger := logrus.New()
    return datastore.NewPostgreDatabase(postgresDSN, logger)
}

// StopPostgres terminates the shared PostgreSQL container. Safe to call multiple times.
func StopPostgres() {
    postgresStopOnce.Do(func() {
        if postgresContainer != nil {
            ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
            defer cancel()
            _ = postgresContainer.Terminate(ctx)
        }
    })
}

func discoverMigrationsDir() string {
    // Try a few likely locations relative to current working directory
    candidates := []string{
        "migrations",
        filepath.Join("..", "migrations"),
        filepath.Join("..", "..", "migrations"),
        filepath.Join("..", "..", "..", "migrations"),
        filepath.Join("..", "..", "..", "..", "migrations"),
    }
    for _, rel := range candidates {
        if info, err := os.Stat(rel); err == nil && info.IsDir() {
            abs, _ := filepath.Abs(rel)
            return abs
        }
    }
    // Fallback: try from project root using env var if provided
    if root := os.Getenv("PROJECT_ROOT"); root != "" {
        p := filepath.Join(root, "migrations")
        if info, err := os.Stat(p); err == nil && info.IsDir() {
            abs, _ := filepath.Abs(p)
            return abs
        }
    }
    return ""
}

// pingPostgresWithRetry attempts to open and ping the database with retries.
func pingPostgresWithRetry(dsn string, attempts int, delay time.Duration) error {
    logger := logrus.New()
    for i := 0; i < attempts; i++ {
        db, err := datastore.NewPostgreDatabase(dsn, logger)
        if err == nil {
            _ = db.Close()
            return nil
        }
        time.Sleep(delay)
    }
    return fmt.Errorf("postgres not ready after %d attempts", attempts)
}


