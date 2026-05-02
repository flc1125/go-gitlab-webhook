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

func hasAttr(attrs []attribute.KeyValue, key string) bool {
	for _, attr := range attrs {
		if string(attr.Key) == key {
			return true
		}
	}

	return false
}
