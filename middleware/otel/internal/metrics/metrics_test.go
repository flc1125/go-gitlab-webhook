package metrics_test

import (
	"context"
	"testing"
	"time"

	"github.com/flc1125/go-gitlab-webhook/middleware/otel/v3/internal/eventmeta"
	"github.com/flc1125/go-gitlab-webhook/middleware/otel/v3/internal/metrics"
	"github.com/flc1125/go-gitlab-webhook/middleware/otel/v3/internal/semconv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func TestRecord(t *testing.T) {
	reader := sdkmetric.NewManualReader()
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	recorder := metrics.New(provider, "test/instrumentation", "v0.0.0")
	metadata := eventmeta.Metadata{
		SpanName: "gitlab.webhook.push",
		Attributes: []attribute.KeyValue{
			semconv.WebhookEventType("push"),
			semconv.WebhookObjectKind("push"),
			semconv.ProjectPath("group/project"),
			semconv.Ref("refs/heads/main"),
		},
	}

	recorder.RecordActive(context.Background(), metadata, 1)
	recorder.Record(context.Background(), metadata, time.Second, nil)
	recorder.RecordActive(context.Background(), metadata, -1)

	rm := collectMetrics(t, reader)
	events := metricByName(t, rm, "gitlab.webhook.events")
	eventSum, ok := events.Data.(metricdata.Sum[int64])
	require.True(t, ok)
	require.Len(t, eventSum.DataPoints, 1)
	assert.Equal(t, int64(1), eventSum.DataPoints[0].Value)
	assertAttrString(t, eventSum.DataPoints[0].Attributes, "gitlab.webhook.event_type", "push")
	assertAttrString(t, eventSum.DataPoints[0].Attributes, "gitlab.webhook.object_kind", "push")
	assertAttrString(t, eventSum.DataPoints[0].Attributes, "gitlab.webhook.result", "success")
	assertNoAttr(t, eventSum.DataPoints[0].Attributes, "gitlab.project.path")
	assertNoAttr(t, eventSum.DataPoints[0].Attributes, "gitlab.ref")

	duration := metricByName(t, rm, "gitlab.webhook.event.duration")
	histogram, ok := duration.Data.(metricdata.Histogram[float64])
	require.True(t, ok)
	require.Len(t, histogram.DataPoints, 1)
	assert.Equal(t, uint64(1), histogram.DataPoints[0].Count)

	active := metricByName(t, rm, "gitlab.webhook.active_events")
	activeSum, ok := active.Data.(metricdata.Sum[int64])
	require.True(t, ok)
	require.Len(t, activeSum.DataPoints, 1)
	assert.Equal(t, int64(0), activeSum.DataPoints[0].Value)
	assertAttrString(t, activeSum.DataPoints[0].Attributes, "gitlab.webhook.event_type", "push")
	assertNoAttr(t, activeSum.DataPoints[0].Attributes, "gitlab.webhook.result")
}

func collectMetrics(t *testing.T, reader *sdkmetric.ManualReader) metricdata.ResourceMetrics {
	t.Helper()

	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(t.Context(), &rm))

	return rm
}

func metricByName(t *testing.T, rm metricdata.ResourceMetrics, name string) metricdata.Metrics {
	t.Helper()

	for _, scope := range rm.ScopeMetrics {
		for _, metric := range scope.Metrics {
			if metric.Name == name {
				return metric
			}
		}
	}

	require.Failf(t, "metric not found", "metric %q was not collected", name)
	return metricdata.Metrics{}
}

func assertAttrString(t *testing.T, attrs attribute.Set, key, expected string) {
	t.Helper()

	value, ok := attrs.Value(attribute.Key(key))
	require.Truef(t, ok, "attribute %q was not found", key)
	assert.Equal(t, expected, value.AsString())
}

func assertNoAttr(t *testing.T, attrs attribute.Set, key string) {
	t.Helper()

	_, ok := attrs.Value(attribute.Key(key))
	assert.Falsef(t, ok, "attribute %q should not be recorded", key)
}
