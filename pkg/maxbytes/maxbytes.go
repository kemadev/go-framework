package maxbytes

import (
	"net/http"
)

// WrapHandler returns sets max bytes for request body, see [net/http.MaxBytesHandler].
func WrapHandler(h http.Handler, n int64) http.Handler {
	return http.MaxBytesHandler(h, n)
}

// NewMiddleware returns n middleware with max bytes for request body, see [net/http.MaxBytesHandler].
func NewMiddleware(n int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return WrapHandler(next, n)
	}
}
