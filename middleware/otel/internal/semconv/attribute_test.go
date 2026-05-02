package semconv_test

import (
	"testing"

	"github.com/flc1125/go-gitlab-webhook/middleware/otel/v3/internal/semconv"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
)

func TestAttributes(t *testing.T) {
	tests := []struct {
		name string
		got  attribute.KeyValue
		want attribute.KeyValue
	}{
		{
			name: "WebhookEventType",
			got:  semconv.WebhookEventType("push"),
			want: attribute.String("gitlab.webhook.event_type", "push"),
		},
		{
			name: "WebhookObjectKind",
			got:  semconv.WebhookObjectKind("merge_request"),
			want: attribute.String("gitlab.webhook.object_kind", "merge_request"),
		},
		{
			name: "WebhookEventName",
			got:  semconv.WebhookEventName("tag_push"),
			want: attribute.String("gitlab.webhook.event_name", "tag_push"),
		},
		{
			name: "WebhookAction",
			got:  semconv.WebhookAction("open"),
			want: attribute.String("gitlab.webhook.action", "open"),
		},
		{
			name: "WebhookStatus",
			got:  semconv.WebhookStatus("200"),
			want: attribute.String("gitlab.webhook.status", "200"),
		},
		{
			name: "WebhookResult",
			got:  semconv.WebhookResult("success"),
			want: attribute.String("gitlab.webhook.result", "success"),
		},
		{
			name: "WebhookGoType",
			got:  semconv.WebhookGoType("*gitlab.PushEvent"),
			want: attribute.String("gitlab.webhook.go_type", "*gitlab.PushEvent"),
		},
		{
			name: "ProjectID",
			got:  semconv.ProjectID(123),
			want: attribute.Int64("gitlab.project.id", 123),
		},
		{
			name: "ProjectName",
			got:  semconv.ProjectName("project"),
			want: attribute.String("gitlab.project.name", "project"),
		},
		{
			name: "ProjectPath",
			got:  semconv.ProjectPath("group/project"),
			want: attribute.String("gitlab.project.path", "group/project"),
		},
		{
			name: "ProjectNamespaceID",
			got:  semconv.ProjectNamespaceID(456),
			want: attribute.Int64("gitlab.project.namespace_id", 456),
		},
		{
			name: "ProjectVisibility",
			got:  semconv.ProjectVisibility("private"),
			want: attribute.String("gitlab.project.visibility", "private"),
		},
		{
			name: "GroupID",
			got:  semconv.GroupID(789),
			want: attribute.Int64("gitlab.group.id", 789),
		},
		{
			name: "GroupName",
			got:  semconv.GroupName("group"),
			want: attribute.String("gitlab.group.name", "group"),
		},
		{
			name: "GroupPath",
			got:  semconv.GroupPath("parent/group"),
			want: attribute.String("gitlab.group.path", "parent/group"),
		},
		{
			name: "ParentGroupID",
			got:  semconv.ParentGroupID(321),
			want: attribute.Int64("gitlab.parent_group.id", 321),
		},
		{
			name: "ParentGroupPath",
			got:  semconv.ParentGroupPath("parent"),
			want: attribute.String("gitlab.parent_group.path", "parent"),
		},
		{
			name: "Ref",
			got:  semconv.Ref("refs/heads/main"),
			want: attribute.String("gitlab.ref", "refs/heads/main"),
		},
		{
			name: "SHA",
			got:  semconv.SHA("abc123"),
			want: attribute.String("gitlab.sha", "abc123"),
		},
		{
			name: "EmojiAwardableType",
			got:  semconv.EmojiAwardableType("merge_request"),
			want: attribute.String("gitlab.emoji.awardable_type", "merge_request"),
		},
		{
			name: "FeatureFlagID",
			got:  semconv.FeatureFlagID(654),
			want: attribute.Int64("gitlab.feature_flag.id", 654),
		},
		{
			name: "FeatureFlagName",
			got:  semconv.FeatureFlagName("flag"),
			want: attribute.String("gitlab.feature_flag.name", "flag"),
		},
		{
			name: "FeatureFlagActive",
			got:  semconv.FeatureFlagActive(true),
			want: attribute.Bool("gitlab.feature_flag.active", true),
		},
		{
			name: "IssueIID",
			got:  semconv.IssueIID(101),
			want: attribute.Int64("gitlab.issue.iid", 101),
		},
		{
			name: "IssueState",
			got:  semconv.IssueState("opened"),
			want: attribute.String("gitlab.issue.state", "opened"),
		},
		{
			name: "JobID",
			got:  semconv.JobID(202),
			want: attribute.Int64("gitlab.job.id", 202),
		},
		{
			name: "JobName",
			got:  semconv.JobName("test"),
			want: attribute.String("gitlab.job.name", "test"),
		},
		{
			name: "JobStage",
			got:  semconv.JobStage("build"),
			want: attribute.String("gitlab.job.stage", "build"),
		},
		{
			name: "MergeRequestIID",
			got:  semconv.MergeRequestIID(303),
			want: attribute.Int64("gitlab.merge_request.iid", 303),
		},
		{
			name: "MergeRequestState",
			got:  semconv.MergeRequestState("merged"),
			want: attribute.String("gitlab.merge_request.state", "merged"),
		},
		{
			name: "MergeRequestMergeStatus",
			got:  semconv.MergeRequestMergeStatus("can_be_merged"),
			want: attribute.String("gitlab.merge_request.merge_status", "can_be_merged"),
		},
		{
			name: "MergeRequestDetailedMergeStatus",
			got:  semconv.MergeRequestDetailedMergeStatus("mergeable"),
			want: attribute.String("gitlab.merge_request.detailed_merge_status", "mergeable"),
		},
		{
			name: "MilestoneIID",
			got:  semconv.MilestoneIID(404),
			want: attribute.Int64("gitlab.milestone.iid", 404),
		},
		{
			name: "MilestoneState",
			got:  semconv.MilestoneState("active"),
			want: attribute.String("gitlab.milestone.state", "active"),
		},
		{
			name: "NoteNoteableType",
			got:  semconv.NoteNoteableType("Issue"),
			want: attribute.String("gitlab.note.noteable_type", "Issue"),
		},
		{
			name: "PipelineID",
			got:  semconv.PipelineID(505),
			want: attribute.Int64("gitlab.pipeline.id", 505),
		},
		{
			name: "PipelineIID",
			got:  semconv.PipelineIID(606),
			want: attribute.Int64("gitlab.pipeline.iid", 606),
		},
		{
			name: "PipelineSource",
			got:  semconv.PipelineSource("push"),
			want: attribute.String("gitlab.pipeline.source", "push"),
		},
		{
			name: "PushTotalCommitsCount",
			got:  semconv.PushTotalCommitsCount(2),
			want: attribute.Int64("gitlab.push.total_commits_count", 2),
		},
		{
			name: "ReleaseTag",
			got:  semconv.ReleaseTag("v1.0.0"),
			want: attribute.String("gitlab.release.tag", "v1.0.0"),
		},
		{
			name: "ResourceAccessTokenID",
			got:  semconv.ResourceAccessTokenID(707),
			want: attribute.Int64("gitlab.resource_access_token.id", 707),
		},
		{
			name: "VulnerabilityState",
			got:  semconv.VulnerabilityState("detected"),
			want: attribute.String("gitlab.vulnerability.state", "detected"),
		},
		{
			name: "VulnerabilityReport",
			got:  semconv.VulnerabilityReport("sast"),
			want: attribute.String("gitlab.vulnerability.report_type", "sast"),
		},
		{
			name: "VulnerabilityLevel",
			got:  semconv.VulnerabilityLevel("high"),
			want: attribute.String("gitlab.vulnerability.severity", "high"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.got)
		})
	}
}
