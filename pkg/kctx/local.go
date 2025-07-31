// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package kctx

import "context"

type localKey struct {
	name string
}

// Local retrieves a request-scoped value, set by [LocalSet].
func (c *Kctx) Local(name string) any {
	return c.Value(localKey{name: name})
}

// Local sets a request-scoped value, which can be retrieved using [Local] later on.
func (c *Kctx) LocalSet(name string, value any) {
	c.Context = context.WithValue(c.Context, localKey{name: name}, value)
}
