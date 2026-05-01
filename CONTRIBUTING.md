# Contributing

Thank you for contributing to this project.

## Development

Before submitting changes, run:

```bash
make lint
make test
```

## Adding Support for a New Event

If the upstream GitLab Go client already supports parsing a webhook event, but this repository's dispatcher does not yet route it, follow the steps below to add support.

### 1. Confirm upstream parsing support first

Before making changes in this repository, confirm that the upstream `ParseWebhook(...)` already supports the event:

- Docs: https://pkg.go.dev/gitlab.com/gitlab-org/api/client-go/v2
- Source: https://gitlab.com/gitlab-org/api/client-go/-/blob/main/event_parsing.go

Dispatcher support in this repository should only be added after the upstream client can already parse the webhook payload into the corresponding event struct.

### 2. Add the listener interface

Add the corresponding listener interface in `listeners.go`.

Example:

```go
type EmojiListener interface {
	OnEmoji(ctx context.Context, event *gitlab.EmojiEvent) error
}
```

Naming conventions:

- use `XxxListener` for the listener interface
- use `OnXxx` for the handler method
- use the upstream event struct type from the GitLab Go client

### 3. Register the event in the dispatcher

Update `dispatcher.go` to include the full dispatcher flow:

- add a listener slice field to `Dispatcher`
- add `RegisterXxxListener(...)`
- add auto-registration support in `RegisterListeners(...)`
- add a type branch in `Dispatch(...)`
- add `processXxxEvent(...)`

Follow the existing event implementations in this repository to keep naming and structure consistent.

### 4. Add tests

Extend `dispatcher_test.go` to cover the new event. At minimum, tests should verify:

- the webhook payload is dispatched correctly through `DispatchRequest(...)`
- the listener is actually invoked
- key event fields are asserted as expected

If there is no existing fixture for the event, add a minimal valid payload under `internal/testdata/webhooks/`.

### 5. Scope of a typical event support change

A complete event support change usually includes:

- listener definition
- dispatcher registration and routing
- tests
- a webhook fixture, if needed

Prefer keeping one PR focused on one event. That makes review and regression validation much easier.

### Reference Implementations

You can use existing supported events as references, for example:

- `BuildEvent`
- `ReleaseEvent`
- `WikiPageEvent`
- `EmojiEvent`
