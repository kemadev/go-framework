// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"net/http"

	"github.com/kemadev/go-framework/pkg/kctx"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
	"go.opentelemetry.io/otel/trace"
)

// LoggingMiddleware logs request details.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := kctx.FromRequest(r)
		if c == nil {
			c.Logger(packageName).Warn("failure finding kctx in request context")
			c = &kctx.Kctx{}
		}

		spanCtx := trace.SpanFromContext(c).SpanContext()
		fmt.Printf("[LOGGING] TraceID: %s, SpanID: %s, Method: %s, Path: %s\n",
			spanCtx.TraceID().String(),
			spanCtx.SpanID().String(),
			r.Method,
			r.URL.Path,
		)

		next.ServeHTTP(w, r)
	})
}

// AuthMiddleware simulates authentication.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := kctx.FromRequest(r)
		if c == nil {
			c.Logger(packageName).Warn("failure finding kctx in request context")
			c = &kctx.Kctx{}
		}

		userID := "whoever"
		span := c.Span(r)
		span.SetAttributes(
			semconv.UserID(userID),
			attribute.Bool("auth.authenticated", true),
		)

		mem, err := baggage.NewMember(string(semconv.UserIDKey), userID)
		if err != nil {
			fmt.Printf("Failed to create baggage member: %v\n", err)
			next.ServeHTTP(w, r)

			return
		}

		newCtx, err := c.BaggageSet(r, mem)
		if err != nil {
			fmt.Printf("Failed to set baggage: %v\n", err)
			next.ServeHTTP(w, r)

			return
		}

		c.Context = newCtx

		spanCtx := span.SpanContext()
		fmt.Printf("[AUTH] TraceID: %s, SpanID: %s, User: "+userID+"\n",
			spanCtx.TraceID().String(),
			spanCtx.SpanID().String(),
		)

		c.LocalSet("user", userID)

		next.ServeHTTP(w, r.WithContext(c.Context))
	})
}

// TimingMiddleware logs timing.
func TimingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := kctx.FromRequest(r)
		if c == nil {
			c.Logger(packageName).Warn("failure finding kctx in request context")
			c = &kctx.Kctx{}
		}

		spanCtx := trace.SpanFromContext(c).SpanContext()
		fmt.Printf("[TIMING] TraceID: %s, SpanID: %s, Starting request\n",
			spanCtx.TraceID().String(),
			spanCtx.SpanID().String(),
		)

		next.ServeHTTP(w, r)

		fmt.Printf("[TIMING] TraceID: %s, SpanID: %s, Request completed\n",
			spanCtx.TraceID().String(),
			spanCtx.SpanID().String(),
		)
	})
}
