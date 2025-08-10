package locking_test

import (
    "os"
    "testing"

    "github.com/chud-lori/go-boilerplate/internal/testutils"
)

func TestMain(m *testing.M) {
    _ = testutils.StartRedisOnce()
    code := m.Run()
    testutils.StopRedis()
    os.Exit(code)
}


