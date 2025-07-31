// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package log

import (
	"log/slog"
)

// Logger retrieves the logger nammed after name. Name should be the package name.
func Logger(name string) *slog.Logger {
	return GetPackageLogger(name)
}
