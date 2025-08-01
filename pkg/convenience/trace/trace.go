// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package trace

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/trace"
)

// Span return the span from kctx.
func Span(c context.Context) trace.Span {
	return trace.SpanFromContext(c)
}

// SpanCtx return the span context from kctx. If you already have a reference to the span, prefer
// using [go.opentelemetry.io/otel/trace].Span.SpanContext().
func SpanCtx(c context.Context) trace.SpanContext {
	return Span(c).SpanContext()
}

// SpanSetAttrs sets attributes for a span. If you already have a reference to the span, prefer
// using [go.opentelemetry.io/otel/trace].Span.SetAttributes().
func SpanSetAttrs(c context.Context, kv ...attribute.KeyValue) {
	span := Span(c)
	span.SetAttributes(kv...)
}

// Baggage return the baggage from kctx.
func Baggage(c context.Context) baggage.Baggage {
	return baggage.FromContext(c)
}

// BaggageSetMembers sets baggage members for a span. If you already have a reference to the baggage, prefer
// using [go.opentelemetry.io/otel/baggage].Baggage.SetMember()
// Please not that returned context needs to be propagated in order for the baggage to be propagated, too.
func BaggageSet(c context.Context, members ...baggage.Member) (context.Context, error) {
	bag := baggage.FromContext(c)

	var err error

	for _, member := range members {
		bag, err = bag.SetMember(member)
		if err != nil {
			return nil, fmt.Errorf("error setting baggage member: %w", err)
		}
	}

	return baggage.ContextWithBaggage(c, bag), nil
}
