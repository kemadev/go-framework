// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package kctx

import (
	"log/slog"

	"github.com/kemadev/go-framework/pkg/log"
)

// Logger retrieves the logger nammed after name.
func (c *Kctx) Logger(name string) *slog.Logger {
	return log.GetPackageLogger(name)
}
