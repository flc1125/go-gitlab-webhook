package otel

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type config struct {
	tracerProvider trace.TracerProvider
}

// Option configures the OpenTelemetry middleware.
type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (fn optionFunc) apply(cfg *config) {
	fn(cfg)
}

// WithTracerProvider configures the tracer provider used by the middleware.
func WithTracerProvider(provider trace.TracerProvider) Option {
	return optionFunc(func(c *config) {
		if provider != nil {
			c.tracerProvider = provider
		}
	})
}

func newConfig(opts ...Option) *config {
	cfg := &config{
		tracerProvider: otel.GetTracerProvider(),
	}

	for _, opt := range opts {
		opt.apply(cfg)
	}

	return cfg
}
