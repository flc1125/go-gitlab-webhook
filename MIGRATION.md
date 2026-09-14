# Migrating from v3 to v4

Version 4 uses GitLab Go Client v3. Use Go 1.26.0 or later, as required by both libraries.

## Update dependencies and imports

Update your application imports using this mapping:

| v3 application imports | v4 application imports |
| --- | --- |
| `github.com/flc1125/go-gitlab-webhook/v3` | `github.com/flc1125/go-gitlab-webhook/v4` |
| `github.com/flc1125/go-gitlab-webhook/middleware/otel/v3` | `github.com/flc1125/go-gitlab-webhook/middleware/otel/v4` |
| `gitlab.com/gitlab-org/api/client-go/v2` | `gitlab.com/gitlab-org/api/client-go/v3` |

Once v4.0.0 is published, install the new modules:

```shell
go get github.com/flc1125/go-gitlab-webhook/v4@v4.0.0
```

If you use the OpenTelemetry middleware, also run:

```shell
go get github.com/flc1125/go-gitlab-webhook/middleware/otel/v4@v4.0.0
```

Update all listener method signatures and typed middleware to use event types from `client-go/v3`. Go treats types from `/v2` and `/v3` as different types; a listener using `/v2` events does not implement the v4 listener interface. Add compile-time interface assertions to catch missed imports:

```go
package webhooks

import (
	"context"

	gitlabwebhook "github.com/flc1125/go-gitlab-webhook/v4"
	gitlab "gitlab.com/gitlab-org/api/client-go/v3"
)

type mergeListener struct{}

var _ gitlabwebhook.MergeListener = (*mergeListener)(nil)

func (l *mergeListener) OnMerge(ctx context.Context, event *gitlab.MergeEvent) error {
	return nil
}
```

Dispatcher methods, listener method names, registration options, token validation, payload limits, and middleware composition retain their existing usage.

## Update merge request reviewer types

GitLab Go Client v3 changes these fields from `[]*gitlab.EventUser` to `[]*gitlab.EventReviewer`:

- `MergeEvent.Reviewers`
- `MergeEventChangesReviewers.Previous`
- `MergeEventChangesReviewers.Current`

Update helper function parameters and struct literals that explicitly use the old reviewer type. `EventReviewer` includes the existing user fields and adds `State` and `ReRequested`. These fields are available to listeners and typed middleware after webhook parsing.

If your application also uses the GitLab REST client, check the [upstream v3 migration guide](https://gitlab.com/gitlab-org/api/client-go/-/blob/v3.9.0/docs/release-3.0-migration.md) for changes to other APIs.

## Update OpenTelemetry scope filters

The middleware instrumentation scope name changes to `github.com/flc1125/go-gitlab-webhook/middleware/otel/v4`, and its instrumentation version becomes `4.0.0`. Update SDK views, collector filters, or dashboards that match the old scope name or version.

Span names, metric names, and semantic attribute keys retain their existing names.

## Verify your application

Run dependency cleanup and your application's tests in each Go module:

```shell
go mod tidy
go test -race ./...
```

Applications staying on `go-gitlab-webhook/v3` continue to use `client-go/v2`.
