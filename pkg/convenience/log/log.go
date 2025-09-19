// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package log

import (
	"log/slog"

	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

// Logger retrieves the logger nammed after name. Name should be the package name.
func Logger(name string) *slog.Logger {
	return GetPackageLogger(name)
}

// ErrLog retrieves the logger nammed after name and logs [msg], attaching [err] as attribute.
func ErrLog(name string, msg string, err error) {
	GetPackageLogger(
		name,
	).Error("error occured", slog.String(string(semconv.ErrorMessageKey), err.Error()))
}
