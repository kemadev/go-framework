package timeout

import (
	"net/http"
	"time"
)

// WrapHandler returns an handler wrapping [handler] with a timeout set to [timeout].
func WrapHandler(h http.Handler, t time.Duration) http.Handler {
	return http.TimeoutHandler(h, t, http.StatusText(http.StatusServiceUnavailable))
}

// NewMiddleware returns an middleware with a timeout set to [timeout].
func NewMiddleware(t time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return WrapHandler(next, t)
	}
}
