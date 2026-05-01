package otel

import (
	"context"
	"reflect"

	gitlabwebhook "github.com/flc1125/go-gitlab-webhook/v3"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	instrumentationName = "github.com/flc1125/go-gitlab-webhook/middleware/otel/v3"

	attributeEventType  = "gitlab.webhook.event_type"
	attributeObjectKind = "gitlab.webhook.object_kind"
	attributeEventName  = "gitlab.webhook.event_name"
)

// Middleware returns middleware that traces webhook event dispatch.
func Middleware(opts ...Option) gitlabwebhook.Middleware {
	cfg := newConfig(opts...)
	tracer := cfg.tracerProvider.Tracer(
		instrumentationName,
		trace.WithInstrumentationVersion(Version()),
	)

	return func(next gitlabwebhook.HandlerFunc) gitlabwebhook.HandlerFunc {
		return func(ctx context.Context, event any) error {
			ctx, span := tracer.Start(ctx,
				defaultSpanName(event),
				trace.WithAttributes(spanAttributes(event)...),
			)
			defer span.End()

			err := next(ctx, event)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}

			return err
		}
	}
}

func spanAttributes(event any) []attribute.KeyValue {
	attributes := []attribute.KeyValue{
		attribute.String(attributeEventType, eventTypeName(event)),
	}
	attributes = append(attributes, eventAttributes(event)...)

	return attributes
}

func defaultSpanName(event any) string {
	return "gitlab.webhook " + eventTypeName(event)
}

func eventTypeName(event any) string {
	t := eventType(event)
	if t == nil {
		return ""
	}

	if t.Name() == "" {
		return t.String()
	}

	return t.Name()
}

func eventAttributes(event any) []attribute.KeyValue {
	v := eventValue(event)
	if !v.IsValid() {
		return nil
	}

	attributes := make([]attribute.KeyValue, 0, 2)
	if value, ok := stringField(v, "ObjectKind"); ok {
		attributes = append(attributes, attribute.String(attributeObjectKind, value))
	}
	if value, ok := stringField(v, "EventName"); ok {
		attributes = append(attributes, attribute.String(attributeEventName, value))
	}

	return attributes
}

func eventType(event any) reflect.Type {
	if event == nil {
		return nil
	}

	t := reflect.TypeOf(event)
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	return t
}

func eventValue(event any) reflect.Value {
	if event == nil {
		return reflect.Value{}
	}

	v := reflect.ValueOf(event)
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return reflect.Value{}
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return reflect.Value{}
	}

	return v
}

func stringField(v reflect.Value, name string) (string, bool) {
	field := v.FieldByName(name)
	if !field.IsValid() || field.Kind() != reflect.String {
		return "", false
	}

	value := field.String()
	if value == "" {
		return "", false
	}

	return value, true
}
