package otel

import (
	"net/http"

	"github.com/kemadev/go-framework/pkg/router"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

type PatternHolder struct {
	Pattern string
}

type PatternKey struct{}

type Middleware func(http.Handler) http.Handler

const packageName = "github.com/kemadev/go-framework/pkg/convenience/otel"

// WrapMux wraps a router with OpenTelemetry HTTP instrumentation.
func WrapMux(mux *router.Router, packageName string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		otelhttp.NewHandler(
			mux,
			packageName,
			otelhttp.WithSpanNameFormatter(
				func(operation string, r *http.Request) string {
					pattern := r.Pattern
					if pattern != "" {
						return pattern + " (mux)"
					}
					return operation
				},
			),
		).ServeHTTP(w, r)
	})
}

// WrapMux wraps a handler with an OpenTelemetry span.
func WrapHandler(
	pattern string,
	handler func(w http.ResponseWriter, r *http.Request),
) (string, http.HandlerFunc) {
	return pattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, span := otel.Tracer(pattern).
			Start(r.Context(), pattern)
		defer span.End()

		handler(w, r.WithContext(c))
	})
}

func WrapMiddleware(
	spanName string,
	middleware Middleware,
) Middleware {
	tracer := otel.Tracer(packageName)

	return func(next http.Handler) http.Handler {
		wrappedHandler := middleware(next)

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, span := tracer.Start(r.Context(), spanName)
			defer span.End()

			wrappedHandler.ServeHTTP(w, r.WithContext(c))
		})
	}
}
