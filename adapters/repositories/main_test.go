package repositories_test

import (
	"os"
	"testing"

	"github.com/chud-lori/go-boilerplate/internal/testutils"
)

func TestMain(m *testing.M) {
	// Ensure shared Postgres is up for all repository tests
	_ = testutils.StartPostgresOnce()
	code := m.Run()
	testutils.StopPostgres()
	os.Exit(code)
}
