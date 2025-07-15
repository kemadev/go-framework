package log

import (
	"log/slog"
	"os"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// CreateFallbackLogger returns a fallback logger that is used when the OpenTelemetry logger is not available.
// It uses the slog package to create a JSON logger that writes to stdout.
func CreateFallbackLogger() *slog.Logger {
	return slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		),
	)
}

// AddAttributesToInstruments adds attributes to the given meter, tracer, and logger, and returns the updated instances.
func AddAttributesToInstruments(
	name string,
	meter metric.Meter,
	tracer trace.Span,
	logger *slog.Logger,
	attributes []attribute.KeyValue,
) (metric.Meter, trace.Span, *slog.Logger) {
	for _, attr := range attributes {
		var slogAttr interface{}
		switch attr.Value.Type() {
		case attribute.STRING:
			slogAttr = slog.String(string(attr.Key), attr.Value.AsString())
		case attribute.INT64:
			slogAttr = slog.Int64(string(attr.Key), attr.Value.AsInt64())
		case attribute.FLOAT64:
			slogAttr = slog.Float64(string(attr.Key), attr.Value.AsFloat64())
		case attribute.BOOL:
			slogAttr = slog.Bool(string(attr.Key), attr.Value.AsBool())
		case attribute.STRINGSLICE:
			slogAttr = slog.String(string(attr.Key), strings.Join(attr.Value.AsStringSlice(), ""))
		case attribute.INT64SLICE:
			slogAttr = slog.String(string(attr.Key), strings.Join(attr.Value.AsStringSlice(), ""))
		case attribute.FLOAT64SLICE:
			slogAttr = slog.String(string(attr.Key), strings.Join(attr.Value.AsStringSlice(), ""))
		case attribute.BOOLSLICE:
			slogAttr = slog.String(string(attr.Key), strings.Join(attr.Value.AsStringSlice(), ""))
		default:
			slogAttr = attr.Value.AsInterface()
		}
		meter = otel.Meter(name, metric.WithInstrumentationAttributes(attr))
		tracer.SetAttributes(attr)
		logger = logger.With(slogAttr)
	}
	return meter, tracer, logger
}
