package route

import (
	"net/http"

	"github.com/kemadev/go-framework/pkg/config"
)

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

type Server struct {
	config *config.Config
}

// NewServer creates a new server with dependencies
func NewServer(cfg *config.Config) *Server {
	return &Server{
		config: cfg,
	}
}

// GetConfig returns the server's config
func (s *Server) GetConfig() *config.Config {
	return s.config
}

// RoutesToRegister is a slice of HTTPRoute.
// It should be used as a convenience type to pass a list of routes to the HTTP server.
type RoutesToRegister []Route

// RouteWithDependencies represents a route that needs access to server dependencies
type RouteWithDependencies struct {
	Pattern     string
	HandlerFunc func(*Server) http.HandlerFunc
}

// RoutesWithDependencies is a slice of RouteWithDependencies
type RoutesWithDependencies []RouteWithDependencies

// CreateRoute creates a RouteWithDependencies from a pattern and handler function
func CreateRoute(pattern string, handlerFunc func(*Server) http.HandlerFunc) RouteWithDependencies {
	return RouteWithDependencies{
		Pattern:     pattern,
		HandlerFunc: handlerFunc,
	}
}

// ServerHandler is a convenience type for handlers that need server dependencies
type ServerHandler func(*Server, http.ResponseWriter, *http.Request)

// CreateRouteWithHandler creates a RouteWithDependencies from a pattern and ServerHandler
func CreateRouteWithHandler(pattern string, handler ServerHandler) RouteWithDependencies {
	return RouteWithDependencies{
		Pattern: pattern,
		HandlerFunc: func(server *Server) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				handler(server, w, r)
			}
		},
	}
}
