package otel

import (
	"context"
	"time"

	"github.com/flc1125/go-gitlab-webhook/middleware/otel/v3/internal/eventmeta"
	"github.com/flc1125/go-gitlab-webhook/middleware/otel/v3/internal/metrics"
	gitlabwebhook "github.com/flc1125/go-gitlab-webhook/v3"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	instrumentationName = "github.com/flc1125/go-gitlab-webhook/middleware/otel/v3"
)

// Middleware returns a [gitlabwebhook.Middleware] that traces and records
// metrics for webhook event dispatch.
//
// The middleware creates a span for each dispatched event, records errors on the
// span, and records event metrics with the configured meter provider.
func Middleware(opts ...Option) gitlabwebhook.Middleware {
	cfg := newConfig(opts...)
	tracer := cfg.tracerProvider.Tracer(
		instrumentationName,
		trace.WithInstrumentationVersion(Version()),
	)
	metricRecorder := metrics.New(cfg.meterProvider, instrumentationName, Version())

	return func(next gitlabwebhook.HandlerFunc) gitlabwebhook.HandlerFunc {
		return func(ctx context.Context, event any) error {
			start := time.Now()
			metadata := eventmeta.Extract(event)
			metricRecorder.RecordActive(ctx, metadata, 1)
			defer metricRecorder.RecordActive(ctx, metadata, -1)

			ctx, span := tracer.Start(ctx,
				metadata.SpanName,
				trace.WithAttributes(metadata.Attributes...),
			)
			defer span.End()

			err := next(ctx, event)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}

			metricRecorder.Record(ctx, metadata, time.Since(start), err)

			return err
		}
	}
}
