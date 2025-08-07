// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package router

import (
	"embed"
	"net/http"
	"slices"

	"github.com/kemadev/go-framework/pkg/convenience/render"
)

const ServerRootSpanName = "server"

type PatternHolder struct {
	Pattern string
}
type PatternKey struct{}

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

// Handle registers a handler for a pattern.
func (r *Router) Handle(pattern string, h http.Handler) {
	for _, mw := range slices.Backward(r.routeChain) {
		h = mw(h)
	}

	r.ServeMux.Handle(pattern, h)
}

// HandleStatic registers a handler for a pattern, serving static files from [fs].
func (r *Router) HandleStatic(pattern string, fs embed.FS) {
	handler := http.FileServer(http.FS(fs))
	r.Handle(pattern, handler)
}

// HandleHtmlTmpl registers a handler for a pattern, executing [html/template.Template] files from [fs], using [renderer] with
// data from [dataFunc], using template nammed after [net/http.Request.PathValue] of [pathValueKey] suffixed by `.html`.
func (r *Router) HandleHtmlTmpl(
	pattern string,
	fs embed.FS,
	renderer *render.TemplateRenderer,
	dataFunc func(r *http.Request) (any, error),
	pathValueKey string,
) {
	handler := renderer.HandlerFuncWithData(
		dataFunc,
		pathValueKey,
		".html",
	)
	r.Handle(pattern, handler)
}

// ServeHTTP implements http.Handler, applying global middleware.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var h http.Handler = r.ServeMux
	for _, mw := range slices.Backward(r.globalChain) {
		h = mw(h)
	}

	h.ServeHTTP(w, req)
}
