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
	"go.opentelemetry.io/otel/trace"
)

func main() {
	r := router.New()

	// Add middlewares
	r.UseInstrumented("logging-middleware", LoggingMiddleware)
	r.UseInstrumented("auth-middleware", AuthMiddleware)
	r.UseInstrumented("timing-middleware", TimingMiddleware)

	// Default handler when nothing matches
	r.HandleInstrumented("/", http.NotFoundHandler())

	r.HandleInstrumented("GET /foo", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get user from context (set by AuthMiddleware)
		user := ctx.Value("user")

		// Get span context for logging
		spanCtx := trace.SpanFromContext(ctx).SpanContext()
		fmt.Printf("[HANDLER] TraceID: %s, SpanID: %s, User: %v\n",
			spanCtx.TraceID().String(),
			spanCtx.SpanID().String(),
			user,
		)

		w.Write([]byte(fmt.Sprintf("Hello, %v! TraceID: %s", user, spanCtx.TraceID().String())))
	}))

	server.Run(r)
}
