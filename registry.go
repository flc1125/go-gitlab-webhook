package gitlabwebhook

import (
	"context"
)

type listenerRegistration func(*Dispatcher, any)
type eventDispatcher func(*Dispatcher, context.Context, any) (bool, error)

var dispatcherListenerRegistrations = []listenerRegistration{
	registerListenerType[BuildListener](func(d *Dispatcher, listener BuildListener) {
		d.RegisterBuildListener(listener)
	}),
	registerListenerType[CommitCommentListener](func(d *Dispatcher, listener CommitCommentListener) {
		d.RegisterCommitCommentListener(listener)
	}),
	registerListenerType[DeploymentListener](func(d *Dispatcher, listener DeploymentListener) {
		d.RegisterDeploymentListener(listener)
	}),
	registerListenerType[EmojiListener](func(d *Dispatcher, listener EmojiListener) {
		d.RegisterEmojiListener(listener)
	}),
	registerListenerType[FeatureFlagListener](func(d *Dispatcher, listener FeatureFlagListener) {
		d.RegisterFeatureFlagListener(listener)
	}),
	registerListenerType[GroupResourceAccessTokenListener](func(d *Dispatcher, listener GroupResourceAccessTokenListener) {
		d.RegisterGroupResourceAccessTokenListener(listener)
	}),
	registerListenerType[IssueCommentListener](func(d *Dispatcher, listener IssueCommentListener) {
		d.RegisterIssueCommentListener(listener)
	}),
	registerListenerType[IssueListener](func(d *Dispatcher, listener IssueListener) {
		d.RegisterIssueListener(listener)
	}),
	registerListenerType[JobListener](func(d *Dispatcher, listener JobListener) {
		d.RegisterJobListener(listener)
	}),
	registerListenerType[MemberListener](func(d *Dispatcher, listener MemberListener) {
		d.RegisterMemberListener(listener)
	}),
	registerListenerType[MilestoneListener](func(d *Dispatcher, listener MilestoneListener) {
		d.RegisterMilestoneListener(listener)
	}),
	registerListenerType[MergeCommentListener](func(d *Dispatcher, listener MergeCommentListener) {
		d.RegisterMergeCommentListener(listener)
	}),
	registerListenerType[MergeListener](func(d *Dispatcher, listener MergeListener) {
		d.RegisterMergeListener(listener)
	}),
	registerListenerType[PipelineListener](func(d *Dispatcher, listener PipelineListener) {
		d.RegisterPipelineListener(listener)
	}),
	registerListenerType[ProjectListener](func(d *Dispatcher, listener ProjectListener) {
		d.RegisterProjectListener(listener)
	}),
	registerListenerType[ProjectResourceAccessTokenListener](func(d *Dispatcher, listener ProjectResourceAccessTokenListener) {
		d.RegisterProjectResourceAccessTokenListener(listener)
	}),
	registerListenerType[PushListener](func(d *Dispatcher, listener PushListener) {
		d.RegisterPushListener(listener)
	}),
	registerListenerType[ReleaseListener](func(d *Dispatcher, listener ReleaseListener) {
		d.RegisterReleaseListener(listener)
	}),
	registerListenerType[SnippetCommentListener](func(d *Dispatcher, listener SnippetCommentListener) {
		d.RegisterSnippetCommentListener(listener)
	}),
	registerListenerType[SubGroupListener](func(d *Dispatcher, listener SubGroupListener) {
		d.RegisterSubGroupListener(listener)
	}),
	registerListenerType[TagListener](func(d *Dispatcher, listener TagListener) {
		d.RegisterTagListener(listener)
	}),
	registerListenerType[VulnerabilityListener](func(d *Dispatcher, listener VulnerabilityListener) {
		d.RegisterVulnerabilityListener(listener)
	}),
	registerListenerType[WikiPageListener](func(d *Dispatcher, listener WikiPageListener) {
		d.RegisterWikiPageListener(listener)
	}),
}

var dispatcherEventDispatchers = []eventDispatcher{
	dispatchEventType(func(d *Dispatcher) []BuildListener { return d.buildListeners }, BuildListener.OnBuild),
	dispatchEventType(func(d *Dispatcher) []CommitCommentListener { return d.commitCommentListeners }, CommitCommentListener.OnCommitComment),
	dispatchEventType(func(d *Dispatcher) []DeploymentListener { return d.deploymentListeners }, DeploymentListener.OnDeployment),
	dispatchEventType(func(d *Dispatcher) []EmojiListener { return d.emojiListeners }, EmojiListener.OnEmoji),
	dispatchEventType(func(d *Dispatcher) []FeatureFlagListener { return d.featureFlagListeners }, FeatureFlagListener.OnFeatureFlag),
	dispatchEventType(func(d *Dispatcher) []GroupResourceAccessTokenListener { return d.groupResourceAccessTokenListeners }, GroupResourceAccessTokenListener.OnGroupResourceAccessToken),
	dispatchEventType(func(d *Dispatcher) []IssueCommentListener { return d.issueCommentListeners }, IssueCommentListener.OnIssueComment),
	dispatchEventType(func(d *Dispatcher) []IssueListener { return d.issueListeners }, IssueListener.OnIssue),
	dispatchEventType(func(d *Dispatcher) []JobListener { return d.jobListeners }, JobListener.OnJob),
	dispatchEventType(func(d *Dispatcher) []MemberListener { return d.memberListeners }, MemberListener.OnMember),
	dispatchEventType(func(d *Dispatcher) []MilestoneListener { return d.milestoneListeners }, MilestoneListener.OnMilestone),
	dispatchEventType(func(d *Dispatcher) []MergeCommentListener { return d.mergeCommentListeners }, MergeCommentListener.OnMergeComment),
	dispatchEventType(func(d *Dispatcher) []MergeListener { return d.mergeListeners }, MergeListener.OnMerge),
	dispatchEventType(func(d *Dispatcher) []PipelineListener { return d.pipelineListeners }, PipelineListener.OnPipeline),
	dispatchEventType(func(d *Dispatcher) []ProjectListener { return d.projectListeners }, ProjectListener.OnProject),
	dispatchEventType(func(d *Dispatcher) []ProjectResourceAccessTokenListener { return d.projectResourceAccessTokenListeners }, ProjectResourceAccessTokenListener.OnProjectResourceAccessToken),
	dispatchEventType(func(d *Dispatcher) []PushListener { return d.pushListeners }, PushListener.OnPush),
	dispatchEventType(func(d *Dispatcher) []ReleaseListener { return d.releaseListeners }, ReleaseListener.OnRelease),
	dispatchEventType(func(d *Dispatcher) []SnippetCommentListener { return d.snippetCommentListeners }, SnippetCommentListener.OnSnippetComment),
	dispatchEventType(func(d *Dispatcher) []SubGroupListener { return d.subGroupListeners }, SubGroupListener.OnSubGroup),
	dispatchEventType(func(d *Dispatcher) []TagListener { return d.tagListeners }, TagListener.OnTag),
	dispatchEventType(func(d *Dispatcher) []VulnerabilityListener { return d.vulnerabilityListeners }, VulnerabilityListener.OnVulnerability),
	dispatchEventType(func(d *Dispatcher) []WikiPageListener { return d.wikiPageListeners }, WikiPageListener.OnWikiPage),
}

func registerListenerType[L any](register func(*Dispatcher, L)) listenerRegistration {
	return func(d *Dispatcher, listener any) {
		typedListener, ok := listener.(L)
		if !ok {
			return
		}

		register(d, typedListener)
	}
}

func dispatchEventType[E any, L any](
	listeners func(*Dispatcher) []L,
	handler func(L, context.Context, E) error,
) eventDispatcher {
	return func(d *Dispatcher, ctx context.Context, event any) (bool, error) {
		typedEvent, ok := event.(E)
		if !ok {
			return false, nil
		}

		return true, processEvent(ctx, listeners(d), handler, typedEvent)
	}
}
