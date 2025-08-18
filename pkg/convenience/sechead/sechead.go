package sechead

import (
	"log/slog"
	"net/http"

	"github.com/kemadev/go-framework/pkg/convenience/log"
)

const packageName = "github.com/kemadev/go-framework/pkg/convenience/sechead"

// NewMiddleware returns a middleware adding security headers.
func NewMiddleware(conf SecurityHeadersConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for key, val := range conf.Headers() {
				if len(val) != 1 {
					log.GetPackageLogger(packageName).Error("multiple values found in header", slog.String("header", key))
				}

				w.Header().Add(key, val[0])
			}

			next.ServeHTTP(w, r)
		})
	}
}
