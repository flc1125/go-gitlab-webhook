# Gitlab Webhook Dispatcher 🚀

![Supported Go Versions](https://img.shields.io/badge/Go-%3E%3D1.25.0-blue)
[![Package Version](https://badgen.net/github/release/flc1125/go-gitlab-webhook/stable)](https://github.com/flc1125/go-gitlab-webhook/releases)
[![GoDoc](https://pkg.go.dev/badge/github.com/flc1125/go-gitlab-webhook/v3)](https://pkg.go.dev/github.com/flc1125/go-gitlab-webhook/v3)
[![codecov](https://codecov.io/gh/flc1125/go-gitlab-webhook/graph/badge.svg?token=QPTHZ5L9GT)](https://codecov.io/gh/flc1125/go-gitlab-webhook)
[![Go Report Card](https://goreportcard.com/badge/github.com/flc1125/go-gitlab-webhook)](https://goreportcard.com/report/github.com/flc1125/go-gitlab-webhook)
[![lint](https://github.com/flc1125/go-gitlab-webhook/actions/workflows/lint.yml/badge.svg)](https://github.com/flc1125/go-gitlab-webhook/actions/workflows/lint.yml)
[![tests](https://github.com/flc1125/go-gitlab-webhook/actions/workflows/test.yml/badge.svg)](https://github.com/flc1125/go-gitlab-webhook/actions/workflows/test.yml)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)

This is a simple webhook dispatcher for Gitlab. It listens for incoming webhooks and dispatches them to the appropriate handler.

## ✨ Features

- 📋 Very convenient registration of listeners
- 🔄 A single listener can implement multiple different webhook functions
- ⚡ Support asynchronous and efficient processing
- 🚀 Multiple dispatch methods
- 🔐 Token validation support for secure webhook handling

## 📦 Installation

```shell
go get github.com/flc1125/go-gitlab-webhook/v3
```

## 🔗 Compatibility

| go-gitlab-webhook | GitLab Go Client |
| --- | --- |
| `1.x` | [xanzy/go-gitlab](https://github.com/xanzy/go-gitlab) |
| `2.x` | [gitlab-org/api/client-go](https://gitlab.com/gitlab-org/api/client-go) |
| `3.x` | [gitlab-org/api/client-go/v2](https://gitlab.com/gitlab-org/api/client-go/v2) |

## 💻 Usage

```go
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/flc1125/go-gitlab-webhook/v3"
	"gitlab.com/gitlab-org/api/client-go/v2"
)

var (
	_ gitlabwebhook.BuildListener         = (*testBuildListener)(nil)
	_ gitlabwebhook.CommitCommentListener = (*testCommitCommentListener)(nil)
	_ gitlabwebhook.BuildListener         = (*testBuildAndCommitCommentListener)(nil)
	_ gitlabwebhook.CommitCommentListener = (*testBuildAndCommitCommentListener)(nil)
)

type testBuildListener struct{}

func (l *testBuildListener) OnBuild(ctx context.Context, event *gitlab.BuildEvent) error {
	// do something
	return nil
}

type testCommitCommentListener struct{}

func (l *testCommitCommentListener) OnCommitComment(ctx context.Context, event *gitlab.CommitCommentEvent) error {
	// do something
	return nil
}

type testBuildAndCommitCommentListener struct{}

func (l *testBuildAndCommitCommentListener) OnBuild(ctx context.Context, event *gitlab.BuildEvent) error {
	// do something
	return nil
}

func (l *testBuildAndCommitCommentListener) OnCommitComment(ctx context.Context, event *gitlab.CommitCommentEvent) error {
	// do something
	return nil
}

func main() {
	dispatcher := gitlabwebhook.NewDispatcher(
		gitlabwebhook.RegisterListeners(
			&testBuildListener{},
			&testCommitCommentListener{},
			&testBuildAndCommitCommentListener{},
		),
		gitlabwebhook.WithMiddlewares(
			loggingMiddleware,
			// Only runs for push events.
			gitlabwebhook.MiddlewareForEvent(func(ctx context.Context, event *gitlab.PushEvent) error {
				log.Printf("push event: %s", event.Project.PathWithNamespace)
				return nil
			}),
		),
	)

	dispatcher.RegisterListeners(
		&testBuildListener{},
		&testCommitCommentListener{},
		&testBuildAndCommitCommentListener{},
	)

	http.Handle("/webhook", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := dispatcher.DispatchRequest(r,
			gitlabwebhook.DispatchRequestWithToken("your-secret-token"), // validate token, if needed
			gitlabwebhook.DispatchRequestWithContext(context.Background()), // custom context
		); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func loggingMiddleware(next gitlabwebhook.HandlerFunc) gitlabwebhook.HandlerFunc {
	return func(ctx context.Context, event any) error {
		log.Printf("received webhook event: %T", event)
		return next(ctx, event)
	}
}
```

## 📜 License

MIT License. See [LICENSE](LICENSE) for the full license text.

## 💖 Thanks

- [gitlab-org/api/client-go](https://gitlab.com/gitlab-org/api/client-go)(Formerly known as: [xanzy/go-gitlab](https://github.com/xanzy/go-gitlab)): Go library for accessing the GitLab API
- [stretchr/testify](github.com/stretchr/testify): A toolkit with common assertions and mocks that plays nicely with the standard library
