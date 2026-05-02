// Package otel provides OpenTelemetry middleware for GitLab webhook dispatch.
//
// Use [Middleware] to add tracing and metrics to a [gitlabwebhook.Dispatcher].
// Tracer and meter providers can be configured with [WithTracerProvider] and
// [WithMeterProvider].
package otel
