package otel

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type config struct {
	tracerProvider trace.TracerProvider
	meterProvider  metric.MeterProvider
}

// Option configures [Middleware].
//
// Nil options are ignored.
type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (fn optionFunc) apply(cfg *config) {
	fn(cfg)
}

// WithTracerProvider configures the tracer provider used by the middleware.
//
// If provider is nil, [Middleware] keeps using the global OpenTelemetry
// tracer provider.
func WithTracerProvider(provider trace.TracerProvider) Option {
	return optionFunc(func(c *config) {
		if provider != nil {
			c.tracerProvider = provider
		}
	})
}

// WithMeterProvider configures the meter provider used by the middleware.
//
// If provider is nil, [Middleware] keeps using the global OpenTelemetry meter
// provider.
func WithMeterProvider(provider metric.MeterProvider) Option {
	return optionFunc(func(c *config) {
		if provider != nil {
			c.meterProvider = provider
		}
	})
}

func newConfig(opts ...Option) *config {
	cfg := &config{
		tracerProvider: otel.GetTracerProvider(),
		meterProvider:  otel.GetMeterProvider(),
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt.apply(cfg)
	}

	return cfg
}
