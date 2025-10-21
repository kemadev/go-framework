package otelfailsafe

import (
	"context"
	"fmt"
	"time"

	"github.com/failsafe-go/failsafe-go"
	"github.com/failsafe-go/failsafe-go/adaptivelimiter"
	"github.com/failsafe-go/failsafe-go/adaptivethrottler"
	"github.com/failsafe-go/failsafe-go/bulkhead"
	"github.com/failsafe-go/failsafe-go/cachepolicy"
	"github.com/failsafe-go/failsafe-go/circuitbreaker"
	"github.com/failsafe-go/failsafe-go/fallback"
	"github.com/failsafe-go/failsafe-go/hedgepolicy"
	"github.com/failsafe-go/failsafe-go/ratelimiter"
	"github.com/failsafe-go/failsafe-go/retrypolicy"
	"github.com/failsafe-go/failsafe-go/timeout"
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
	fallbackCounter   metric.Int64Counter
	timeoutCounter    metric.Int64Counter
	cacheCounter      metric.Int64Counter
	successCounter    metric.Int64Counter
	failureCounter    metric.Int64Counter

	adaptiveLimit       metric.Int64Gauge
	circuitbreakerState metric.Int64Gauge
}

const packageName = "github.com/kemadev/go-framework/pkg/otelfailsafe"

type PolicyEngine[R any] struct {
	nameAttributeMeasurmentOption metric.MeasurementOption
	Metrics                       *FailsafeMetrics
}

func NewPolicyEngine[R any](name string) (PolicyEngine[R], error) {
	metrics, err := newFailsafeMetrics(name)
	if err != nil {
		return PolicyEngine[R]{}, fmt.Errorf("error creating failsafe metrics: %w", err)
	}

	return PolicyEngine[R]{
		nameAttributeMeasurmentOption: metric.WithAttributes(attribute.String("failsafe.executor.name", name)),
		Metrics:                       metrics,
	}, nil
}

func newFailsafeMetrics(name string) (*FailsafeMetrics, error) {
	meter := otel.GetMeterProvider().Meter(packageName)

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

	executionDuration, err := meter.Float64Histogram(
		"failsafe.execution.duration",
		metric.WithDescription("Duration of failsafe executions"),
		metric.WithUnit("ms"),
	)
	if err != nil {
		return nil, err
	}

	retryCounter, err := meter.Int64Counter(
		"failsafe.retries.total",
		metric.WithDescription("Total number of failsafe retry attempts"),
		metric.WithUnit("{retry}"),
	)
	if err != nil {
		return nil, err
	}

	hedgeCounter, err := meter.Int64Counter(
		"failsafe.hedges.total",
		metric.WithDescription("Total number of failsafe hedge attempts"),
		metric.WithUnit("{hedge}"),
	)
	if err != nil {
		return nil, err
	}

	fallbackCounter, err := meter.Int64Counter(
		"failsafe.fallbacks.total",
		metric.WithDescription("Total number of fallbacks"),
		metric.WithUnit("{hedge}"),
	)
	if err != nil {
		return nil, err
	}

	timeoutCounter, err := meter.Int64Counter(
		"failsafe.timeouts.total",
		metric.WithDescription("Total number of timeouts"),
		metric.WithUnit("{hedge}"),
	)
	if err != nil {
		return nil, err
	}

	cacheCounter, err := meter.Int64Counter(
		"failsafe.caches.total",
		metric.WithDescription("Total number of caches"),
		metric.WithUnit("{hedge}"),
	)
	if err != nil {
		return nil, err
	}

	failureCounter, err := meter.Int64Counter(
		"failsafe.failures.total",
		metric.WithDescription("Total number of failsafe failures"),
		metric.WithUnit("{failure}"),
	)
	if err != nil {
		return nil, err
	}

	successCounter, err := meter.Int64Counter(
		"failsafe.successs.total",
		metric.WithDescription("Total number of failsafe successes"),
		metric.WithUnit("{success}"),
	)
	if err != nil {
		return nil, err
	}

	adaptiveLimit, err := meter.Int64Gauge(
		"failsafe.adaptivelimit.limit",
		metric.WithDescription("Adaptive limit limit"),
		metric.WithUnit("{limit}"),
	)
	if err != nil {
		return nil, err
	}

	circuitbreakerState, err := meter.Int64Gauge(
		"failsafe.circuitbreaker.state",
		metric.WithDescription(
			fmt.Sprintf(
				"Circuit breaker state (%d=%s %d=%s %d=%s)",
				circuitbreaker.ClosedState,
				circuitbreaker.ClosedState.String(),
				circuitbreaker.OpenState,
				circuitbreaker.OpenState.String(),
				circuitbreaker.HalfOpenState,
				circuitbreaker.HalfOpenState.String(),
			),
		),
		metric.WithUnit("{state}"),
	)
	if err != nil {
		return nil, err
	}

	return &FailsafeMetrics{
		executionCounter:    executionCounter,
		attemptCounter:      attemptCounter,
		executionDuration:   executionDuration,
		retryCounter:        retryCounter,
		hedgeCounter:        hedgeCounter,
		fallbackCounter:     fallbackCounter,
		timeoutCounter:      timeoutCounter,
		cacheCounter:        cacheCounter,
		successCounter:      successCounter,
		failureCounter:      failureCounter,
		adaptiveLimit:       adaptiveLimit,
		circuitbreakerState: circuitbreakerState,
	}, nil
}

func (p PolicyEngine[R]) NewCircuitBreakerBuilder() circuitbreaker.Builder[R] {
	return circuitbreaker.NewBuilder[R]().
		OnStateChanged(func(sce circuitbreaker.StateChangedEvent) {
			ctx := context.Background()
			p.Metrics.adaptiveLimit.Record(ctx, int64(sce.NewState), p.nameAttributeMeasurmentOption)
		})
}

func (p PolicyEngine[R]) NewAdaptiveLimiterBuilder() adaptivelimiter.Builder[R] {
	return adaptivelimiter.NewBuilder[R]().
		OnLimitChanged(func(event adaptivelimiter.LimitChangedEvent) {
			ctx := context.Background()
			p.Metrics.adaptiveLimit.Record(ctx, int64(event.NewLimit), p.nameAttributeMeasurmentOption)
		})
}

func (p PolicyEngine[R]) NewAdaptiveThrottlerBuilder() adaptivethrottler.Builder[R] {
	return adaptivethrottler.NewBuilder[R]()
}

func (p PolicyEngine[R]) NewTimeoutBuilder(timeLimit time.Duration) timeout.Builder[R] {
	return timeout.NewBuilder[R](timeLimit).
		OnTimeoutExceeded(func(event failsafe.ExecutionDoneEvent[R]) {
			ctx := context.Background()
			p.Metrics.timeoutCounter.Add(ctx, 1, p.nameAttributeMeasurmentOption)
		})
}

func (p PolicyEngine[R]) NewFallbackWithErrorBuilder(err error) fallback.Builder[R] {
	return fallback.NewBuilderWithError[R](err).
		OnFallbackExecuted(func(event failsafe.ExecutionDoneEvent[R]) {
			ctx := context.Background()
			p.Metrics.failureCounter.Add(ctx, 1, p.nameAttributeMeasurmentOption)
		})
}

func (p PolicyEngine[R]) NewFallbackWithFuncBuilder(f func(exec failsafe.Execution[R]) (R, error)) fallback.Builder[R] {
	return fallback.NewBuilderWithFunc(f).
		OnFallbackExecuted(func(event failsafe.ExecutionDoneEvent[R]) {
			ctx := context.Background()
			p.Metrics.failureCounter.Add(ctx, 1, p.nameAttributeMeasurmentOption)
		})
}

func (p PolicyEngine[R]) NewFallbackWithResultBuilder(result R) fallback.Builder[R] {
	return fallback.NewBuilderWithResult(result).
		OnFallbackExecuted(func(event failsafe.ExecutionDoneEvent[R]) {
			ctx := context.Background()
			p.Metrics.failureCounter.Add(ctx, 1, p.nameAttributeMeasurmentOption)
		})
}

func (p PolicyEngine[R]) NewRetryBuilder() retrypolicy.Builder[R] {
	return retrypolicy.NewBuilder[R]()
}

func (p PolicyEngine[R]) NewRateLimiterBurstyBuilder(maxExecutions uint, period time.Duration) ratelimiter.Builder[R] {
	return ratelimiter.NewBurstyBuilder[R](maxExecutions, period)
}

func (p PolicyEngine[R]) NewRateLimiterSmoothBuilder(maxExecutions uint, period time.Duration) ratelimiter.Builder[R] {
	return ratelimiter.NewSmoothBuilder[R](maxExecutions, period)
}

func (p PolicyEngine[R]) NewRateLimiterSmoothWithMaxRateBuilder(maxRate time.Duration) ratelimiter.Builder[R] {
	return ratelimiter.NewSmoothBuilderWithMaxRate[R](maxRate)
}

func (p PolicyEngine[R]) NewHedgeWithDelayBuilder(delay time.Duration) hedgepolicy.Builder[R] {
	return hedgepolicy.NewBuilderWithDelay[R](delay)
}

func (p PolicyEngine[R]) NewHedgeWithDelayFuncBuilder(delayFunc failsafe.DelayFunc[R]) hedgepolicy.Builder[R] {
	return hedgepolicy.NewBuilderWithDelayFunc(delayFunc)
}

func (p PolicyEngine[R]) NewBulkheadBuilder(maxConcurrency uint) bulkhead.Builder[R] {
	return bulkhead.NewBuilder[R](maxConcurrency)
}

func (p PolicyEngine[R]) NewCacheBuilder(cache cachepolicy.Cache[R]) cachepolicy.Builder[R] {
	return cachepolicy.NewBuilder(cache).
		OnCacheMiss(func(event failsafe.ExecutionEvent[R]) {
			ctx := context.Background()
			p.Metrics.cacheCounter.Add(ctx, 1, p.nameAttributeMeasurmentOption)
		}).
		OnCacheHit(func(event failsafe.ExecutionDoneEvent[R]) {
			ctx := context.Background()
			p.Metrics.cacheCounter.Add(ctx, 1, p.nameAttributeMeasurmentOption)
		})
}

func (p PolicyEngine[R]) NewExecutor(policies ...failsafe.Policy[R]) failsafe.Executor[R] {
	metricRegister := func(success bool) func(ede failsafe.ExecutionDoneEvent[R]) {
		return func(ede failsafe.ExecutionDoneEvent[R]) {
			ctx := context.Background()

			p.Metrics.attemptCounter.Add(ctx, int64(ede.Attempts()), p.nameAttributeMeasurmentOption)

			p.Metrics.executionCounter.Add(ctx, int64(ede.Executions()), p.nameAttributeMeasurmentOption)
			p.Metrics.executionDuration.Record(ctx, float64(ede.ElapsedTime()*time.Millisecond), p.nameAttributeMeasurmentOption)

			p.Metrics.retryCounter.Add(ctx, int64(ede.Retries()), p.nameAttributeMeasurmentOption)
			p.Metrics.hedgeCounter.Add(ctx, int64(ede.Hedges()), p.nameAttributeMeasurmentOption)

			if success {
				p.Metrics.successCounter.Add(ctx, 1, p.nameAttributeMeasurmentOption)
			} else {
				p.Metrics.failureCounter.Add(ctx, 1, p.nameAttributeMeasurmentOption)
			}
		}
	}

	return failsafe.With(policies...).
		OnSuccess(metricRegister(true)).
		OnFailure(metricRegister(false))
}
