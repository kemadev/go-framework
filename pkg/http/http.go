// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

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
