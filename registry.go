package gitlabwebhook

import (
	"context"
)

type eventDescriptor interface {
	register(*Dispatcher, any)
	dispatch(*Dispatcher, context.Context, any) (bool, error)
}

type typedEventDescriptor[E any, L any] struct {
	listeners func(*Dispatcher) []L
	set       func(*Dispatcher, []L)
	handler   func(L, context.Context, E) error
}

var dispatcherEvents = []eventDescriptor{
	newEventDescriptor(func(d *Dispatcher) []BuildListener { return d.buildListeners }, func(d *Dispatcher, v []BuildListener) { d.buildListeners = v }, BuildListener.OnBuild),
	newEventDescriptor(func(d *Dispatcher) []CommitCommentListener { return d.commitCommentListeners }, func(d *Dispatcher, v []CommitCommentListener) { d.commitCommentListeners = v }, CommitCommentListener.OnCommitComment),
	newEventDescriptor(func(d *Dispatcher) []DeploymentListener { return d.deploymentListeners }, func(d *Dispatcher, v []DeploymentListener) { d.deploymentListeners = v }, DeploymentListener.OnDeployment),
	newEventDescriptor(func(d *Dispatcher) []EmojiListener { return d.emojiListeners }, func(d *Dispatcher, v []EmojiListener) { d.emojiListeners = v }, EmojiListener.OnEmoji),
	newEventDescriptor(func(d *Dispatcher) []FeatureFlagListener { return d.featureFlagListeners }, func(d *Dispatcher, v []FeatureFlagListener) { d.featureFlagListeners = v }, FeatureFlagListener.OnFeatureFlag),
	newEventDescriptor(func(d *Dispatcher) []GroupResourceAccessTokenListener { return d.groupResourceAccessTokenListeners }, func(d *Dispatcher, v []GroupResourceAccessTokenListener) { d.groupResourceAccessTokenListeners = v }, GroupResourceAccessTokenListener.OnGroupResourceAccessToken),
	newEventDescriptor(func(d *Dispatcher) []IssueCommentListener { return d.issueCommentListeners }, func(d *Dispatcher, v []IssueCommentListener) { d.issueCommentListeners = v }, IssueCommentListener.OnIssueComment),
	newEventDescriptor(func(d *Dispatcher) []IssueListener { return d.issueListeners }, func(d *Dispatcher, v []IssueListener) { d.issueListeners = v }, IssueListener.OnIssue),
	newEventDescriptor(func(d *Dispatcher) []JobListener { return d.jobListeners }, func(d *Dispatcher, v []JobListener) { d.jobListeners = v }, JobListener.OnJob),
	newEventDescriptor(func(d *Dispatcher) []MemberListener { return d.memberListeners }, func(d *Dispatcher, v []MemberListener) { d.memberListeners = v }, MemberListener.OnMember),
	newEventDescriptor(func(d *Dispatcher) []MilestoneListener { return d.milestoneListeners }, func(d *Dispatcher, v []MilestoneListener) { d.milestoneListeners = v }, MilestoneListener.OnMilestone),
	newEventDescriptor(func(d *Dispatcher) []MergeCommentListener { return d.mergeCommentListeners }, func(d *Dispatcher, v []MergeCommentListener) { d.mergeCommentListeners = v }, MergeCommentListener.OnMergeComment),
	newEventDescriptor(func(d *Dispatcher) []MergeListener { return d.mergeListeners }, func(d *Dispatcher, v []MergeListener) { d.mergeListeners = v }, MergeListener.OnMerge),
	newEventDescriptor(func(d *Dispatcher) []PipelineListener { return d.pipelineListeners }, func(d *Dispatcher, v []PipelineListener) { d.pipelineListeners = v }, PipelineListener.OnPipeline),
	newEventDescriptor(func(d *Dispatcher) []ProjectListener { return d.projectListeners }, func(d *Dispatcher, v []ProjectListener) { d.projectListeners = v }, ProjectListener.OnProject),
	newEventDescriptor(func(d *Dispatcher) []ProjectResourceAccessTokenListener { return d.projectResourceAccessTokenListeners }, func(d *Dispatcher, v []ProjectResourceAccessTokenListener) { d.projectResourceAccessTokenListeners = v }, ProjectResourceAccessTokenListener.OnProjectResourceAccessToken),
	newEventDescriptor(func(d *Dispatcher) []PushListener { return d.pushListeners }, func(d *Dispatcher, v []PushListener) { d.pushListeners = v }, PushListener.OnPush),
	newEventDescriptor(func(d *Dispatcher) []ReleaseListener { return d.releaseListeners }, func(d *Dispatcher, v []ReleaseListener) { d.releaseListeners = v }, ReleaseListener.OnRelease),
	newEventDescriptor(func(d *Dispatcher) []SnippetCommentListener { return d.snippetCommentListeners }, func(d *Dispatcher, v []SnippetCommentListener) { d.snippetCommentListeners = v }, SnippetCommentListener.OnSnippetComment),
	newEventDescriptor(func(d *Dispatcher) []SubGroupListener { return d.subGroupListeners }, func(d *Dispatcher, v []SubGroupListener) { d.subGroupListeners = v }, SubGroupListener.OnSubGroup),
	newEventDescriptor(func(d *Dispatcher) []TagListener { return d.tagListeners }, func(d *Dispatcher, v []TagListener) { d.tagListeners = v }, TagListener.OnTag),
	newEventDescriptor(func(d *Dispatcher) []VulnerabilityListener { return d.vulnerabilityListeners }, func(d *Dispatcher, v []VulnerabilityListener) { d.vulnerabilityListeners = v }, VulnerabilityListener.OnVulnerability),
	newEventDescriptor(func(d *Dispatcher) []WikiPageListener { return d.wikiPageListeners }, func(d *Dispatcher, v []WikiPageListener) { d.wikiPageListeners = v }, WikiPageListener.OnWikiPage),
}

func (d typedEventDescriptor[E, L]) register(dispatcher *Dispatcher, listener any) {
	typedListener, ok := listener.(L)
	if !ok {
		return
	}

	d.set(dispatcher, append(d.listeners(dispatcher), typedListener))
}

func (d typedEventDescriptor[E, L]) dispatch(dispatcher *Dispatcher, ctx context.Context, event any) (bool, error) {
	typedEvent, ok := event.(E)
	if !ok {
		return false, nil
	}

	return true, processEvent(ctx, d.listeners(dispatcher), d.handler, typedEvent)
}

func newEventDescriptor[E any, L any](
	listeners func(*Dispatcher) []L,
	set func(*Dispatcher, []L),
	handler func(L, context.Context, E) error,
) eventDescriptor {
	return typedEventDescriptor[E, L]{listeners: listeners, set: set, handler: handler}
}
