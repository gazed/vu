// Copyright © 2024 Galvanized Logic Inc.

package vu

import (
	"log/slog"
	"os"
	"testing"
)

// TestMain is called by "go test" instead of running the tests individually.
// It is used to setup and teardown state for all tests.
func TestMain(m *testing.M) {

	// configure the default logger to log everything during tests.
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))
	m.Run() // run individual tests

	// no teardown for now.
}
