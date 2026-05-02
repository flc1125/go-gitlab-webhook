package metrics

import (
	"context"
	"time"

	"github.com/flc1125/go-gitlab-webhook/middleware/otel/v3/internal/eventmeta"
	"github.com/flc1125/go-gitlab-webhook/middleware/otel/v3/internal/semconv"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
)

const (
	webhookResultSuccess = "success"
	webhookResultError   = "error"
)

// Metrics records GitLab webhook event metrics.
type Metrics struct {
	events   metric.Int64Counter
	active   metric.Int64UpDownCounter
	duration metric.Float64Histogram
}

// New creates the OpenTelemetry instruments used by the middleware.
func New(provider metric.MeterProvider, instrumentationName, instrumentationVersion string) Metrics {
	meter := provider.Meter(
		instrumentationName,
		metric.WithInstrumentationVersion(instrumentationVersion),
	)

	events, err := meter.Int64Counter(
		"gitlab.webhook.events",
		metric.WithDescription("Total number of GitLab webhook events handled."),
		metric.WithUnit("{event}"),
	)
	if err != nil {
		events = noop.Int64Counter{}
	}

	active, err := meter.Int64UpDownCounter(
		"gitlab.webhook.active_events",
		metric.WithDescription("Number of GitLab webhook events currently being handled."),
		metric.WithUnit("{event}"),
	)
	if err != nil {
		active = noop.Int64UpDownCounter{}
	}

	duration, err := meter.Float64Histogram(
		"gitlab.webhook.event.duration",
		metric.WithDescription("Duration of GitLab webhook event handling."),
		metric.WithUnit("s"),
	)
	if err != nil {
		duration = noop.Float64Histogram{}
	}

	return Metrics{
		events:   events,
		active:   active,
		duration: duration,
	}
}

// RecordActive records the current number of webhook events being handled.
func (m Metrics) RecordActive(ctx context.Context, metadata eventmeta.Metadata, delta int64) {
	m.active.Add(ctx, delta, metric.WithAttributes(metricAttributes(metadata.Attributes)...))
}

// Record records the handling result and duration for a webhook event.
func (m Metrics) Record(ctx context.Context, metadata eventmeta.Metadata, elapsed time.Duration, err error) {
	result := webhookResultSuccess
	if err != nil {
		result = webhookResultError
	}

	attrs := append(metricAttributes(metadata.Attributes), semconv.WebhookResult(result))
	opts := metric.WithAttributes(attrs...)

	m.events.Add(ctx, 1, opts)
	m.duration.Record(ctx, elapsed.Seconds(), opts)
}

func metricAttributes(attrs []attribute.KeyValue) []attribute.KeyValue {
	values := make([]attribute.KeyValue, 0, len(attrs))
	for _, attr := range attrs {
		switch attr.Key {
		case semconv.WebhookEventTypeKey,
			semconv.WebhookObjectKindKey,
			semconv.WebhookEventNameKey,
			semconv.WebhookActionKey,
			semconv.WebhookStatusKey:
			values = append(values, attr)
		}
	}

	return values
}
