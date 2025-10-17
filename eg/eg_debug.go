// Copyright © 2024 Galvanized Logic Inc.

//go:build debug

package main

// eg_debug.go turns on structured logging debug logs
// when building with "go build -tags debug"

import (
	"log/slog"
	"os"
)

// init runs after globals vars are initialized and before main is called
// so this overrides the configLogging called at startup.
func init() {
	configLogging = func() {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))
	}
}
