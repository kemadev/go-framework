// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

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
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/log/global"
	nometric "go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/exemplar"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
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
			attribute.String(string(semconv.DeploymentEnvironmentNameKey), conf.RuntimeEnv),
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

	// Set up meter provider.
	meterProvider, err := newMeterProvider(ctx, res, conf)
	if err != nil {
		handleErr(err)

		return nil, fmt.Errorf("otel meter provider: %w", err)
	}

	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	// Set up trace provider.
	tracerProvider, err := newTracerProvider(ctx, res, conf)
	if err != nil {
		handleErr(err)

		return nil, fmt.Errorf("otel tracer provider: %w", err)
	}

	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

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
func newLoggerProvider(
	ctx context.Context,
	res *resource.Resource,
	conf config.Config,
) (*log.LoggerProvider, error) {
	stdoutExporter, err := klog.NewExporter()
	if err != nil {
		return nil, fmt.Errorf("otel logger init: %w", err)
	}

	stdoutSimpleProcessor := log.NewSimpleProcessor(
		stdoutExporter,
	)

	grpcExporter, err := otlploggrpc.New(
		ctx,
		otlploggrpc.WithCompressor(conf.OtelExporterCompression),
		otlploggrpc.WithEndpointURL(conf.OtelEndpointURL.String()),
	)
	if err != nil {
		return nil, fmt.Errorf("otel logger init: %w", err)
	}

	grpcBatchProcessor := log.NewBatchProcessor(
		grpcExporter,
	)

	// Log Info by default, Debug for dev
	logLevel := minsev.SeverityInfo
	if conf.RuntimeEnv == config.EnvDev {
		logLevel = minsev.SeverityDebug
	}

	// Wrap the processor so that it filters by severity level
	stdoutProcessor := minsev.NewLogProcessor(stdoutSimpleProcessor, logLevel)

	if conf.RuntimeEnv == config.EnvDev {
		provider := log.NewLoggerProvider(
			log.WithResource(res),
			log.WithProcessor(stdoutProcessor),
		)
		return provider, nil
	}

	grpcProcessor := minsev.NewLogProcessor(grpcBatchProcessor, logLevel)

	provider := log.NewLoggerProvider(
		log.WithResource(res),
		log.WithProcessor(stdoutProcessor),
		log.WithProcessor(grpcProcessor),
	)

	return provider, nil
}

// newMeterProvider returns a new OpenTelemetry meter provider, and an error if any occurred during the setup.
// The meter provider is configured to batch export metrics to the OpenTelemetry collector, or synchrounously to stdout
// if conf.RuntimeEnv is set to [github.com/kemadev/go-framework/pkg/config.EnvDev].
func newMeterProvider(
	ctx context.Context,
	res *resource.Resource,
	conf config.Config,
) (*metric.MeterProvider, error) {
	var exporter metric.Exporter

	if conf.RuntimeEnv == config.EnvDev {
		if conf.MetricsExportInterval <= 0 {
			prov := nometric.NewMeterProvider()
			return &metric.MeterProvider{
				MeterProvider: prov,
			}, nil
		} else {
			exp, err := stdoutmetric.New(
				stdoutmetric.WithPrettyPrint(),
			)
			if err != nil {
				return nil, fmt.Errorf("otel metric init: %w", err)
			}
			exporter = exp
		}
	} else {
		exp, err := otlpmetricgrpc.New(
			ctx,
			otlpmetricgrpc.WithCompressor(conf.OtelExporterCompression),
			otlpmetricgrpc.WithEndpointURL(conf.OtelEndpointURL.String()),
		)
		if err != nil {
			return nil, fmt.Errorf("otel metric init: %w", err)
		}
		exporter = exp
	}

	proc := metric.NewPeriodicReader(
		exporter,
		metric.WithInterval(time.Duration(conf.MetricsExportInterval)*time.Second),
	)

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(proc),
		metric.WithResource(res),
		metric.WithExemplarFilter(exemplar.TraceBasedFilter),
	)

	return meterProvider, nil
}

// newTracerProvider returns a new OpenTelemetry tracer provider, and an error if any occurred during the setup.
// The tracer provider is configured to batch export traces to the OpenTelemetry collector
func newTracerProvider(
	ctx context.Context,
	res *resource.Resource,
	conf config.Config,
) (*trace.TracerProvider, error) {
	exp, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithCompressor(conf.OtelExporterCompression),
		otlptracegrpc.WithEndpointURL(conf.OtelEndpointURL.String()),
	)
	if err != nil {
		return nil, fmt.Errorf("otel tracer init: %w", err)
	}

	batcher := trace.WithBatcher(exp)

	if conf.RuntimeEnv == config.EnvDev {
		batcher = trace.WithSyncer(exp)
	}

	tracerProvider := trace.NewTracerProvider(
		batcher,
		trace.WithSampler(
			trace.ParentBased(
				trace.TraceIDRatioBased(conf.TracesSampleRatio),
			),
		),
		trace.WithResource(res),
	)

	return tracerProvider, nil
}
