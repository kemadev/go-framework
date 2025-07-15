// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package http

import (
	"context"
	"encoding/json"
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
	"github.com/kemadev/go-framework/pkg/otel"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

var (
	// ErrJSONEncodingFail is a sentinel error that indicates that JSON encoding failed.
	ErrJSONEncodingFail = fmt.Errorf("failed to encode JSON response")
	// ErrInternalServerError is a sentinel error that indicates an internal server error.
	// It uses conventional HTTP status text for internal server error.
	ErrInternalServerError = fmt.Errorf("%s", http.StatusText(http.StatusInternalServerError))
)

const (
	// ContentTypeHeaderKey is the HTTP header key for the content type.
	ContentTypeHeaderKey = "Content-Type"
	// ContentTypeJSON is the HTTP content type for JSON.
	ContentTypeJSON = "application/json"
	// ContentTypePlain is the HTTP content type for plain text.
	ContentTypePlain = "text/plain"
)

// ClientInfo contains information about an HTTP client request, as well as the
// required instrumentation to log and trace it.
type ClientInfo struct {
	// Ctx is the context of the HTTP request.
	// It is used to propagate the context across the request.
	// It is also used to cancel the request if needed.
	Ctx context.Context
	// Writer is the HTTP response writer.
	// It is used to write the response to the client.
	Writer http.ResponseWriter
	// Logger is the logger used to log the request and response.
	// It should be an instrumented logger.
	Logger *slog.Logger
	// Span is the trace span used to trace the request.
	// It should be an instrumented span.
	Span trace.Span
}

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

// RoutesToRegister is a slice of HTTPRoute.
// It should be used as a convinience type to pass a list of routes to the HTTP server.
type RoutesToRegister []Route

// Run starts the HTTP server and registers the routes.
// It handles the shutdown (can be caused by SIGINT) gracefully and sets up OpenTelemetry.
// It returns an error if the server fails to start or if the shutdown fails.
// It should be called from the main function of the application.
// It is a blocking call and will not return until the server is shut down.
func Run(routes RoutesToRegister, conf *config.Config) error {
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

	// Start HTTP server.
	srv := &http.Server{
		// Use any host, let Kubernetes handle the routing.
		Addr:         ":" + strconv.Itoa(conf.HTTPServePort),
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		ReadTimeout:  time.Duration(conf.HTTPReadTimeout) * time.Second,
		WriteTimeout: time.Duration(conf.HTTPWriteTimeout) * time.Second,
		Handler:      newHTTPHandler(routes),
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
func newHTTPHandler(routes RoutesToRegister) http.Handler {
	mux := http.NewServeMux()

	// handleFunc is a replacement for [net/http.mux.HandleFunc] which enriches the
	// handler's HTTP instrumentation with the pattern as the [net/http.route].
	handleFunc := func(pattern string, handlerFunc func(http.ResponseWriter, *http.Request)) {
		// Configure the [net/http.route] with automatic instrumentation.
		handler := otelhttp.WithRouteTag(pattern, http.HandlerFunc(handlerFunc))
		mux.Handle(pattern, handler)
	}

	// Register handlers.
	for _, route := range routes {
		handleFunc(route.Pattern, route.HandlerFunc)
	}

	// Add HTTP instrumentation for the whole server.
	handler := otelhttp.NewHandler(mux, "/")

	return handler
}

// SendJSONResponse sends a JSON response to the client.
// It sets the content type to JSON and writes the response with the given status code.
// It uses [encoding/json.Encoder] to encode the data to JSON.
// If the JSON encoding fails, it sends an internal server error response.
// It should be used to send JSON responses to the client.
// No further writing should be done after calling this function.
func SendJSONResponse(
	clientInfo ClientInfo,
	statusCode int,
	data any,
) {
	clientInfo.Writer.Header().Set(ContentTypeHeaderKey, ContentTypeJSON)
	clientInfo.Writer.WriteHeader(statusCode)

	err := json.NewEncoder(clientInfo.Writer).Encode(data)
	if err != nil {
		SendErrorResponse(
			clientInfo,
			http.StatusInternalServerError,
			ErrInternalServerError,
			err,
		)
	}
}

// ErrToSend is sent to client but not logged,
// errToLog is logged but not sent to client. Those two can be identical or different,
// allowing to send a different error message to the client than the one logged, thus
// allowing to hide information from the client while logging it.
func SendErrorResponse(
	clientInfo ClientInfo,
	statusCode int,
	errToSend error,
	errToLog error,
) {
	errMap := map[string]string{
		"error": errToSend.Error(),
	}

	clientInfo.Writer.Header().Set(ContentTypeHeaderKey, ContentTypeJSON)
	clientInfo.Writer.WriteHeader(statusCode)
	InstrumentError(clientInfo.Ctx, clientInfo.Logger, clientInfo.Span, errToLog)

	err := json.NewEncoder(clientInfo.Writer).Encode(errMap)
	if err != nil {
		InstrumentError(clientInfo.Ctx, clientInfo.Logger, clientInfo.Span, err)
		http.Error(
			clientInfo.Writer,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
	}
}

// InstrumentError is a helper function that logs and traces an error.
// It should be used to log and trace errors that occur during the handling of an HTTP request.
func InstrumentError(
	ctx context.Context,
	logger *slog.Logger,
	span trace.Span,
	err error,
) {
	errString := err.Error()
	logger.ErrorContext(
		ctx,
		"error occurred",
		// TODO use semconv value once released, see https://opentelemetry.io/docs/specs/semconv/attributes-registry/error/#error-message
		slog.String("error.message", errString),
	)
	span.RecordError(err)
}

func SendResponse(
	clientInfo ClientInfo,
	statusCode int,
	contentType string,
	data []byte,
) {
	clientInfo.Writer.Header().Set(ContentTypeHeaderKey, contentType)
	clientInfo.Writer.WriteHeader(statusCode)

	_, err := clientInfo.Writer.Write(data)
	if err != nil {
		SendErrorResponse(
			clientInfo,
			http.StatusInternalServerError,
			ErrInternalServerError,
			err,
		)
	}
}
