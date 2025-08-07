/*
Copyright 2025 kemadev
SPDX-License-Identifier: MPL-2.0
*/

package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/kemadev/go-framework/pkg/convenience/local"
	"github.com/kemadev/go-framework/pkg/convenience/otel"
	"github.com/kemadev/go-framework/pkg/convenience/render"
	"github.com/kemadev/go-framework/pkg/convenience/req"
	"github.com/kemadev/go-framework/pkg/convenience/trace"
	"github.com/kemadev/go-framework/pkg/monitoring"
	"github.com/kemadev/go-framework/pkg/router"
	"github.com/kemadev/go-framework/pkg/server"
	"github.com/kemadev/go-framework/web"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

const packageName = "github.com/kemadev/go-framework/cmd/main"

func main() {
	r := router.New()

	// Add middlewares
	r.Use(otel.WrapMiddleware("logging", LoggingMiddleware))

	// Add monitoring endpoints
	r.Handle(
		monitoring.LivenessHandler(
			func() monitoring.CheckResults { return monitoring.CheckResults{} },
		),
	)
	r.Handle(
		monitoring.ReadinessHandler(
			func() monitoring.CheckResults { return monitoring.CheckResults{} },
		),
	)

	// Add handlers
	r.Handle(
		otel.WrapHandler("GET /foo/{bar}", http.HandlerFunc(FooBar)),
	)
	r.Handle(
		otel.WrapHandler(
			"GET /tester",
			http.HandlerFunc(Tester),
		),
	)

	// Create groups
	r.Group(func(r *router.Router) {
		r.Use(otel.WrapMiddleware("auth", AuthMiddleware))

		r.Group(func(r *router.Router) {
			r.Use(otel.WrapMiddleware("timing", TimingMiddleware))

			r.Handle(
				otel.WrapHandler(
					"GET /auth/{bar}",
					http.HandlerFunc(FooBar),
				),
			)
		})
	})

	// Handle static assets
	r.Handle(
		otel.WrapHandler(
			"GET /static/",
			http.FileServerFS(web.GetStaticFS()).ServeHTTP,
		),
	)

	// Handle template assets
	tmplFS := web.GetTmplFS()
	renderer, _ := render.New(tmplFS)
	r.Handle(
		otel.WrapHandler(
			"GET /tmpl/",
			renderer.HandlerFuncWithData(
				func(r *http.Request) (any, error) {
					return map[string]any{
						"WorldName": "WoRlD",
					}, nil
				},
			),
		),
	)

	server.Run(otel.WrapMux(r, packageName))
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

	fmt.Fprintf(w, "Hello, %v! TraceID: %s", user, spanCtx.TraceID().String())
}

type TesterPayload2 struct {
	Foo    string
	hidden string
}

type TesterPayload struct {
	hidden string
	Foo    string `json:"hello"`
	Quux   TesterPayload2
}

func Tester(w http.ResponseWriter, r *http.Request) {
	ip, err := req.Hosts(r)
	if err != nil {
		slog.Debug(err.Error())
		return
	}
	slog.Debug(fmt.Sprintf("%s", ip))
}
