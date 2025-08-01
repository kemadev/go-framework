// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

//
// [this nice article]: https://www.alexedwards.net/blog/organize-your-go-middleware-without-dependencies
// [Alex Edwards]: https://www.alexedwards.net/
package router

import (
	"net/http"
	"slices"
)

const ServerRootSpanName = "server"

type PatternHolder struct {
	Pattern string
}
type PatternKey struct{}

// Router is an HTTP router.
type Router struct {
	routeChain []func(http.Handler) http.Handler
	*http.ServeMux
}

// New returns a new HTTP router.
func New() *Router {
	return &Router{ServeMux: http.NewServeMux()}
}

// Use appends [mw] to the routers chain.
func (r *Router) Use(mw ...func(http.Handler) http.Handler) {
	r.routeChain = append(r.routeChain, mw...)
}

// Group adds all routers down the chain to a group. All members of a group inherits from
// their parent's routers chain.
func (r *Router) Group(group func(r *Router)) {
	subRouter := &Router{
		routeChain: slices.Clone(r.routeChain),
		ServeMux:   r.ServeMux,
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

// ServeHTTP implements http.Handler, applying global middleware.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var h http.Handler = r.ServeMux
	for _, mw := range slices.Backward(r.routeChain) {
		h = mw(h)
	}

	h.ServeHTTP(w, req)
}
