package otel

import (
	"context"
	"net/http"

	"github.com/kemadev/go-framework/pkg/router"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

const ServerRootSpanName = "server"

type PatternHolder struct {
	Pattern string
}
type PatternKey struct{}

func WrapMux(mux *router.Router, packageName string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		holder := &PatternHolder{}
		c := context.WithValue(r.Context(), PatternKey{}, holder)

		c, span := otel.Tracer(packageName).
			Start(c, packageName)
		defer span.End()

		mux.ServeHTTP(w, r.WithContext(c))

		holder, ok := c.Value(PatternKey{}).(*PatternHolder)
		if ok {
			pattern := holder.Pattern
			if pattern == "" {
				span.SetName(packageName)
				return
			}
			span.SetName(pattern + " - " + ServerRootSpanName)
		}
	})
}

func WrapHandler(
	pattern string,
	handler func(w http.ResponseWriter, r *http.Request),
) (string, http.HandlerFunc) {
	return pattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		holder, ok := r.Context().Value(PatternKey{}).(*PatternHolder)
		if ok {
			holder.Pattern = pattern
		}

		h := otelhttp.NewHandler(http.HandlerFunc(handler), pattern)

		h.ServeHTTP(w, r)
	})
}
