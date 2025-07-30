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

// LoggingMiddleware logs request details
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := kctx.FromRequest(r)

		spanCtx := trace.SpanFromContext(ctx).SpanContext()
		fmt.Printf("[LOGGING] TraceID: %s, SpanID: %s, Method: %s, Path: %s\n",
			spanCtx.TraceID().String(),
			spanCtx.SpanID().String(),
			r.Method,
			r.URL.Path,
		)

		next.ServeHTTP(w, r)
	})
}

// AuthMiddleware simulates authentication
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := kctx.FromRequest(r)

		userId := "whoever"
		span := ctx.Span(r)
		span.SetAttributes(
			semconv.UserID(userId),
			attribute.Bool("auth.authenticated", true),
		)

		mem, err := baggage.NewMember(string(semconv.UserIDKey), userId)
		if err != nil {
			fmt.Printf("Failed to create baggage member: %v\n", err)
			next.ServeHTTP(w, r)
			return
		}

		err = ctx.BaggageSet(r, mem)
		if err != nil {
			fmt.Printf("Failed to set baggage: %v\n", err)
			next.ServeHTTP(w, r)
			return
		}

		spanCtx := span.SpanContext()
		fmt.Printf("[AUTH] TraceID: %s, SpanID: %s, User: kema-dev\n",
			spanCtx.TraceID().String(),
			spanCtx.SpanID().String(),
		)

		ctx.LocalSet("user", "kema-dev")

		next.ServeHTTP(w, r)
	})
}

// TimingMiddleware logs timing
func TimingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := kctx.FromRequest(r)

		spanCtx := trace.SpanFromContext(ctx).SpanContext()
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
