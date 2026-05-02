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

// Middleware returns middleware that traces webhook event dispatch.
func Middleware(opts ...Option) gitlabwebhook.Middleware {
	cfg := newConfig(opts...)
	tracer := cfg.tracerProvider.Tracer(
		instrumentationName,
		trace.WithInstrumentationVersion(Version()),
	)
	metricRecorder := metrics.New(cfg.meterProvider, instrumentationName, Version())

	return func(next gitlabwebhook.HandlerFunc) gitlabwebhook.HandlerFunc {
		return func(ctx context.Context, event any) error {
			metadata := eventmeta.Extract(event)
			start := time.Now()
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
