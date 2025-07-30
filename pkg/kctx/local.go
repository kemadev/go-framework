// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package kctx

import "context"

type localKey struct {
	name string
}

// Local retrieves a request-scoped value, set by [LocalSet].
func (ctx *Kctx) Local(name string) any {
	return ctx.Value(localKey{name: name})
}

// Local sets a request-scoped value, which can be retrieved using [Local] later on.
func (ctx *Kctx) LocalSet(name string, value any) {
	ctx.Context = context.WithValue(ctx.Context, localKey{name: name}, value)
}
