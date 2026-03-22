package gitlabwebhook

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	gitlab "gitlab.com/gitlab-org/api/client-go/v2"
)

func TestDispatcher_DispatchWebhook_Fixtures(t *testing.T) {
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

			err := dispatcher.DispatchWebhook(
				newDispatcherContext(context.Background()),
				tt.eventType,
				loadFixture(t, tt.filepath),
			)

			assert.NoError(t, err)
			assert.EqualValues(t, 2, listener.calls.Load())
		})
	}
}

func TestDispatcher_Dispatch_UnsupportedEvent(t *testing.T) {
	dispatcher := NewDispatcher()

	err := dispatcher.Dispatch(context.Background(), struct{}{})

	assert.ErrorIs(t, err, ErrUnsupportedEvent)
}

func TestDispatcher_DispatchWebhook_InvalidPayload(t *testing.T) {
	dispatcher := NewDispatcher()

	err := dispatcher.DispatchWebhook(context.Background(), gitlab.EventTypePush, []byte(`{"invalid"`))

	assert.Error(t, err)
	assert.NotErrorIs(t, err, ErrUnsupportedEvent)
}

func TestDispatcher_DispatchWebhook_NoListeners(t *testing.T) {
	dispatcher := NewDispatcher()

	err := dispatcher.DispatchWebhook(
		context.Background(),
		gitlab.EventTypePush,
		loadFixture(t, "testdata/webhooks/push.json"),
	)

	assert.NoError(t, err)
}

func TestDispatcher_DispatchWebhook_ListenerErrorsJoined(t *testing.T) {
	errOne := errors.New("listener one failed")
	errTwo := errors.New("listener two failed")

	dispatcher := NewDispatcher(RegisterListeners(
		&errorPushListener{err: errOne},
		&errorPushListener{err: errTwo},
	))

	err := dispatcher.DispatchWebhook(
		context.Background(),
		gitlab.EventTypePush,
		loadFixture(t, "testdata/webhooks/push.json"),
	)

	assert.Error(t, err)
	assert.ErrorIs(t, err, errOne)
	assert.ErrorIs(t, err, errTwo)
	assert.Contains(t, err.Error(), errOne.Error())
	assert.Contains(t, err.Error(), errTwo.Error())
}

func TestDispatcher_Dispatch_ParsedEvent_NoListeners(t *testing.T) {
	dispatcher := NewDispatcher()

	event, err := gitlab.ParseWebhook(gitlab.EventTypePush, loadFixture(t, "testdata/webhooks/push.json"))
	if err != nil {
		t.Fatalf("parse push fixture: %v", err)
	}

	err = dispatcher.Dispatch(context.Background(), event)

	assert.NoError(t, err)
}

func TestExpectString(t *testing.T) {
	assert.NoError(t, expectString("field", "same", "same"))

	err := expectString("field", "left", "right")
	assert.EqualError(t, err, fmt.Sprintf("field: got %q, want %q", "left", "right"))
}
