package kctx

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/trace"
)

// Span return the span from kctx
func (ctx *Kctx) Span(r *http.Request) trace.Span {
	return trace.SpanFromContext(r.Context())
}

// SpanCtx return the span context from kctx. If you already have a reference to the span, prefer
// using [go.opentelemetry.io/otel/trace].Span.SpanContext()
func (ctx *Kctx) SpanCtx(r *http.Request) trace.SpanContext {
	return ctx.Span(r).SpanContext()
}

// SpanSetAttrs sets attributes for a span. If you already have a reference to the span, prefer
// using [go.opentelemetry.io/otel/trace].Span.SetAttributes()
func (ctx *Kctx) SpanSetAttrs(r *http.Request, kv ...attribute.KeyValue) {
	span := ctx.Span(r)
	span.SetAttributes(kv...)
}

// Baggage return the baggage from kctx
func (ctx *Kctx) Baggage(r *http.Request) baggage.Baggage {
	return baggage.FromContext(r.Context())
}

// BaggageSetMembers sets baggage members for a span. If you already have a reference to the baggage, prefer
// using [go.opentelemetry.io/otel/baggage].Baggage.SetMember()
func (ctx *Kctx) BaggageSet(r *http.Request, members ...baggage.Member) error {
	bag := baggage.FromContext(r.Context())
	for _, member := range members {
		var err error
		bag, err = bag.SetMember(member)
		if err != nil {
			return fmt.Errorf("error setting baggage member: %w", err)
		}
	}

	return nil
}
