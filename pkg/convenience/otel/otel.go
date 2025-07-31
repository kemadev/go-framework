package otel

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const ServerRootSpanName = "server"

type PatternHolder struct {
	Pattern string
}
type PatternKey struct{}

func WrapMux(mux *chi.Mux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		otelhttp.NewHandler(mux, ServerRootSpanName, otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string { return r.Pattern })).
			ServeHTTP(w, r)
	})
}

func WrapHandler(
	pattern string,
	handler func(w http.ResponseWriter, r *http.Request),
) (string, http.HandlerFunc) {
	return pattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	})
}
