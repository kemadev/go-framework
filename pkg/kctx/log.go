package kctx

import (
	"log/slog"

	"github.com/kemadev/go-framework/pkg/log"
)

// Logger retrieves the logger nammed after name
func (ctx *Kctx) Logger(name string) *slog.Logger {
	return log.GetPackageLogger(name)
}
