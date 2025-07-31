/*
Copyright 2025 kemadev
SPDX-License-Identifier: MPL-2.0
*/

package main

import (
	"fmt"
	"net/http"

	"github.com/kemadev/go-framework/pkg/kctx"
	"github.com/kemadev/go-framework/pkg/router"
	"github.com/kemadev/go-framework/pkg/server"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	app := router.New()

	// Add middlewares
	app.UseInstrumented("kctx-middleware", kctx.Middleware)
	app.UseInstrumented("logging-middleware", LoggingMiddleware)

	// Create groups
	app.Group(func(r *router.Router) {
		r.UseInstrumented("auth-middleware", AuthMiddleware)

		r.Group(func(r *router.Router) {
			r.UseInstrumented("timing-middleware", TimingMiddleware)

			r.HandleInstrumented(
				"GET /auth/2/{bar}",
				http.HandlerFunc(FooBar),
			)
		})
	})

	// Add handlers
	app.HandleInstrumented(
		"GET /foo/{bar}",
		http.HandlerFunc(FooBar),
	)

	server.Run(app.ServerHandlerInstrumented())
}

func FooBar(w http.ResponseWriter, r *http.Request) {
	c := kctx.FromRequest(r)

	// Get user from context (set by AuthMiddleware)
	user := c.Local("user")

	// Get span context for logging
	span := c.Span(r)
	spanCtx := span.SpanContext()
	fmt.Printf("[HANDLER] TraceID: %s, SpanID: %s, User: %v\n",
		spanCtx.TraceID().String(),
		spanCtx.SpanID().String(),
		user,
	)

	bag := c.Baggage(r)
	span.AddEvent(
		"handling this...",
		trace.WithAttributes(semconv.UserID(bag.Member(string(semconv.UserIDKey)).Value())),
	)

	span.SetAttributes(attribute.String("bar", r.PathValue("bar")))

	fmt.Fprintf(w, "Hello, %v! TraceID: %s", user, spanCtx.TraceID().String())
}
