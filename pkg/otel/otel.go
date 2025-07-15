package otel

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kemadev/go-framework/pkg/config"
	klog "github.com/kemadev/go-framework/pkg/log"
	"go.opentelemetry.io/contrib/processors/minsev"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

var (
	// ErrOtelExporterCompressionInvalid is a sentinel error that indicates that the OpenTelemetry exporter compression is invalid.
	ErrOtelCompressionInvalid = fmt.Errorf("otel compression is invalid")
	// ErrBuildInfoFail is a sentinel error that indicates that the build info extraction failed.
	ErrBuildInfoFail = fmt.Errorf("build info extraction failed")
)

// SetupOTelSDK returns a function that can be called to shut down the OpenTelemetry SDK, and an error if any occurred during the setup.
// The function returned by SetupOTelSDK should be called to shut down the OpenTelemetry SDK.
// It sets up the OpenTelemetry SDK with the provided configuration.
// If it does not return an error, a propoer call to shutdown is needed to
// clean up the OpenTelemetry SDK.
func SetupOTelSDK(
	ctx context.Context,
	conf config.Config,
) (func(context.Context) error, error) {
	var err error

	var shutdownFuncs []func(context.Context) error

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown := func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}

		shutdownFuncs = nil

		return err
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// Set up propagator.
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	// Set up resource in order to enrich telemetry data.
	res, err := resource.New(
		context.Background(),
		resource.WithSchemaURL(semconv.SchemaURL),
		resource.WithFromEnv(),
		resource.WithHost(),
		resource.WithOS(),
		resource.WithProcess(),
		resource.WithContainer(),
		resource.WithTelemetrySDK(),
		resource.WithAttributes(
			attribute.String(string(semconv.ServiceNamespaceKey), conf.AppNamespace),
			attribute.String(string(semconv.ServiceNameKey), conf.AppName),
			attribute.String(string(semconv.ServiceVersionKey), conf.AppVersion),
			attribute.String(string(semconv.DeploymentEnvironmentKey), conf.RuntimeEnv),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("otel resource creation: %w", err)
	}

	// Set up logger provider.
	loggerProvider, err := newLoggerProvider(ctx, res, conf)
	if err != nil {
		handleErr(err)

		return nil, fmt.Errorf("otel logger provider: %w", err)
	}

	shutdownFuncs = append(shutdownFuncs, loggerProvider.Shutdown)
	global.SetLoggerProvider(loggerProvider)

	// Set up trace provider.
	tracerProvider, err := newTracerProvider(ctx, res, conf)
	if err != nil {
		handleErr(err)

		return nil, fmt.Errorf("otel tracer provider: %w", err)
	}

	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	// Set up meter provider.
	meterProvider, err := newMeterProvider(ctx, res, conf)
	if err != nil {
		handleErr(err)

		return nil, fmt.Errorf("otel meter provider: %w", err)
	}

	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	return shutdown, nil
}

// newPropagator returns a new OpenTelemetry propagator.
func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

// newLoggerProvider returns a new OpenTelemetry logger provider, and an error if any occurred during the setup.
// The logger provider is configured to batch export logs to the OpenTelemetry collector, or synchrounously to stdout using
// [github.com/kemadev/go-framework/pkg/log] if conf.RuntimeEnv is set to [github.com/kemadev/go-framework/pkg/config.Env_dev].
func newLoggerProvider(
	ctx context.Context,
	res *resource.Resource,
	conf config.Config,
) (*log.LoggerProvider, error) {
	var processor log.Processor

	// Export to stdout w/o batching in dev, batch export to collector via OLTP otherwise
	if conf.RuntimeEnv == config.Env_dev {
		exp, err := klog.NewExporter()
		if err != nil {
			return nil, fmt.Errorf("otel logger init: %w", err)
		}

		p := log.NewSimpleProcessor(
			exp,
		)
		processor = p
	} else {
		exp, err := otlploggrpc.New(
			ctx,
			otlploggrpc.WithCompressor(conf.OtelExporterCompression),
			otlploggrpc.WithEndpointURL(conf.OtelEndpointUrl.String()),
		)
		if err != nil {
			return nil, fmt.Errorf("otel logger init: %w", err)
		}

		p := log.NewBatchProcessor(
			exp,
		)
		processor = p
	}

	// Log Info by default, Debug for dev
	logLevel := minsev.SeverityInfo
	if conf.RuntimeEnv == config.Env_dev {
		logLevel = minsev.SeverityDebug
	}

	// Wrap the processor so that it filters by severity level
	processorChain := minsev.NewLogProcessor(processor, logLevel)

	provider := log.NewLoggerProvider(
		log.WithProcessor(processorChain),
		log.WithResource(res),
	)

	return provider, nil
}

// newTracerProvider returns a new OpenTelemetry tracer provider, and an error if any occurred during the setup.
// The tracer provider is configured to batch export traces to the OpenTelemetry collector.
func newTracerProvider(
	ctx context.Context,
	res *resource.Resource,
	conf config.Config,
) (*trace.TracerProvider, error) {
	exp, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithCompressor(conf.OtelExporterCompression),
		otlptracegrpc.WithEndpointURL(conf.OtelEndpointUrl.String()),
	)
	if err != nil {
		return nil, fmt.Errorf("otel tracer init: %w", err)
	}

	tracerProvider := trace.NewTracerProvider(
		// Always sample is required for tail sampling
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(exp),
		trace.WithResource(res),
	)

	return tracerProvider, nil
}

// newMeterProvider returns a new OpenTelemetry meter provider, and an error if any occurred during the setup.
// The meter provider is configured to batch export metrics to the OpenTelemetry collector.
func newMeterProvider(
	ctx context.Context,
	res *resource.Resource,
	conf config.Config,
) (*metric.MeterProvider, error) {
	exp, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithCompressor(conf.OtelExporterCompression),
		otlpmetricgrpc.WithEndpointURL(conf.OtelEndpointUrl.String()),
	)
	if err != nil {
		return nil, fmt.Errorf("otel metric init: %w", err)
	}

	// Shorter interval in dev
	var batchInterval time.Duration

	switch conf.RuntimeEnv {
	case config.Env_dev:
		batchInterval = 5 * time.Second
	default:
		batchInterval = 30 * time.Second
	}

	processor := metric.NewPeriodicReader(
		exp,
		metric.WithInterval(batchInterval),
	)

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(processor),
		metric.WithResource(res),
	)

	return meterProvider, nil
}
