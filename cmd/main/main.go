/*
Copyright 2025 kemadev
SPDX-License-Identifier: MPL-2.0
*/

package main

import (
	"fmt"
	"net/http"

	"github.com/kemadev/go-framework/pkg/router"
	"github.com/kemadev/go-framework/pkg/server"
	"go.opentelemetry.io/otel/baggage"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	r := router.New()

	// Add middlewares
	r.UseInstrumented("logging-middleware", LoggingMiddleware)
	r.UseInstrumented("auth-middleware", AuthMiddleware)
	r.UseInstrumented("timing-middleware", TimingMiddleware)

	r.HandleInstrumented(
		"GET /foo/{bar}",
		http.HandlerFunc(FooBar),
	)

	server.Run(r)
}

func FooBar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user from context (set by AuthMiddleware)
	user := ctx.Value("user")

	// Get span context for logging
	span := trace.SpanFromContext(ctx)
	spanCtx := span.SpanContext()
	fmt.Printf("[HANDLER] TraceID: %s, SpanID: %s, User: %v\n",
		spanCtx.TraceID().String(),
		spanCtx.SpanID().String(),
		user,
	)
	bag := baggage.FromContext(ctx)
	span.AddEvent(
		"handling this...",
		trace.WithAttributes(semconv.UserID(bag.Member(string(semconv.UserIDKey)).Value())),
	)

	w.Write([]byte(fmt.Sprintf("Hello, %v! TraceID: %s", user, spanCtx.TraceID().String())))
}
