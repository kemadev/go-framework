// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package kctx

import (
	"context"
	"net/http"
	"sync"
)

const packageName = "github.com/kemadev/go-framework/pkg/kctx"

var kctxPool = sync.Pool{
	New: func() interface{} {
		return &Kctx{}
	},
}

type contextKey struct{}

var kctxKey = contextKey{}

type Kctx struct {
	context.Context
	w http.ResponseWriter
	r *http.Request
}

func (c *Kctx) release() {
	c.Context = nil
	c.w = nil
	c.r = nil
	kctxPool.Put(c)
}

// FromRequest extracts Kctx from [net/http.Request]. If context is not found, it returns nil.
func FromRequest(r *http.Request) *Kctx {
	return fromContext(r.Context())
}

// FromContext extracts Kctx from [context.Context]. If context is not found, it returns nil.
func FromContext(c context.Context) *Kctx {
	return fromContext(c)
}

func fromContext(c context.Context) *Kctx {
	kctx, ok := c.Value(kctxKey).(*Kctx)
	if !ok {
		return nil
	}
	return kctx
}

func normalizeHeaders(headers http.Header) {
	for key := range headers {
		canonical := http.CanonicalHeaderKey(key)
		if key != canonical {
			headers[canonical] = headers[key]
			delete(headers, key)
		}
	}
}

// Middleware manages Kctx lifecycle in a [sync.Pool]. It populates initializes a kctx instance,
// populates kctx key in context and propagates it down the chain. It also normalizes request headers.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := kctxPool.Get().(*Kctx)
		defer c.release()

		normalizeHeaders(r.Header)

		reqCtx := r.Context()
		c.Context = reqCtx
		c.w = w
		c.r = r

		ctx := context.WithValue(reqCtx, kctxKey, c)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
