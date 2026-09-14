package gitlabwebhook_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	gitlabwebhook "github.com/flc1125/go-gitlab-webhook/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gitlab "gitlab.com/gitlab-org/api/client-go/v3"
)

type mergeReviewerListener struct {
	event *gitlab.MergeEvent
}

var _ gitlabwebhook.MergeListener = (*mergeReviewerListener)(nil)

func (l *mergeReviewerListener) OnMerge(_ context.Context, event *gitlab.MergeEvent) error {
	l.event = event
	return nil
}

func TestDispatcherMergeReviewers(t *testing.T) {
	payload, err := os.ReadFile("testdata/webhooks/merge_request_reviewers.json")
	require.NoError(t, err)

	tests := []struct {
		name     string
		dispatch func(*testing.T, *gitlabwebhook.Dispatcher) error
	}{
		{
			name: "Dispatch",
			dispatch: func(t *testing.T, d *gitlabwebhook.Dispatcher) error {
				event, err := gitlab.ParseWebhook(gitlab.EventTypeMergeRequest, payload)
				require.NoError(t, err)
				return d.Dispatch(t.Context(), event)
			},
		},
		{
			name: "DispatchWebhook",
			dispatch: func(t *testing.T, d *gitlabwebhook.Dispatcher) error {
				return d.DispatchWebhook(t.Context(), gitlab.EventTypeMergeRequest, payload)
			},
		},
		{
			name: "DispatchRequest",
			dispatch: func(t *testing.T, d *gitlabwebhook.Dispatcher) error {
				req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/webhook", bytes.NewReader(payload))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("X-Gitlab-Event", string(gitlab.EventTypeMergeRequest))
				return d.DispatchRequest(req)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			listener := &mergeReviewerListener{}
			var middlewareEvent *gitlab.MergeEvent
			d := gitlabwebhook.NewDispatcher(
				gitlabwebhook.RegisterListeners(listener),
				gitlabwebhook.WithMiddlewares(
					gitlabwebhook.MiddlewareForEvent(func(_ context.Context, event *gitlab.MergeEvent) error {
						middlewareEvent = event
						return nil
					}),
				),
			)

			require.NoError(t, tt.dispatch(t, d))
			require.NotNil(t, listener.event)
			assert.Same(t, listener.event, middlewareEvent)
			assert.Equal(t, "update", listener.event.ObjectAttributes.Action)
			assert.Equal(t, []*gitlab.EventReviewer{
				{ID: 42, Name: "Reviewer", Username: "reviewer", State: "unreviewed", ReRequested: true},
			}, listener.event.Reviewers)
			assert.Equal(t, []*gitlab.EventReviewer{
				{ID: 42, Name: "Reviewer", Username: "reviewer", State: "approved", ReRequested: false},
			}, listener.event.Changes.Reviewers.Previous)
			assert.Equal(t, listener.event.Reviewers, listener.event.Changes.Reviewers.Current)
		})
	}
}
