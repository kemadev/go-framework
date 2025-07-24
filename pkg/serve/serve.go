package serve

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/kemadev/go-framework/pkg/config"
	"github.com/kemadev/go-framework/pkg/monitoring"
	"github.com/kemadev/go-framework/pkg/otel"
	"github.com/kemadev/go-framework/pkg/route"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// Run starts the HTTP server and registers the routes.
// It handles the shutdown (can be caused by SIGINT) gracefully and sets up OpenTelemetry.
// It returns an error if the server fails to start or if the shutdown fails.
// It should be called from the main function of the application.
// It is a blocking call and will not return until the server is shut down.
func Run(
	routes route.RoutesToRegister,
	dependencyRoutes route.RoutesWithDependencies,
	conf *config.Config,
) error {
	// Handle SIGINT gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Set up OpenTelemetry.
	otelShutdown, err := otel.SetupOTelSDK(ctx, *conf)
	if err != nil {
		return fmt.Errorf("otel setup: %w", err)
	}
	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	// Create instrumented slog logger and set it as default
	logger := otelslog.NewLogger(
		"default",
		otelslog.WithSource(true),
		otelslog.WithVersion(conf.AppVersion),
	)
	slog.SetDefault(logger)

	// Add monitoring routes
	dependencyRoutes = append(dependencyRoutes, monitoring.Routes()...)

	server := route.NewServer(conf)

	// Start HTTP server.
	srv := &http.Server{
		// Use any host, let Kubernetes handle the routing.
		Addr:         ":" + strconv.Itoa(conf.HTTPServePort),
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		ReadTimeout:  time.Duration(conf.HTTPReadTimeout) * time.Second,
		WriteTimeout: time.Duration(conf.HTTPWriteTimeout) * time.Second,
		Handler:      newHTTPHandler(server, routes, dependencyRoutes),
	}
	srvErr := make(chan error, 1)

	go func() {
		srvErr <- srv.ListenAndServe()
	}()

	// Wait for interruption.
	select {
	case err = <-srvErr:
		// Error when starting HTTP server.
		return fmt.Errorf("http server error: %w", err)
	case <-ctx.Done():
		// Wait for first SIGINT
		// Stop receiving signal notifications as soon as possible.
		stop()
	}

	// When Shutdown is called, ListenAndServe immediately returns ErrServerClosed.
	err = srv.Shutdown(context.Background())
	if err != nil {
		// Error when shutting down HTTP server.
		return fmt.Errorf("http server shutdown error: %w", err)
	}

	return nil
}

// newHTTPHandler returns a new HTTP handler with the given routes.
// It uses the [net/http.ServeMux] to register the routes and adds OpenTelemetry instrumentation.
func newHTTPHandler(
	server *route.Server,
	routes route.RoutesToRegister,
	dependencyRoutes route.RoutesWithDependencies,
) http.Handler {
	mux := http.NewServeMux()

	// handleFunc is a replacement for [net/http.mux.HandleFunc] which enriches the
	// handler's HTTP instrumentation with the pattern of the [net/http.route].
	handleFunc := func(pattern string, handlerFunc func(http.ResponseWriter, *http.Request)) {
		// Configure the [net/http.route] with automatic instrumentation.
		handler := otelhttp.WithRouteTag(pattern, http.HandlerFunc(handlerFunc))
		mux.Handle(pattern, handler)
	}

	// Register regular handlers (without dependency injection).
	for _, route := range routes {
		handleFunc(route.Pattern, route.HandlerFunc)
	}

	// Register handlers with dependency injection.
	for _, route := range dependencyRoutes {
		handlerWithDeps := route.HandlerFunc(server)
		handleFunc(route.Pattern, handlerWithDeps)
	}

	// Add HTTP instrumentation for the whole server.
	handler := otelhttp.NewHandler(mux, "/")

	return handler
}
