package otelfailsafe

import (
	"context"
	"fmt"
	"time"

	"github.com/failsafe-go/failsafe-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type FailsafeMetrics struct {
	attemptCounter    metric.Int64Counter
	executionCounter  metric.Int64Counter
	executionDuration metric.Float64Histogram
	retryCounter      metric.Int64Counter
	hedgeCounter      metric.Int64Counter
	successCounter    metric.Int64Counter
	failureCounter    metric.Int64Counter
}

const packageName = "github.com/kemadev/go-framework/pkg/otelfailsafe"

func newFailsafeMetrics(name string) (*FailsafeMetrics, error) {
	meter := otel.GetMeterProvider().Meter(packageName + "-" + name)

	attemptCounter, err := meter.Int64Counter(
		"failsafe.attempts.total",
		metric.WithDescription("Total number of failsafe attempts"),
		metric.WithUnit("{attempt}"),
	)
	if err != nil {
		return nil, err
	}

	executionCounter, err := meter.Int64Counter(
		"failsafe.executions.total",
		metric.WithDescription("Total number of failsafe executions"),
		metric.WithUnit("{execution}"),
	)
	if err != nil {
		return nil, err
	}

	retryCounter, err := meter.Int64Counter(
		"failsafe.retries.total",
		metric.WithDescription("Total number of retry attempts"),
		metric.WithUnit("{retry}"),
	)
	if err != nil {
		return nil, err
	}

	hedgeCounter, err := meter.Int64Counter(
		"failsafe.hedges.total",
		metric.WithDescription("Total number of hedge attempts"),
		metric.WithUnit("{hedge}"),
	)
	if err != nil {
		return nil, err
	}

	executionDuration, err := meter.Float64Histogram(
		"failsafe.execution.duration",
		metric.WithDescription("Duration of failsafe executions"),
		metric.WithUnit("ms"),
	)
	if err != nil {
		return nil, err
	}

	successCounter, err := meter.Int64Counter(
		"failsafe.successs.total",
		metric.WithDescription("Total number of success attempts"),
		metric.WithUnit("{success}"),
	)
	if err != nil {
		return nil, err
	}

	failureCounter, err := meter.Int64Counter(
		"failsafe.failures.total",
		metric.WithDescription("Total number of failure attempts"),
		metric.WithUnit("{failure}"),
	)
	if err != nil {
		return nil, err
	}

	return &FailsafeMetrics{
		executionCounter:  executionCounter,
		attemptCounter:    attemptCounter,
		executionDuration: executionDuration,
		retryCounter:      retryCounter,
		hedgeCounter:      hedgeCounter,
		successCounter:    successCounter,
		failureCounter:    failureCounter,
	}, nil
}

func NewExecutor[R any](name string, policies ...failsafe.Policy[R]) (failsafe.Executor[R], error) {
	meter, err := newFailsafeMetrics(name)
	if err != nil {
		return nil, fmt.Errorf("error creating failsafe metrics: %w", err)
	}

	executorNameAttribute := attribute.String("failsafe.executor.name", name)

	return failsafe.With(policies...).
		OnDone(func(ede failsafe.ExecutionDoneEvent[R]) {
			ctx := context.Background()
			meter.attemptCounter.Add(ctx, int64(ede.Attempts()), metric.WithAttributes(executorNameAttribute))
			meter.executionCounter.Add(ctx, int64(ede.Executions()), metric.WithAttributes(executorNameAttribute))
			meter.executionDuration.Record(ctx, float64(ede.ElapsedTime()*time.Millisecond), metric.WithAttributes(executorNameAttribute))
			meter.retryCounter.Add(ctx, int64(ede.Retries()), metric.WithAttributes(executorNameAttribute))
			meter.hedgeCounter.Add(ctx, int64(ede.Hedges()), metric.WithAttributes(executorNameAttribute))
		}).
		OnSuccess(func(ede failsafe.ExecutionDoneEvent[R]) {
			ctx := context.Background()
			meter.successCounter.Add(ctx, 1, metric.WithAttributes(executorNameAttribute))
		}).
		OnFailure(func(ede failsafe.ExecutionDoneEvent[R]) {
			ctx := context.Background()
			meter.failureCounter.Add(ctx, 1, metric.WithAttributes(executorNameAttribute))
		}), nil
}
