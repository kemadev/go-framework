// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"log/slog"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"

	khttp "github.com/kemadev/go-framework/pkg/http"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

// NOTE This file should be placed in a dedicated package as per Go standards, it is
// placed in main package for ease of demo purpose only

// Name for telemetry, should reflect go import structure.
const name = "github.com/kemadev/go-framework/cmd/main"

// Telemetry components should be package-scoped.
var (
	// Tracing instrument.
	tracer = otel.Tracer(name)
	// Metring instrument.
	meter = otel.Meter(name)
	// Logging instrument.
	logger = otelslog.NewLogger(
		"default",
		otelslog.WithSource(true),
	)
	// Count of rolls.
	rollCnt metric.Int64Counter
)

// init initializes the telemetry components.
func init() {
	var err error
	// Other types are available, for recommendations see https://opentelemetry.io/docs/specs/otel/metrics/api/
	rollCnt, err = meter.Int64Counter("dice.rolls",
		metric.WithDescription("The number of rolls by roll value"),
		metric.WithUnit("{roll}"))
	if err != nil {
		panic(err)
	}
}

// rolldice is the handler for the /rolldice/{player} endpoint.
func rolldice(w http.ResponseWriter, r *http.Request) {
	// Init trace span
	ctx, span := tracer.Start(
		// Get propagated context from parent, or create new one if there is none
		r.Context(),
		// Use a meaningful name, see https://opentelemetry.io/docs/specs/otel/trace/api/#span
		"roll",
		// Set appropriate kind, see https://opentelemetry.io/docs/specs/otel/trace/api/#spankind
		trace.WithSpanKind(trace.SpanKindServer),
	)
	defer span.End()

	kclient := khttp.HTTPClientInfo{
		Ctx:    ctx,
		Writer: w,
		Logger: logger,
		Span:   span,
	}

	// Please note that you can also link spans using `span.AddLink()`, see https://opentelemetry.io/docs/specs/otel/trace/api/#link
	// Using following links whenever possible is encouraged
	// Users - https://opentelemetry.io/docs/specs/semconv/attributes-registry/user/
	// Web sessions - https://opentelemetry.io/docs/specs/semconv/attributes-registry/session/
	// Feature flags - https://opentelemetry.io/docs/specs/semconv/attributes-registry/feature-flag/

	// Run business logic
	// nosemgrep: go.lang.security.audit.crypto.math_random.math-random-used // We are ok using pseudo-random number generator here, this is a demo
	roll := 1 + rand.Intn(6)

	var msg, player string

	player = r.PathValue("player")
	if player != "" {
		msg = player + " is rolling the dice"
	} else {
		msg = "Anonymous player is rolling the dice"
	}

	// You can add span events, see https://opentelemetry.io/docs/concepts/signals/traces/#span-events
	span.AddEvent("roll computed")

	// Use logger as normal `*slog.Logger`, still benefit from instrumentation
	// Attributes namming should respect semantic conventions as much as possible, see https://opentelemetry.io/docs/specs/semconv/
	// e.g. enduser.id comes from https://opentelemetry.io/docs/specs/semconv/attributes-registry/enduser/
	// Same goes for metrics and traces, use semantic conventions whenever you can!
	// Please note that some attributes are already packaged by the exporters thanks to their resource (see
	// https://opentelemetry.io/docs/languages/go/resources/) being automatically instrumented during OpenTelemetry SDK initialization.
	// Some more instrumentation is handled by using instrumented libraries like HTTP server, slog logger, ...
	// That is, only include attributes your programs *handles*, not attributes defining what your program *is*, e.g. do not include
	// HTTP status code which is already provided by HTTP server instrumentation, but include enduser ID which you got from business logic
	logger.InfoContext(
		ctx,
		msg,
		slog.Int("result", roll),
		slog.String(string(semconv.EnduserIDKey), player),
	)

	// Enrich trace, once again using a meaningful attribute name
	rollValueAttr := attribute.Int("roll.value", roll)
	span.SetAttributes(rollValueAttr)

	// You can call function, just remember to pass context!
	subTask(ctx, roll)

	// Add metric
	rollCnt.Add(ctx, 1, metric.WithAttributes(rollValueAttr))

	// Use helper function to send JSON responses, automatically instrumented. You
	// can use `w` directly, but requires more boilerplate. `khttp.SendErrorResponse` can
	// also be used to return errors to clients
	khttp.SendJSONResponse(
		kclient,
		http.StatusOK,
		map[string]int{
			"roll": roll,
		},
	)

	// Remember to set status
	span.SetStatus(codes.Ok, "roll successful")
}

// subTask is a helper function to show how to call another service.
func subTask(ctx context.Context, roll int) (*http.Response, error) {
	// You can create nested spans, see https://opentelemetry.io/docs/languages/go/instrumentation/#create-nested-spans
	// Nested spans are useful to track different parts of the handling, such as external calls like in this example
	ctx, childSpan := tracer.Start(ctx, "subtask")
	defer childSpan.End()

	// Calling an instrumented receiver using an instrumented library, e.g. configured `otelhttp` client calling another configured
	// `otelhttp` server, context will propagate as needed, linking spans automatically. When calling non-instrumented clients,
	// you need to pass context manually, client needs to retrieve it and use it, e.g. by using `trace.ContextWithRemoteSpanContext()`
	// (using the retreived context) as context when calling `tracer.Start()`
	return otelhttp.Post(
		ctx,
		(&url.URL{
			Scheme: "http",
			Host:   "localhost:8080",
			Path:   "/" + strconv.Itoa(roll),
		}).String(),
		"application/octet-stream",
		http.NoBody,
	)
}
