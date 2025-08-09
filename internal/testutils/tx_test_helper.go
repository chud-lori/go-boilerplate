package testutils

import (
	"context"
	"testing"

	"github.com/chud-lori/go-boilerplate/domain/ports"
	"github.com/chud-lori/go-boilerplate/pkg/logger"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

type RepoSetupFunc func(db ports.Database) (any, error)

// Generic transaction test wrapper
func WithTransactionTest[T any](
	t *testing.T,
	setupRepo func(db ports.Database) (T, error),
	testFunc func(ctx context.Context, repo T, tx ports.Transaction),
) {
    // Use shared PostgreSQL container to avoid per-test container spin-up
    sharedDB, err := OpenSharedPostgres()
    require.NoError(t, err)
    t.Cleanup(func() { sharedDB.Close() })

	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.New())
    tx, err := sharedDB.BeginTx(ctx)
	require.NoError(t, err)

    repo, err := setupRepo(sharedDB)
	require.NoError(t, err)

	testFunc(ctx, repo, tx)

	_ = tx.Rollback()
}
