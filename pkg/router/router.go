// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

// [this nice article]: https://www.alexedwards.net/blog/organize-your-go-middleware-without-dependencies
// [Alex Edwards]: https://www.alexedwards.net/
package router

import (
	"context"
	"net/http"
	"slices"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

const ServerRootSpanName = "server"

type chain []func(http.Handler) http.Handler

func (c chain) thenFunc(h http.HandlerFunc) http.Handler {
	return c.then(h)
}

func (c chain) then(h http.Handler) http.Handler {
	for _, mw := range slices.Backward(c) {
		h = mw(h)
	}

	return h
}

type PatternHolder struct {
	Pattern string
}
type PatternKey struct{}

func injectPattern(pattern string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		holder, ok := r.Context().Value(PatternKey{}).(*PatternHolder)
		if ok {
			holder.Pattern = pattern
		}

		next.ServeHTTP(w, r)
	})
}

// Router is an HTTP router.
type Router struct {
	globalChain []func(http.Handler) http.Handler
	routeChain  []func(http.Handler) http.Handler
	isSubRouter bool
	*http.ServeMux
}

// New returns a new HTTP router.
func New() *Router {
	return &Router{ServeMux: http.NewServeMux()}
}

// Use appends [mw] to the routers chain.
func (r *Router) Use(mw ...func(http.Handler) http.Handler) {
	if r.isSubRouter {
		r.routeChain = append(r.routeChain, mw...)
	} else {
		r.globalChain = append(r.globalChain, mw...)
	}
}

// UseInstrumented appends [mw] to the routers chain, wrapping the handler with OpenTelemetry instrumentation.
func (r *Router) UseInstrumented(name string, mw func(http.Handler) http.Handler) {
	instrumentedMw := func(next http.Handler) http.Handler {
		return otelhttp.NewHandler(mw(next), name)
	}
	if r.isSubRouter {
		r.routeChain = append(r.routeChain, instrumentedMw)
	} else {
		r.globalChain = append(r.globalChain, instrumentedMw)
	}
}

// Group adds all routers down the chain to a group. All members of a group inherits from
// their parent's routers chain.
func (r *Router) Group(group func(r *Router)) {
	subRouter := &Router{
		routeChain:  slices.Clone(r.routeChain),
		isSubRouter: true,
		ServeMux:    r.ServeMux,
	}
	group(subRouter)
}

// HandleFunc registers a handler function for a pattern.
func (r *Router) HandleFunc(pattern string, h http.HandlerFunc) {
	r.Handle(pattern, h)
}

// Handle registers a handler for a pattern, automatically injecting the pattern into the context.
func (r *Router) Handle(pattern string, h http.Handler) {
	for _, mw := range slices.Backward(r.routeChain) {
		h = mw(h)
	}

	r.ServeMux.Handle(pattern, h)
}

// HandleInstrumented registers a handler with otelhttp instrumentation and pattern injection.
func (r *Router) HandleInstrumented(pattern string, h http.Handler) {
	h = injectPattern(pattern, h)
	r.Handle(pattern, otelhttp.NewHandler(h, pattern))
}

// ServeHTTP implements http.Handler, applying global middleware.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var h http.Handler = r.ServeMux
	for _, mw := range slices.Backward(r.globalChain) {
		h = mw(h)
	}

	h.ServeHTTP(w, req)
}

// ServerHandlerInstrumented returns an instrumented handler for use as http.Server.Handler.
func (r *Router) ServerHandlerInstrumented() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		holder := &PatternHolder{}
		ctx := context.WithValue(req.Context(), PatternKey{}, holder)

		name := req.Method + " - " + ServerRootSpanName

		ctx, span := otel.Tracer(ServerRootSpanName).Start(ctx, name)
		defer span.End()

		r.ServeHTTP(w, req.WithContext(ctx))

		if holder.Pattern != "" {
			span.SetName(holder.Pattern + " - " + ServerRootSpanName)
		}
	})
}
