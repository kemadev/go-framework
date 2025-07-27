// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package log

import (
	"log/slog"
	"os"

	"github.com/kemadev/go-framework/pkg/config"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
)

// CreateFallbackLogger returns a fallback logger that is used when the OpenTelemetry logger is not available.
// It uses the slog package to create a JSON logger that writes to stdout.
func CreateFallbackLogger(conf config.Runtime) *slog.Logger {
	var handler slog.Handler
	if conf.IsLocalEnvironment() {
		handler = slog.NewTextHandler(os.Stdout,
			&slog.HandlerOptions{
				Level:     slog.LevelDebug,
				AddSource: true,
			},
		)
	} else {
		handler = slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level:     slog.LevelInfo,
				AddSource: true,
			},
		)
	}
	return slog.New(handler).With(
		slog.String(string(semconv.ServiceNameKey), conf.AppName),
		slog.String(string(semconv.ServiceNamespaceKey), conf.AppNamespace),
		slog.String(string(semconv.ServiceVersionKey), conf.AppVersion),
	)
}
