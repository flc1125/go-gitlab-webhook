package gitlabwebhook

import (
	"context"

	gitlab "gitlab.com/gitlab-org/api/client-go/v2"
)

// BuildListener handles GitLab build webhook events.
type BuildListener interface {
	OnBuild(ctx context.Context, event *gitlab.BuildEvent) error
}

// CommitCommentListener handles GitLab commit comment webhook events.
type CommitCommentListener interface {
	OnCommitComment(ctx context.Context, event *gitlab.CommitCommentEvent) error
}

// DeploymentListener handles GitLab deployment webhook events.
type DeploymentListener interface {
	OnDeployment(ctx context.Context, event *gitlab.DeploymentEvent) error
}

// EmojiListener handles GitLab emoji webhook events.
type EmojiListener interface {
	OnEmoji(ctx context.Context, event *gitlab.EmojiEvent) error
}

// FeatureFlagListener handles GitLab feature flag webhook events.
type FeatureFlagListener interface {
	OnFeatureFlag(ctx context.Context, event *gitlab.FeatureFlagEvent) error
}

// GroupResourceAccessTokenListener handles GitLab group resource access token webhook events.
type GroupResourceAccessTokenListener interface {
	OnGroupResourceAccessToken(ctx context.Context, event *gitlab.GroupResourceAccessTokenEvent) error
}

// IssueCommentListener handles GitLab issue comment webhook events.
type IssueCommentListener interface {
	OnIssueComment(ctx context.Context, event *gitlab.IssueCommentEvent) error
}

// IssueListener handles GitLab issue webhook events.
type IssueListener interface {
	OnIssue(ctx context.Context, event *gitlab.IssueEvent) error
}

// JobListener handles GitLab job webhook events.
type JobListener interface {
	OnJob(ctx context.Context, event *gitlab.JobEvent) error
}

// MemberListener handles GitLab member webhook events.
type MemberListener interface {
	OnMember(ctx context.Context, event *gitlab.MemberEvent) error
}

// MilestoneListener handles GitLab milestone webhook events.
type MilestoneListener interface {
	OnMilestone(ctx context.Context, event *gitlab.MilestoneWebhookEvent) error
}

// MergeCommentListener handles GitLab merge request comment webhook events.
type MergeCommentListener interface {
	OnMergeComment(ctx context.Context, event *gitlab.MergeCommentEvent) error
}

// MergeListener handles GitLab merge request webhook events.
type MergeListener interface {
	OnMerge(ctx context.Context, event *gitlab.MergeEvent) error
}

// PipelineListener handles GitLab pipeline webhook events.
type PipelineListener interface {
	OnPipeline(ctx context.Context, event *gitlab.PipelineEvent) error
}

// ProjectListener handles GitLab project webhook events.
type ProjectListener interface {
	OnProject(ctx context.Context, event *gitlab.ProjectWebhookEvent) error
}

// ProjectResourceAccessTokenListener handles GitLab project resource access token webhook events.
type ProjectResourceAccessTokenListener interface {
	OnProjectResourceAccessToken(ctx context.Context, event *gitlab.ProjectResourceAccessTokenEvent) error
}

// PushListener handles GitLab push webhook events.
type PushListener interface {
	OnPush(ctx context.Context, event *gitlab.PushEvent) error
}

// ReleaseListener handles GitLab release webhook events.
type ReleaseListener interface {
	OnRelease(ctx context.Context, event *gitlab.ReleaseEvent) error
}

// SnippetCommentListener handles GitLab snippet comment webhook events.
type SnippetCommentListener interface {
	OnSnippetComment(ctx context.Context, event *gitlab.SnippetCommentEvent) error
}

// SubGroupListener handles GitLab subgroup webhook events.
type SubGroupListener interface {
	OnSubGroup(ctx context.Context, event *gitlab.SubGroupEvent) error
}

// TagListener handles GitLab tag push webhook events.
type TagListener interface {
	OnTag(ctx context.Context, event *gitlab.TagEvent) error
}

// VulnerabilityListener handles GitLab vulnerability webhook events.
type VulnerabilityListener interface {
	OnVulnerability(ctx context.Context, event *gitlab.VulnerabilityEvent) error
}

// WikiPageListener handles GitLab wiki page webhook events.
type WikiPageListener interface {
	OnWikiPage(ctx context.Context, event *gitlab.WikiPageEvent) error
}
