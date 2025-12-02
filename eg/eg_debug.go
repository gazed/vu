// SPDX-FileCopyrightText : © 2024-2025 Galvanized Logic Inc.
// SPDX-License-Identifier: BSD-2-Clause

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
