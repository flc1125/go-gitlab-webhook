package gitlabwebhook

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	gitlab "gitlab.com/gitlab-org/api/client-go/v2"
)

func TestDispatcher_DispatchRequest_Fixtures(t *testing.T) {
	tests := []struct {
		name      string
		eventType gitlab.EventType
		filepath  string
	}{
		{name: "build", eventType: gitlab.EventTypeBuild, filepath: "testdata/webhooks/build.json"},
		{name: "commit comment", eventType: gitlab.EventTypeNote, filepath: "testdata/webhooks/note_commit.json"},
		{name: "deployment", eventType: gitlab.EventTypeDeployment, filepath: "testdata/webhooks/deployment.json"},
		{name: "emoji", eventType: gitlab.EventTypeEmoji, filepath: "testdata/webhooks/emoji.json"},
		{name: "feature flag", eventType: gitlab.EventTypeFeatureFlag, filepath: "testdata/webhooks/feature_flag.json"},
		{name: "group resource access token", eventType: gitlab.EventTypeResourceAccessToken, filepath: "testdata/webhooks/resource_access_token_group.json"},
		{name: "issue comment", eventType: gitlab.EventTypeNote, filepath: "testdata/webhooks/note_issue.json"},
		{name: "issue", eventType: gitlab.EventTypeIssue, filepath: "testdata/webhooks/issue.json"},
		{name: "job", eventType: gitlab.EventTypeJob, filepath: "testdata/webhooks/job.json"},
		{name: "member", eventType: gitlab.EventTypeMember, filepath: "testdata/webhooks/member.json"},
		{name: "milestone", eventType: gitlab.EventTypeMilestone, filepath: "testdata/webhooks/milestone.json"},
		{name: "merge comment", eventType: gitlab.EventTypeNote, filepath: "testdata/webhooks/note_merge_request.json"},
		{name: "merge", eventType: gitlab.EventTypeMergeRequest, filepath: "testdata/webhooks/merge_request.json"},
		{name: "pipeline", eventType: gitlab.EventTypePipeline, filepath: "testdata/webhooks/pipeline.json"},
		{name: "project", eventType: gitlab.EventTypeProject, filepath: "testdata/webhooks/project.json"},
		{name: "project resource access token", eventType: gitlab.EventTypeResourceAccessToken, filepath: "testdata/webhooks/resource_access_token_project.json"},
		{name: "push", eventType: gitlab.EventTypePush, filepath: "testdata/webhooks/push.json"},
		{name: "release", eventType: gitlab.EventTypeRelease, filepath: "testdata/webhooks/release.json"},
		{name: "snippet comment", eventType: gitlab.EventTypeNote, filepath: "testdata/webhooks/note_snippet.json"},
		{name: "subgroup", eventType: gitlab.EventTypeSubGroup, filepath: "testdata/webhooks/subgroup.json"},
		{name: "tag", eventType: gitlab.EventTypeTagPush, filepath: "testdata/webhooks/tag_push.json"},
		{name: "vulnerability", eventType: gitlab.EventTypeVulnerability, filepath: "testdata/webhooks/vulnerability.json"},
		{name: "wiki page", eventType: gitlab.EventTypeWikiPage, filepath: "testdata/webhooks/wiki_page.json"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			listener := &eventAssertingListener{}
			dispatcher := NewDispatcher(RegisterListeners(listener))
			dispatcher.RegisterListeners(listener)

			req := newWebhookRequest(t, tt.eventType, loadFixture(t, tt.filepath))

			err := dispatcher.DispatchRequest(
				req,
				DispatchRequestWithContext(newDispatcherContext(context.Background())),
			)

			assert.NoError(t, err)
			assert.EqualValues(t, 2, listener.calls.Load())
		})
	}
}

func TestDispatcher_DispatchRequestWithToken(t *testing.T) {
	listener := &pushTrackingListener{}
	dispatcher := NewDispatcher(RegisterListeners(listener))

	validToken := "test-secret-token"
	invalidToken := "wrong-token"

	tests := []struct {
		name          string
		token         string
		headerToken   string
		expectedError error
	}{
		{
			name:          "valid token should dispatch successfully",
			token:         validToken,
			headerToken:   validToken,
			expectedError: nil,
		},
		{
			name:          "invalid token should return ErrInvalidToken",
			token:         validToken,
			headerToken:   invalidToken,
			expectedError: ErrInvalidToken,
		},
		{
			name:          "missing token header should return ErrInvalidToken",
			token:         validToken,
			headerToken:   "",
			expectedError: ErrInvalidToken,
		},
		{
			name:          "no token provided should dispatch successfully",
			token:         "",
			headerToken:   "",
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := newWebhookRequest(t, gitlab.EventTypePush, loadFixture(t, "testdata/webhooks/push.json"))

			if tt.headerToken != "" {
				req.Header.Set("X-Gitlab-Token", tt.headerToken)
			}

			var opts []DispatchRequestOption
			if tt.token != "" {
				opts = append(opts, DispatchRequestWithToken(tt.token))
			}

			listener.called.Store(0)

			err := dispatcher.DispatchRequest(req, opts...)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Zero(t, listener.called.Load())
				return
			}

			assert.NoError(t, err)
			assert.EqualValues(t, 1, listener.called.Load())
		})
	}
}

func TestDispatcher_DispatchRequest_ReadBodyError(t *testing.T) {
	dispatcher := NewDispatcher()

	req, err := http.NewRequest(http.MethodPost, "/webhook", nil)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}

	req.Header.Set("X-Gitlab-Event", string(gitlab.EventTypePush))
	req.Body = errReader{}

	err = dispatcher.DispatchRequest(req)

	assert.EqualError(t, err, "read error")
}

func TestDispatcher_DispatchRequest_ListenerErrorsJoined(t *testing.T) {
	errOne := errors.New("listener one failed")
	errTwo := errors.New("listener two failed")

	dispatcher := NewDispatcher(RegisterListeners(
		&errorPushListener{err: errOne},
		&errorPushListener{err: errTwo},
	))

	req := newWebhookRequest(t, gitlab.EventTypePush, loadFixture(t, "testdata/webhooks/push.json"))

	err := dispatcher.DispatchRequest(req)

	assert.Error(t, err)
	assert.ErrorIs(t, err, errOne)
	assert.ErrorIs(t, err, errTwo)
}
