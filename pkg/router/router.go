package router

import (
	"net/http"
	"slices"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

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

// Router is an HTTP router
type Router struct {
	globalChain []func(http.Handler) http.Handler
	routeChain  []func(http.Handler) http.Handler
	isSubRouter bool
	*http.ServeMux
}

// New returns a new HTTP router
func New() *Router {
	return &Router{ServeMux: http.NewServeMux()}
}

// Use appends [mw] to the routers chain
func (r *Router) Use(mw ...func(http.Handler) http.Handler) {
	if r.isSubRouter {
		r.routeChain = append(r.routeChain, mw...)
	} else {
		r.globalChain = append(r.globalChain, mw...)
	}
}

// Group add all routers down the chain to a group. All members of a group inherits from
// their parent's routers chain.
func (r *Router) Group(fn func(r *Router)) {
	subRouter := &Router{
		routeChain:  slices.Clone(r.routeChain),
		isSubRouter: true,
		ServeMux:    r.ServeMux,
	}
	fn(subRouter)
}

// HandleFunc returns a func satisfying [net/http.HandleFunc]
func (r *Router) HandleFunc(pattern string, h http.HandlerFunc) {
	r.Handle(pattern, h)
}

// HandleFunc returns a func satisfying [net/http.HandleFunc], wrapping the handler with OpenTelemetry instrumentation
func (r *Router) HandleFuncOTEL(pattern string, h http.HandlerFunc) {
	r.HandleOTEL(pattern, h)
}

// HandleFunc returns a func satisfying [net/http.Handle]
func (r *Router) Handle(pattern string, h http.Handler) {
	for _, mw := range slices.Backward(r.routeChain) {
		h = mw(h)
	}
	r.ServeMux.Handle(pattern, h)
}

// HandleFunc returns a func satisfying [net/http.Handle], wrapping the handler with OpenTelemetry instrumentation
func (r *Router) HandleOTEL(pattern string, h http.Handler) {
	for _, mw := range slices.Backward(r.routeChain) {
		h = mw(h)
	}
	r.ServeMux.Handle(pattern, otelhttp.NewHandler(h, pattern))
}

// HandleFunc returns a func satisfying [net/http.Handler.ServeHTTP]
func (r *Router) ServeHTTP(w http.ResponseWriter, rq *http.Request) {
	var h http.Handler = r.ServeMux

	for _, mw := range slices.Backward(r.globalChain) {
		h = mw(h)
	}
	h.ServeHTTP(w, rq)
}
