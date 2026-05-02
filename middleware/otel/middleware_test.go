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
