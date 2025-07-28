// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package log

import (
	"log/slog"
	"os"

	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
)

// createFallbackLogger returns a fallback logger that is used when the OpenTelemetry logger is not available.
// It uses the slog package to create a JSON logger that writes to stdout.
func createFallbackLogger() *slog.Logger {
	return slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level:     slog.LevelInfo,
				AddSource: true,
			},
		),
	)
}

// FallbackError logs an error message using a fallback logger. It should only be used when no instrumented logger
// is available, that is, when an unrecoverable error occurs.
func FallbackError(msg string, err error) {
	createFallbackLogger().Error("an unrecoverable error occured", slog.String("Body", msg), slog.String(string(semconv.ErrorMessageKey), err.Error()))
}
