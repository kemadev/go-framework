package kctx

import "github.com/kemadev/go-framework/pkg/log"

// Logger retrieves the logger nammed after name
func (ctx *Kctx) Logger(name string) any {
	return log.GetPackageLogger(name)
}
