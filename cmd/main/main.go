/*
Copyright 2025 kemadev
SPDX-License-Identifier: MPL-2.0
*/

package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/kemadev/go-framework/pkg/header"
	"github.com/kemadev/go-framework/pkg/kctx"
	"github.com/kemadev/go-framework/pkg/monitoring"
	"github.com/kemadev/go-framework/pkg/router"
	"github.com/kemadev/go-framework/pkg/server"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
	"go.opentelemetry.io/otel/trace"
)

const packageName = "github.com/kemadev/go-framework/cmd/main"

func main() {
	app := router.New()

	app.HandleInstrumented(
		monitoring.LivenessHandler(
			func() monitoring.CheckResults { return monitoring.CheckResults{} },
		),
	)
	app.HandleInstrumented(
		monitoring.ReadinessHandler(
			func() monitoring.CheckResults { return monitoring.CheckResults{} },
		),
	)

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
	app.HandleInstrumented(
		"GET /redir",
		http.HandlerFunc(Redir),
	)

	server.Run(app.ServerHandlerInstrumented())
}

func FooBar(w http.ResponseWriter, r *http.Request) {
	c := kctx.FromRequestWarn(r, packageName)

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

	fmt.Println(c.IsMIME(header.ValueAcceptJSON))

	fmt.Fprintf(w, "Hello, %v! TraceID: %s", user, spanCtx.TraceID().String())
}

func Redir(w http.ResponseWriter, r *http.Request) {
	c := kctx.FromRequestWarn(r, packageName)

	c.Redirect(http.StatusPermanentRedirect, url.URL{
		Scheme: "https",
		Host:   "google.com",
	})
}
