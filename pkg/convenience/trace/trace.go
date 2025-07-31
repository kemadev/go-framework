// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package trace

import (
	"context"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/trace"
)

// Span return the span from kctx.
func Span(r *http.Request) trace.Span {
	return trace.SpanFromContext(r.Context())
}

// SpanCtx return the span context from kctx. If you already have a reference to the span, prefer
// using [go.opentelemetry.io/otel/trace].Span.SpanContext().
func SpanCtx(r *http.Request) trace.SpanContext {
	return Span(r).SpanContext()
}

// SpanSetAttrs sets attributes for a span. If you already have a reference to the span, prefer
// using [go.opentelemetry.io/otel/trace].Span.SetAttributes().
func SpanSetAttrs(r *http.Request, kv ...attribute.KeyValue) {
	span := Span(r)
	span.SetAttributes(kv...)
}

// Baggage return the baggage from kctx.
func Baggage(r *http.Request) baggage.Baggage {
	return baggage.FromContext(r.Context())
}

// BaggageSetMembers sets baggage members for a span. If you already have a reference to the baggage, prefer
// using [go.opentelemetry.io/otel/baggage].Baggage.SetMember()
// Please not that returned context needs to be propagated in order for the baggage to be propagated, too.
func BaggageSet(r *http.Request, members ...baggage.Member) (context.Context, error) {
	bag := baggage.FromContext(r.Context())

	var err error

	for _, member := range members {
		bag, err = bag.SetMember(member)
		if err != nil {
			return nil, fmt.Errorf("error setting baggage member: %w", err)
		}
	}

	return baggage.ContextWithBaggage(r.Context(), bag), nil
}
