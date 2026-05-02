# GitLab Webhook OpenTelemetry Middleware

OpenTelemetry tracing and metrics middleware for `github.com/flc1125/go-gitlab-webhook/v3`.

## Installation

```shell
go get github.com/flc1125/go-gitlab-webhook/middleware/otel/v3
```

## Usage

```go
package main

import (
	"context"
	"net/http"

	otelmiddleware "github.com/flc1125/go-gitlab-webhook/middleware/otel/v3"
	gitlabwebhook "github.com/flc1125/go-gitlab-webhook/v3"
	"gitlab.com/gitlab-org/api/client-go/v2"
)

type pushListener struct{}

func (l *pushListener) OnPush(ctx context.Context, event *gitlab.PushEvent) error {
	return nil
}

func main() {
	dispatcher := gitlabwebhook.NewDispatcher(
		gitlabwebhook.RegisterListeners(
			&pushListener{},
		),
		gitlabwebhook.WithMiddlewares(
			otelmiddleware.Middleware(),
		),
	)

	http.Handle("/webhook", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := dispatcher.DispatchRequest(r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}))
}
```

## Tracer Provider

By default, the middleware uses the global OpenTelemetry tracer and meter providers.
Pass providers explicitly when your application wires telemetry without globals:

```go
dispatcher := gitlabwebhook.NewDispatcher(
	gitlabwebhook.WithMiddlewares(
		otelmiddleware.Middleware(
			otelmiddleware.WithTracerProvider(tracerProvider),
			otelmiddleware.WithMeterProvider(meterProvider),
		),
	),
)
```

## Spans

The middleware creates one span per parsed webhook event. Listener execution is
included in that event span.

Span names use stable, low-cardinality values:

```text
gitlab.webhook.push
gitlab.webhook.merge_request open
gitlab.webhook.pipeline success
```

Event details are recorded as attributes, for example:

- `gitlab.webhook.event_type`
- `gitlab.webhook.object_kind`
- `gitlab.webhook.action`
- `gitlab.webhook.status`
- `gitlab.project.id`
- `gitlab.project.path`
- `gitlab.ref`
- `gitlab.sha`

## Metrics

The middleware records basic event handling metrics:

- `gitlab.webhook.events`: count of handled webhook events
- `gitlab.webhook.event.duration`: event handling duration in seconds

Metric attributes are intentionally low-cardinality:

- `gitlab.webhook.event_type`
- `gitlab.webhook.object_kind`
- `gitlab.webhook.event_name`
- `gitlab.webhook.action`
- `gitlab.webhook.status`
- `gitlab.webhook.result`
