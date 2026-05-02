package otel_test

import (
	"context"
	"errors"
	"testing"

	"github.com/flc1125/go-gitlab-webhook/middleware/otel/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gitlab "gitlab.com/gitlab-org/api/client-go/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

func TestMiddlewareCreatesSpan(t *testing.T) {
	recorder := tracetest.NewSpanRecorder()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))

	event := &gitlab.PushEvent{
		ObjectKind:        "push",
		EventName:         "push",
		ProjectID:         123,
		Ref:               "refs/heads/main",
		CheckoutSHA:       "abc123",
		TotalCommitsCount: 2,
		Project: gitlab.PushEventProject{
			Name:              "project",
			PathWithNamespace: "group/project",
		},
	}

	called := false
	handler := otel.Middleware(otel.WithTracerProvider(provider))(func(ctx context.Context, event any) error {
		called = true
		assert.True(t, trace.SpanFromContext(ctx).SpanContext().IsValid())
		return nil
	})

	err := handler(context.Background(), event)
	require.NoError(t, err)
	assert.True(t, called)

	spans := recorder.Ended()
	require.Len(t, spans, 1)

	span := tracetest.SpanStubFromReadOnlySpan(spans[0])
	assert.Equal(t, "gitlab.webhook.push", span.Name)
	assert.Equal(t, "push", attrString(span.Attributes, "gitlab.webhook.event_type"))
	assert.Equal(t, "group/project", attrString(span.Attributes, "gitlab.project.path"))
	assert.Equal(t, int64(2), attrInt64(span.Attributes, "gitlab.push.total_commits_count"))
}

func TestMiddlewareRecordsError(t *testing.T) {
	recorder := tracetest.NewSpanRecorder()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
	expectedErr := errors.New("listener failed")

	handler := otel.Middleware(otel.WithTracerProvider(provider))(func(context.Context, any) error {
		return expectedErr
	})

	err := handler(context.Background(), &gitlab.PushEvent{})
	require.ErrorIs(t, err, expectedErr)

	spans := recorder.Ended()
	require.Len(t, spans, 1)

	span := tracetest.SpanStubFromReadOnlySpan(spans[0])
	assert.Equal(t, codes.Error, span.Status.Code)
	assert.Equal(t, expectedErr.Error(), span.Status.Description)
	assert.NotEmpty(t, span.Events)
}

func TestMiddlewareSkipsNilOption(t *testing.T) {
	recorder := tracetest.NewSpanRecorder()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))

	handler := otel.Middleware(nil, otel.WithTracerProvider(provider))(func(context.Context, any) error {
		return nil
	})

	err := handler(context.Background(), &gitlab.PushEvent{})
	require.NoError(t, err)
	assert.Len(t, recorder.Ended(), 1)
}

func TestMiddlewareRecordsMetrics(t *testing.T) {
	reader := sdkmetric.NewManualReader()
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	event := &gitlab.PushEvent{
		ObjectKind: "push",
		EventName:  "push",
		ProjectID:  123,
		Ref:        "refs/heads/main",
		Project: gitlab.PushEventProject{
			PathWithNamespace: "group/project",
		},
	}

	handler := otel.Middleware(otel.WithMeterProvider(provider))(func(context.Context, any) error {
		return nil
	})

	err := handler(context.Background(), event)
	require.NoError(t, err)

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
}

func attrString(attrs []attribute.KeyValue, key string) string {
	for _, attr := range attrs {
		if string(attr.Key) == key {
			return attr.Value.AsString()
		}
	}

	return ""
}

func attrInt64(attrs []attribute.KeyValue, key string) int64 {
	for _, attr := range attrs {
		if string(attr.Key) == key {
			return attr.Value.AsInt64()
		}
	}

	return 0
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
