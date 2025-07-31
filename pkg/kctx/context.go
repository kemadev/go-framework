// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package kctx

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
)

var kctxPool = sync.Pool{
	New: func() interface{} {
		return &Kctx{}
	},
}

type contextKey struct{}

var kctxKey = contextKey{}

type Kctx struct {
	context.Context
}

func (ctx *Kctx) release() {
	ctx.Context = nil
	kctxPool.Put(ctx)
}

// FromRequest extracts Kctx from [net/http.Request].
func FromRequest(r *http.Request) *Kctx {
	newCtx, found := fromContext(r.Context())
	if !found {
		slog.Warn("kctx not found in context")
	}

	return newCtx
}

// FromContext extracts Kctx from [context.Context].
func FromContext(c context.Context) *Kctx {
	newCtx, found := fromContext(c)
	if !found {
		slog.Warn("kctx not found in context")
	}

	return newCtx
}

func fromContext(c context.Context) (*Kctx, bool) {
	if kctx, ok := c.Value(kctxKey).(*Kctx); ok {
		return kctx, true
	}
	// Fallback
	return &Kctx{Context: c}, false
}

// Middleware manages Kctx lifecycle in a [sync.Pool].
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := kctxPool.Get().(*Kctx)
		defer c.release()

		c.Context = r.Context()

		ctx := context.WithValue(r.Context(), kctxKey, c)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
