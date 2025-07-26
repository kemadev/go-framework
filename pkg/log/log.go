// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package log

import (
	"log/slog"
	"os"
)

// CreateFallbackLogger returns a fallback logger that is used when the OpenTelemetry logger is not available.
// It uses the slog package to create a JSON logger that writes to stdout.
func CreateFallbackLogger() *slog.Logger {
	return slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		),
	)
}
