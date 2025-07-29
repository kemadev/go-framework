package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
	"go.opentelemetry.io/otel/trace"
)

// LoggingMiddleware logs request details
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		spanCtx := trace.SpanFromContext(ctx).SpanContext()
		fmt.Printf("[MIDDLEWARE] TraceID: %s, SpanID: %s, Method: %s, Path: %s\n",
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
		ctx := r.Context()

		userId := "whoever"
		span := trace.SpanFromContext(ctx)
		span.SetAttributes(
			semconv.UserID(userId),
			attribute.Bool("auth.authenticated", true),
		)
		bag := baggage.FromContext(ctx)
		mem, err := baggage.NewMember(string(semconv.UserIDKey), userId)
		if err != nil {
			fmt.Println("failed")
		}
		bag, err = bag.SetMember(mem)
		if err != nil {
			fmt.Println("failed")
		}

		spanCtx := span.SpanContext()
		fmt.Printf("[AUTH] TraceID: %s, SpanID: %s, User: kema-dev\n",
			spanCtx.TraceID().String(),
			spanCtx.SpanID().String(),
		)

		ctx = context.WithValue(ctx, "user", "kema-dev")
		ctx = baggage.ContextWithBaggage(ctx, bag)

		slog.Error("ff?")

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// TimingMiddleware logs timing
func TimingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

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
