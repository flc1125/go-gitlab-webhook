package eventmeta_test

import (
	"testing"

	"github.com/flc1125/go-gitlab-webhook/middleware/otel/v3/internal/eventmeta"
	"github.com/stretchr/testify/assert"
	gitlab "gitlab.com/gitlab-org/api/client-go/v2"
	"go.opentelemetry.io/otel/attribute"
)

func TestExtractPushEvent(t *testing.T) {
	metadata := eventmeta.Extract(&gitlab.PushEvent{
		ObjectKind:        "push",
		EventName:         "push",
		ProjectID:         123,
		Ref:               "refs/heads/main",
		CheckoutSHA:       "abc123",
		TotalCommitsCount: 0,
		Project: gitlab.PushEventProject{
			Name:              "project",
			PathWithNamespace: "group/project",
		},
	})

	assert.Equal(t, "gitlab.webhook.push", metadata.SpanName)
	assert.Equal(t, "push", attrString(metadata.Attributes, "gitlab.webhook.event_type"))
	assert.Equal(t, "push", attrString(metadata.Attributes, "gitlab.webhook.event_name"))
	assert.Equal(t, int64(123), attrInt64(metadata.Attributes, "gitlab.project.id"))
	assert.Equal(t, "group/project", attrString(metadata.Attributes, "gitlab.project.path"))
	assert.Equal(t, int64(0), attrInt64(metadata.Attributes, "gitlab.push.total_commits_count"))
	assert.True(t, hasAttr(metadata.Attributes, "gitlab.push.total_commits_count"))
}

func TestExtractMergeEvent(t *testing.T) {
	metadata := eventmeta.Extract(&gitlab.MergeEvent{
		ObjectKind: "merge_request",
		EventType:  "merge_request",
		Project: gitlab.MergeEventProject{
			ID:                123,
			Name:              "project",
			PathWithNamespace: "group/project",
		},
		ObjectAttributes: gitlab.MergeEventObjectAttributes{
			IID:                 42,
			Action:              "open",
			State:               "opened",
			MergeStatus:         "unchecked",
			DetailedMergeStatus: "checking",
		},
	})

	assert.Equal(t, "gitlab.webhook.merge_request open", metadata.SpanName)
	assert.Equal(t, "merge_request", attrString(metadata.Attributes, "gitlab.webhook.event_type"))
	assert.Equal(t, "open", attrString(metadata.Attributes, "gitlab.webhook.action"))
	assert.Equal(t, int64(42), attrInt64(metadata.Attributes, "gitlab.merge_request.iid"))
	assert.Equal(t, "unchecked", attrString(metadata.Attributes, "gitlab.merge_request.merge_status"))
	assert.Equal(t, "checking", attrString(metadata.Attributes, "gitlab.merge_request.detailed_merge_status"))
}

func TestExtractSupportedEvents(t *testing.T) {
	tests := []struct {
		name    string
		event   any
		span    string
		strings map[string]string
		int64s  map[string]int64
		bools   map[string]bool
	}{
		{
			name: "BuildEvent",
			event: &gitlab.BuildEvent{
				ObjectKind:  "build",
				ProjectID:   123,
				ProjectName: "project",
				Ref:         "refs/heads/main",
				SHA:         "abc123",
				BuildID:     456,
				BuildName:   "test",
				BuildStage:  "build",
				BuildStatus: "success",
			},
			span: "gitlab.webhook.build success",
			strings: map[string]string{
				"gitlab.webhook.object_kind": "build",
				"gitlab.webhook.status":      "success",
				"gitlab.project.name":        "project",
				"gitlab.ref":                 "refs/heads/main",
				"gitlab.sha":                 "abc123",
				"gitlab.job.name":            "test",
				"gitlab.job.stage":           "build",
			},
			int64s: map[string]int64{
				"gitlab.project.id": 123,
				"gitlab.job.id":     456,
			},
		},
		{
			name: "CommitCommentEvent",
			event: &gitlab.CommitCommentEvent{
				ObjectKind: "note",
				ProjectID:  123,
				Project: gitlab.CommitCommentEventProject{
					Name:              "project",
					PathWithNamespace: "group/project",
				},
				ObjectAttributes: gitlab.CommitCommentEventObjectAttributes{
					NoteableType: "Commit",
					Action:       "create",
				},
			},
			span: "gitlab.webhook.note commit create",
			strings: map[string]string{
				"gitlab.webhook.object_kind": "note",
				"gitlab.note.noteable_type":  "Commit",
				"gitlab.webhook.action":      "create",
				"gitlab.project.path":        "group/project",
				"gitlab.project.name":        "project",
			},
			int64s: map[string]int64{"gitlab.project.id": 123},
		},
		{
			name: "DeploymentEvent",
			event: &gitlab.DeploymentEvent{
				ObjectKind: "deployment",
				Status:     "running",
				Ref:        "refs/heads/main",
				Project: gitlab.DeploymentEventProject{
					ID:                123,
					Name:              "project",
					PathWithNamespace: "group/project",
				},
			},
			span: "gitlab.webhook.deployment running",
			strings: map[string]string{
				"gitlab.webhook.object_kind": "deployment",
				"gitlab.webhook.status":      "running",
				"gitlab.project.path":        "group/project",
				"gitlab.ref":                 "refs/heads/main",
			},
			int64s: map[string]int64{"gitlab.project.id": 123},
		},
		{
			name: "EmojiEvent",
			event: &gitlab.EmojiEvent{
				ObjectKind: "emoji",
				ProjectID:  123,
				Project: gitlab.EmojiEventProject{
					Name:              "project",
					PathWithNamespace: "group/project",
				},
				ObjectAttributes: gitlab.EmojiEventObjectAttributes{
					Action:        "award",
					AwardableType: "MergeRequest",
				},
			},
			span: "gitlab.webhook.emoji award",
			strings: map[string]string{
				"gitlab.webhook.object_kind":  "emoji",
				"gitlab.webhook.action":       "award",
				"gitlab.emoji.awardable_type": "MergeRequest",
				"gitlab.project.path":         "group/project",
				"gitlab.project.name":         "project",
				"gitlab.webhook.event_type":   "emoji",
			},
			int64s: map[string]int64{"gitlab.project.id": 123},
		},
		{
			name: "FeatureFlagEvent",
			event: &gitlab.FeatureFlagEvent{
				ObjectKind: "feature_flag",
				Project: gitlab.FeatureFlagEventProject{
					ID:                123,
					Name:              "project",
					PathWithNamespace: "group/project",
				},
				ObjectAttributes: gitlab.FeatureFlagEventObjectAttributes{
					ID:     456,
					Name:   "flag",
					Active: true,
				},
			},
			span: "gitlab.webhook.feature_flag",
			strings: map[string]string{
				"gitlab.webhook.object_kind": "feature_flag",
				"gitlab.feature_flag.name":   "flag",
				"gitlab.project.path":        "group/project",
			},
			int64s: map[string]int64{
				"gitlab.project.id":      123,
				"gitlab.feature_flag.id": 456,
			},
			bools: map[string]bool{"gitlab.feature_flag.active": true},
		},
		{
			name: "GroupResourceAccessTokenEvent",
			event: &gitlab.GroupResourceAccessTokenEvent{
				ObjectKind: "resource_access_token",
				EventName:  "resourceAccessTokenCreated",
				Group: gitlab.GroupResourceAccessTokenEventGroup{
					GroupID:   123,
					GroupName: "group",
					FullPath:  "parent/group",
				},
				ObjectAttributes: gitlab.GroupResourceAccessTokenEventObjectAttributes{ID: 456},
			},
			span: "gitlab.webhook.resource_access_token resource_access_token_created",
			strings: map[string]string{
				"gitlab.webhook.event_name":  "resourceAccessTokenCreated",
				"gitlab.webhook.object_kind": "resource_access_token",
				"gitlab.group.name":          "group",
				"gitlab.group.path":          "parent/group",
			},
			int64s: map[string]int64{
				"gitlab.group.id":                 123,
				"gitlab.resource_access_token.id": 456,
			},
		},
		{
			name: "IssueCommentEvent",
			event: &gitlab.IssueCommentEvent{
				ObjectKind: "note",
				ProjectID:  123,
				Project: gitlab.IssueCommentEventProject{
					Name:              "project",
					PathWithNamespace: "group/project",
				},
				ObjectAttributes: gitlab.IssueCommentEventObjectAttributes{
					NoteableType: "Issue",
					Action:       "update",
				},
				Issue: gitlab.IssueCommentEventIssue{
					IID:   456,
					State: "opened",
				},
			},
			span: "gitlab.webhook.note issue update",
			strings: map[string]string{
				"gitlab.note.noteable_type": "Issue",
				"gitlab.webhook.action":     "update",
				"gitlab.issue.state":        "opened",
			},
			int64s: map[string]int64{
				"gitlab.project.id": 123,
				"gitlab.issue.iid":  456,
			},
		},
		{
			name: "IssueEvent",
			event: &gitlab.IssueEvent{
				ObjectKind: "issue",
				Project: gitlab.IssueEventProject{
					ID:                123,
					Name:              "project",
					PathWithNamespace: "group/project",
				},
				ObjectAttributes: gitlab.IssueEventObjectAttributes{
					IID:    456,
					State:  "opened",
					Action: "reopen",
				},
			},
			span: "gitlab.webhook.issue reopen",
			strings: map[string]string{
				"gitlab.webhook.action": "reopen",
				"gitlab.issue.state":    "opened",
				"gitlab.project.path":   "group/project",
			},
			int64s: map[string]int64{
				"gitlab.project.id": 123,
				"gitlab.issue.iid":  456,
			},
		},
		{
			name: "JobEvent",
			event: &gitlab.JobEvent{
				ObjectKind:  "build",
				ProjectID:   123,
				ProjectName: "project",
				Ref:         "refs/heads/main",
				SHA:         "abc123",
				BuildID:     456,
				BuildName:   "test",
				BuildStage:  "build",
				BuildStatus: "failed",
			},
			span: "gitlab.webhook.job failed",
			strings: map[string]string{
				"gitlab.webhook.object_kind": "build",
				"gitlab.webhook.status":      "failed",
				"gitlab.project.name":        "project",
				"gitlab.job.name":            "test",
			},
			int64s: map[string]int64{
				"gitlab.project.id": 123,
				"gitlab.job.id":     456,
			},
		},
		{
			name: "MemberEvent",
			event: &gitlab.MemberEvent{
				EventName: "user_add_to_group",
				GroupID:   123,
				GroupPath: "parent/group",
				GroupName: "group",
			},
			span: "gitlab.webhook.member user_add_to_group",
			strings: map[string]string{
				"gitlab.webhook.event_name": "user_add_to_group",
				"gitlab.group.path":         "parent/group",
				"gitlab.group.name":         "group",
			},
			int64s: map[string]int64{"gitlab.group.id": 123},
		},
		{
			name: "MilestoneWebhookEvent with group",
			event: &gitlab.MilestoneWebhookEvent{
				ObjectKind: "milestone",
				Action:     "close",
				Project: gitlab.MilestoneEventProject{
					ID:                123,
					Name:              "project",
					PathWithNamespace: "group/project",
				},
				Group: &gitlab.MilestoneEventGroup{
					GroupID:   456,
					GroupName: "group",
					FullPath:  "parent/group",
				},
				ObjectAttributes: gitlab.MilestoneEventObjectAttributes{
					IID:   789,
					State: "closed",
				},
			},
			span: "gitlab.webhook.milestone close",
			strings: map[string]string{
				"gitlab.webhook.action":  "close",
				"gitlab.milestone.state": "closed",
				"gitlab.group.path":      "parent/group",
			},
			int64s: map[string]int64{
				"gitlab.project.id":    123,
				"gitlab.group.id":      456,
				"gitlab.milestone.iid": 789,
			},
		},
		{
			name: "MergeCommentEvent",
			event: &gitlab.MergeCommentEvent{
				ObjectKind: "note",
				ProjectID:  123,
				Project: gitlab.MergeCommentEventProject{
					Name:              "project",
					PathWithNamespace: "group/project",
				},
				ObjectAttributes: gitlab.MergeCommentEventObjectAttributes{
					NoteableType: "MergeRequest",
					Action:       "create",
				},
				MergeRequest: gitlab.MergeCommentEventMergeRequest{
					IID:   456,
					State: "opened",
				},
			},
			span: "gitlab.webhook.note merge_request create",
			strings: map[string]string{
				"gitlab.note.noteable_type":  "MergeRequest",
				"gitlab.webhook.action":      "create",
				"gitlab.merge_request.state": "opened",
				"gitlab.project.path":        "group/project",
			},
			int64s: map[string]int64{
				"gitlab.project.id":        123,
				"gitlab.merge_request.iid": 456,
			},
		},
		{
			name: "PipelineEvent",
			event: &gitlab.PipelineEvent{
				ObjectKind: "pipeline",
				Project: gitlab.PipelineEventProject{
					ID:                123,
					Name:              "project",
					PathWithNamespace: "group/project",
				},
				ObjectAttributes: gitlab.PipelineEventObjectAttributes{
					ID:     456,
					IID:    789,
					Ref:    "refs/heads/main",
					SHA:    "abc123",
					Source: "push",
					Status: "success",
				},
			},
			span: "gitlab.webhook.pipeline success",
			strings: map[string]string{
				"gitlab.webhook.status":  "success",
				"gitlab.pipeline.source": "push",
				"gitlab.ref":             "refs/heads/main",
				"gitlab.sha":             "abc123",
			},
			int64s: map[string]int64{
				"gitlab.project.id":   123,
				"gitlab.pipeline.id":  456,
				"gitlab.pipeline.iid": 789,
			},
		},
		{
			name: "ProjectWebhookEvent",
			event: &gitlab.ProjectWebhookEvent{
				EventName:          "projectCreated Now/foo-bar",
				ProjectID:          123,
				ProjectNamespaceID: 456,
				ProjectVisibility:  "private",
				Name:               "project",
				PathWithNamespace:  "group/project",
			},
			span: "gitlab.webhook.project project_created_now_foo_bar",
			strings: map[string]string{
				"gitlab.webhook.event_name": "projectCreated Now/foo-bar",
				"gitlab.project.visibility": "private",
				"gitlab.project.name":       "project",
				"gitlab.project.path":       "group/project",
			},
			int64s: map[string]int64{
				"gitlab.project.id":           123,
				"gitlab.project.namespace_id": 456,
			},
		},
		{
			name: "ProjectResourceAccessTokenEvent",
			event: &gitlab.ProjectResourceAccessTokenEvent{
				ObjectKind: "resource_access_token",
				EventName:  "resource_access_token_destroyed",
				Project: gitlab.ProjectResourceAccessTokenEventProject{
					ID:                123,
					Name:              "project",
					PathWithNamespace: "group/project",
				},
				ObjectAttributes: gitlab.ProjectResourceAccessTokenEventObjectAttributes{ID: 456},
			},
			span: "gitlab.webhook.resource_access_token resource_access_token_destroyed",
			strings: map[string]string{
				"gitlab.webhook.event_name": "resource_access_token_destroyed",
				"gitlab.project.path":       "group/project",
			},
			int64s: map[string]int64{
				"gitlab.project.id":               123,
				"gitlab.resource_access_token.id": 456,
			},
		},
		{
			name: "ReleaseEvent",
			event: &gitlab.ReleaseEvent{
				ObjectKind: "release",
				Action:     "create",
				Tag:        "v1.0.0",
				Project: gitlab.ReleaseEventProject{
					ID:                123,
					Name:              "project",
					PathWithNamespace: "group/project",
				},
			},
			span: "gitlab.webhook.release create",
			strings: map[string]string{
				"gitlab.webhook.action": "create",
				"gitlab.release.tag":    "v1.0.0",
				"gitlab.project.path":   "group/project",
			},
			int64s: map[string]int64{"gitlab.project.id": 123},
		},
		{
			name: "SnippetCommentEvent",
			event: &gitlab.SnippetCommentEvent{
				ObjectKind: "note",
				ProjectID:  123,
				Project: gitlab.SnippetCommentEventProject{
					Name:              "project",
					PathWithNamespace: "group/project",
				},
				ObjectAttributes: gitlab.SnippetCommentEventObjectAttributes{
					NoteableType: "Snippet",
					Action:       "create",
				},
			},
			span: "gitlab.webhook.note snippet create",
			strings: map[string]string{
				"gitlab.note.noteable_type": "Snippet",
				"gitlab.webhook.action":     "create",
				"gitlab.project.path":       "group/project",
			},
			int64s: map[string]int64{"gitlab.project.id": 123},
		},
		{
			name: "SubGroupEvent",
			event: &gitlab.SubGroupEvent{
				EventName:      "group_rename",
				Name:           "subgroup",
				FullPath:       "parent/subgroup",
				GroupID:        123,
				ParentGroupID:  456,
				ParentFullPath: "parent",
			},
			span: "gitlab.webhook.subgroup group_rename",
			strings: map[string]string{
				"gitlab.webhook.event_name": "group_rename",
				"gitlab.group.name":         "subgroup",
				"gitlab.group.path":         "parent/subgroup",
				"gitlab.parent_group.path":  "parent",
			},
			int64s: map[string]int64{
				"gitlab.group.id":        123,
				"gitlab.parent_group.id": 456,
			},
		},
		{
			name: "TagEvent",
			event: &gitlab.TagEvent{
				ObjectKind:        "tag_push",
				EventName:         "tag_push",
				ProjectID:         123,
				Ref:               "refs/tags/v1.0.0",
				CheckoutSHA:       "abc123",
				TotalCommitsCount: 2,
				Project: gitlab.TagEventProject{
					Name:              "project",
					PathWithNamespace: "group/project",
				},
			},
			span: "gitlab.webhook.tag_push",
			strings: map[string]string{
				"gitlab.webhook.event_name": "tag_push",
				"gitlab.ref":                "refs/tags/v1.0.0",
				"gitlab.sha":                "abc123",
				"gitlab.project.path":       "group/project",
			},
			int64s: map[string]int64{
				"gitlab.project.id":               123,
				"gitlab.push.total_commits_count": 2,
			},
		},
		{
			name: "VulnerabilityEvent",
			event: &gitlab.VulnerabilityEvent{
				ObjectKind: "vulnerability",
				ObjectAttributes: gitlab.VulnerabilityEventObjectAttributes{
					ProjectID:  123,
					State:      "detected",
					Severity:   "high",
					ReportType: "sast",
				},
			},
			span: "gitlab.webhook.vulnerability detected",
			strings: map[string]string{
				"gitlab.webhook.status":            "detected",
				"gitlab.vulnerability.state":       "detected",
				"gitlab.vulnerability.severity":    "high",
				"gitlab.vulnerability.report_type": "sast",
			},
			int64s: map[string]int64{"gitlab.project.id": 123},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metadata := eventmeta.Extract(tt.event)

			assert.Equal(t, tt.span, metadata.SpanName)
			for key, value := range tt.strings {
				assert.Equal(t, value, attrString(metadata.Attributes, key), key)
			}
			for key, value := range tt.int64s {
				assert.Equal(t, value, attrInt64(metadata.Attributes, key), key)
			}
			for key, value := range tt.bools {
				assert.Equal(t, value, attrBool(metadata.Attributes, key), key)
			}
		})
	}
}

func TestExtractWikiPageEventDoesNotEmitUnknownProjectID(t *testing.T) {
	metadata := eventmeta.Extract(&gitlab.WikiPageEvent{
		ObjectKind: "wiki_page",
		Project: gitlab.WikiPageEventProject{
			Name:              "project",
			PathWithNamespace: "group/project",
		},
		ObjectAttributes: gitlab.WikiPageEventObjectAttributes{
			Action: "create",
		},
	})

	assert.Equal(t, "gitlab.webhook.wiki_page create", metadata.SpanName)
	assert.False(t, hasAttr(metadata.Attributes, "gitlab.project.id"))
	assert.Equal(t, "group/project", attrString(metadata.Attributes, "gitlab.project.path"))
	assert.Equal(t, "project", attrString(metadata.Attributes, "gitlab.project.name"))
}

func TestExtractUnknownEvent(t *testing.T) {
	metadata := eventmeta.Extract(struct{}{})

	assert.Equal(t, "gitlab.webhook.unknown", metadata.SpanName)
	assert.Equal(t, "unknown", attrString(metadata.Attributes, "gitlab.webhook.event_type"))
	assert.NotEmpty(t, attrString(metadata.Attributes, "gitlab.webhook.go_type"))
}

func TestExtractUnknownPointerEvent(t *testing.T) {
	metadata := eventmeta.Extract(&unknownEvent{})

	assert.Equal(t, "gitlab.webhook.unknown", metadata.SpanName)
	assert.Equal(t, "unknown", attrString(metadata.Attributes, "gitlab.webhook.event_type"))
	assert.Equal(t, "github.com/flc1125/go-gitlab-webhook/middleware/otel/v3/internal/eventmeta_test.unknownEvent", attrString(metadata.Attributes, "gitlab.webhook.go_type"))
}

func TestExtractNilEvent(t *testing.T) {
	metadata := eventmeta.Extract(nil)

	assert.Equal(t, "gitlab.webhook.unknown", metadata.SpanName)
	assert.Equal(t, "unknown", attrString(metadata.Attributes, "gitlab.webhook.event_type"))
	assert.False(t, hasAttr(metadata.Attributes, "gitlab.webhook.go_type"))
}

type unknownEvent struct{}

func attrString(attrs []attribute.KeyValue, key string) string {
	for _, attr := range attrs {
		if string(attr.Key) == key {
			return attr.Value.AsString()
		}
	}

	return ""
}

func attrInt64(attrs []attribute.KeyValue, key string) int64 {
	for _, attr := range attrs {
		if string(attr.Key) == key {
			return attr.Value.AsInt64()
		}
	}

	return 0
}

func attrBool(attrs []attribute.KeyValue, key string) bool {
	for _, attr := range attrs {
		if string(attr.Key) == key {
			return attr.Value.AsBool()
		}
	}

	return false
}

func hasAttr(attrs []attribute.KeyValue, key string) bool {
	for _, attr := range attrs {
		if string(attr.Key) == key {
			return true
		}
	}

	return false
}
