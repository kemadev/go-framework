package kctx

import "context"

// Local retrieves a request-scoped value, set by [LocalSet]
func (ctx *Kctx) Local(name string) any {
	return ctx.Value(name)
}

// Local sets a request-scoped value, which can be retrieved using [Local] later on
func (ctx *Kctx) LocalSet(name string, value any) {
	ctx.Context = context.WithValue(ctx.Context, name, value)
}
