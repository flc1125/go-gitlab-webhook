package gitlabwebhook

type listenerRegistration func(*Dispatcher, any)

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

func registerListenerType[L any](register func(*Dispatcher, L)) listenerRegistration {
	return func(d *Dispatcher, listener any) {
		typedListener, ok := listener.(L)
		if !ok {
			return
		}

		register(d, typedListener)
	}
}
