/*
Copyright 2025 kemadev
SPDX-License-Identifier: MPL-2.0
*/

package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/kemadev/go-framework/pkg/convenience/headutil"
	"github.com/kemadev/go-framework/pkg/convenience/headval"
	"github.com/kemadev/go-framework/pkg/convenience/local"
	"github.com/kemadev/go-framework/pkg/convenience/resp"
	"github.com/kemadev/go-framework/pkg/convenience/trace"
	"github.com/kemadev/go-framework/pkg/monitoring"
	"github.com/kemadev/go-framework/pkg/router"
	"github.com/kemadev/go-framework/pkg/server"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
	oteltrace "go.opentelemetry.io/otel/trace"
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
	app.UseInstrumented("logging-middleware", LoggingMiddleware)

	// Create groups
	app.Group(func(r *router.Router) {
		r.UseInstrumented("auth-middleware", AuthMiddleware)

		r.Group(func(r *router.Router) {
			r.UseInstrumented("timing-middleware", TimingMiddleware)

			r.HandleInstrumented(
				"GET /auth/{bar}",
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
	// Get user from context (set by AuthMiddleware)
	user := local.Get(r.Context(), "user")

	// Get span context for logging
	span := trace.Span(r.Context())
	spanCtx := span.SpanContext()
	fmt.Printf("[HANDLER] TraceID: %s, SpanID: %s, User: %v\n",
		spanCtx.TraceID().String(),
		spanCtx.SpanID().String(),
		user,
	)

	bag := trace.Baggage(r.Context())
	span.AddEvent(
		"handling this...",
		oteltrace.WithAttributes(semconv.UserID(bag.Member(string(semconv.UserIDKey)).Value())),
	)

	span.SetAttributes(attribute.String("bar", r.PathValue("bar")))

	fmt.Println(headutil.IsMIME(r.Header, headval.AcceptJSON))

	fmt.Fprintf(w, "Hello, %v! TraceID: %s", user, spanCtx.TraceID().String())
}

func Redir(w http.ResponseWriter, r *http.Request) {
	resp.Redirect(w, http.StatusPermanentRedirect, url.URL{
		Scheme: "https",
		Host:   "google.com",
	})
}
