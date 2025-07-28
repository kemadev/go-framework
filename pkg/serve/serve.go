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
	"syscall"
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
	conf *config.Global,
) error {
	// Intercept signals
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
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

	// Add monitoring routes
	dependencyRoutes = append(dependencyRoutes, monitoring.Routes()...)

	server := route.NewServer(conf)

	// Start HTTP server.
	srv := &http.Server{
		// Use any host, let Kubernetes handle the routing.
		Addr:         ":" + strconv.Itoa(conf.Server.BindPort),
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		ReadTimeout:  time.Duration(conf.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(conf.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(conf.Server.IdleTimeout) * time.Second,
		ErrorLog: slog.NewLogLogger(
			otelslog.NewLogger("httpserver").Handler(),
			slog.LevelError,
		),
		Handler: newHTTPHandler(server, routes, dependencyRoutes),
	}

	srvErr := make(chan error, 1)

	go func() {
		srvErr <- srv.ListenAndServe()
	}()

	// Wait for interruption
	select {
	case err = <-srvErr:
		return fmt.Errorf("http server error: %w", err)
	case <-ctx.Done():
		// Stop receiving signal notifications as soon as possible.
		stop()
	}

	err = srv.Shutdown(context.Background())
	if err != nil {
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
		handler := otelhttp.NewHandler(http.HandlerFunc(handlerFunc), pattern)
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

	// Return 404 by default, still instrument the route
	mux.Handle("/", otelhttp.NewHandler(http.NotFoundHandler(), "/"))

	return mux
}
