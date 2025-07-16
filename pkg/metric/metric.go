package metric

import (
	"context"

	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

type noopMetricExporter struct{}

func (noopMetricExporter) Temporality(kind metric.InstrumentKind) metricdata.Temporality {
	return metric.DefaultTemporalitySelector(kind)
}

func (noopMetricExporter) Aggregation(kind metric.InstrumentKind) metric.Aggregation {
	return metric.DefaultAggregationSelector(kind)
}

func (noopMetricExporter) Export(context.Context, *metricdata.ResourceMetrics) error {
	return nil
}

func (noopMetricExporter) ForceFlush(context.Context) error { return nil }

func (noopMetricExporter) Shutdown(context.Context) error { return nil }

// NewNoopExporter returns a noop metric exporter.
func NewNoopExporter() metric.Exporter {
	return noopMetricExporter{}
}
