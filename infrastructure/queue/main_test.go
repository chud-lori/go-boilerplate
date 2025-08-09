package queue_test

import (
	"os"
	"testing"

	"github.com/chud-lori/go-boilerplate/internal/testutils"
)

func TestMain(m *testing.M) {
	_ = testutils.StartRabbitOnce()
	code := m.Run()
	testutils.StopRabbit()
	os.Exit(code)
}
