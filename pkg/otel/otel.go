// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package otel

import (
	"context"
	"encoding/json"
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

// SetupOTelSDK sets up the OpenTelemetry SDK with the provided configuration.
// It returns a function that can be called to shut down the OpenTelemetry SDK, and an error if any occurred during the setup.
// The function returned by SetupOTelSDK should be called to shut down the OpenTelemetry SDK.
// If it does not return an error, a propoer call to shutdown is needed to
// clean up the OpenTelemetry SDK.
func SetupOTelSDK(
	ctx context.Context,
	conf config.Global,
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
	otel.SetTextMapPropagator(newPropagator())

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
			semconv.ServiceName(conf.Runtime.AppName),
			semconv.ServiceNamespace(conf.Runtime.AppNamespace),
			semconv.ServiceVersion(conf.Runtime.AppVersion),
			semconv.DeploymentEnvironmentName(conf.Runtime.Environment),
			attribute.KeyValue{
				Key: "process.config",
				Value: attribute.StringValue(func() string {
					d, err := json.Marshal(conf)
					if err != nil {
						return fmt.Sprintf("%+v", conf)
					}
					return string(d)
				}()),
			},
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
	conf config.Global,
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
		otlploggrpc.WithCompressor(conf.Observability.ExporterCompression),
		otlploggrpc.WithEndpointURL(conf.Observability.EndpointURL),
	)
	if err != nil {
		return nil, fmt.Errorf("otel logger init: %w", err)
	}

	grpcBatchProcessor := log.NewBatchProcessor(
		grpcExporter,
	)

	// Log Info by default, Debug for dev
	logLevel := minsev.SeverityInfo
	if conf.Runtime.IsLocalEnvironment() {
		logLevel = minsev.SeverityDebug
	}

	// Wrap the processor so that it filters by severity level
	stdoutProcessor := minsev.NewLogProcessor(stdoutSimpleProcessor, logLevel)

	// Only output to stdout during local development
	if conf.Runtime.IsLocalEnvironment() {
		provider := log.NewLoggerProvider(
			log.WithResource(res),
			log.WithProcessor(stdoutProcessor),
		)
		return provider, nil
	}

	grpcProcessor := minsev.NewLogProcessor(grpcBatchProcessor, logLevel)

	// Outut as OLTP as well as stdout
	provider := log.NewLoggerProvider(
		log.WithResource(res),
		log.WithProcessor(stdoutProcessor),
		log.WithProcessor(grpcProcessor),
	)

	return provider, nil
}

// newMeterProvider returns a new OpenTelemetry meter provider, and an error if any occurred during the setup.
// The meter provider is configured to batch export metrics to the OpenTelemetry collector, or synchrounously to stdout
// if conf.Runtime.Environment is set to [github.com/kemadev/go-framework/pkg/config.EnvDev].
func newMeterProvider(
	ctx context.Context,
	res *resource.Resource,
	conf config.Global,
) (*metric.MeterProvider, error) {
	var exporter metric.Exporter

	if conf.Runtime.IsLocalEnvironment() {
		// Do not export metrics when export interval is 0 or below
		if conf.Observability.MetricsExportIntervalSeconds <= 0 {
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
			otlpmetricgrpc.WithCompressor(conf.Observability.ExporterCompression),
			otlpmetricgrpc.WithEndpointURL(conf.Observability.EndpointURL),
		)
		if err != nil {
			return nil, fmt.Errorf("otel metric init: %w", err)
		}
		exporter = exp
	}

	proc := metric.NewPeriodicReader(
		exporter,
		metric.WithInterval(
			time.Duration(conf.Observability.MetricsExportIntervalSeconds)*time.Second,
		),
	)

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(proc),
		metric.WithResource(res),
		metric.WithExemplarFilter(exemplar.AlwaysOnFilter),
	)

	return meterProvider, nil
}

// newTracerProvider returns a new OpenTelemetry tracer provider, and an error if any occurred during the setup.
// The tracer provider is configured to batch export traces to the OpenTelemetry collector
func newTracerProvider(
	ctx context.Context,
	res *resource.Resource,
	conf config.Global,
) (*trace.TracerProvider, error) {
	exp, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithCompressor(conf.Observability.ExporterCompression),
		otlptracegrpc.WithEndpointURL(conf.Observability.EndpointURL),
	)
	if err != nil {
		return nil, fmt.Errorf("otel tracer init: %w", err)
	}

	batcher := trace.WithBatcher(exp)

	if conf.Runtime.IsLocalEnvironment() {
		batcher = trace.WithSyncer(exp)
	}

	tracerProvider := trace.NewTracerProvider(
		batcher,
		trace.WithSampler(
			trace.ParentBased(
				trace.TraceIDRatioBased(float64(conf.Observability.TracingSamplePercent/100)),
			),
		),
		trace.WithResource(res),
	)

	return tracerProvider, nil
}
