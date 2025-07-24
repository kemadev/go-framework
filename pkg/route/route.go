package route

import "net/http"

// Route contains the pattern and the handler function for an HTTP route.
// The attributes are passed to [net/http.ServeMux.Handle], see its package's documentation for more information.
type Route struct {
	// Pattern is the pattern for the HTTP route.
	// See [net/http.ServeMux.Handle] for more information on how to use it.
	Pattern string
	// HandlerFunc is the handler function for the HTTP route.
	// See [net/http.ServeMux.Handle] for more information on how to use it.
	HandlerFunc func(http.ResponseWriter, *http.Request)
}
