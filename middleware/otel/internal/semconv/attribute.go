package semconv

import "go.opentelemetry.io/otel/attribute"

const (
	WebhookEventTypeKey  = attribute.Key("gitlab.webhook.event_type")
	WebhookObjectKindKey = attribute.Key("gitlab.webhook.object_kind")
	WebhookEventNameKey  = attribute.Key("gitlab.webhook.event_name")
	WebhookActionKey     = attribute.Key("gitlab.webhook.action")
	WebhookStatusKey     = attribute.Key("gitlab.webhook.status")
	WebhookResultKey     = attribute.Key("gitlab.webhook.result")
	WebhookGoTypeKey     = attribute.Key("gitlab.webhook.go_type")

	ProjectIDKey          = attribute.Key("gitlab.project.id")
	ProjectNameKey        = attribute.Key("gitlab.project.name")
	ProjectPathKey        = attribute.Key("gitlab.project.path")
	ProjectNamespaceIDKey = attribute.Key("gitlab.project.namespace_id")
	ProjectVisibilityKey  = attribute.Key("gitlab.project.visibility")

	GroupIDKey   = attribute.Key("gitlab.group.id")
	GroupNameKey = attribute.Key("gitlab.group.name")
	GroupPathKey = attribute.Key("gitlab.group.path")

	ParentGroupIDKey   = attribute.Key("gitlab.parent_group.id")
	ParentGroupPathKey = attribute.Key("gitlab.parent_group.path")

	RefKey = attribute.Key("gitlab.ref")
	SHAKey = attribute.Key("gitlab.sha")

	EmojiAwardableTypeKey = attribute.Key("gitlab.emoji.awardable_type")

	FeatureFlagIDKey     = attribute.Key("gitlab.feature_flag.id")
	FeatureFlagNameKey   = attribute.Key("gitlab.feature_flag.name")
	FeatureFlagActiveKey = attribute.Key("gitlab.feature_flag.active")

	IssueIIDKey   = attribute.Key("gitlab.issue.iid")
	IssueStateKey = attribute.Key("gitlab.issue.state")

	JobIDKey    = attribute.Key("gitlab.job.id")
	JobNameKey  = attribute.Key("gitlab.job.name")
	JobStageKey = attribute.Key("gitlab.job.stage")

	MergeRequestIIDKey                 = attribute.Key("gitlab.merge_request.iid")
	MergeRequestStateKey               = attribute.Key("gitlab.merge_request.state")
	MergeRequestMergeStatusKey         = attribute.Key("gitlab.merge_request.merge_status")
	MergeRequestDetailedMergeStatusKey = attribute.Key("gitlab.merge_request.detailed_merge_status")

	MilestoneIIDKey   = attribute.Key("gitlab.milestone.iid")
	MilestoneStateKey = attribute.Key("gitlab.milestone.state")

	NoteNoteableTypeKey = attribute.Key("gitlab.note.noteable_type")

	PipelineIDKey     = attribute.Key("gitlab.pipeline.id")
	PipelineIIDKey    = attribute.Key("gitlab.pipeline.iid")
	PipelineSourceKey = attribute.Key("gitlab.pipeline.source")

	PushTotalCommitsCountKey = attribute.Key("gitlab.push.total_commits_count")

	ReleaseTagKey = attribute.Key("gitlab.release.tag")

	ResourceAccessTokenIDKey = attribute.Key("gitlab.resource_access_token.id")

	VulnerabilityStateKey  = attribute.Key("gitlab.vulnerability.state")
	VulnerabilityReportKey = attribute.Key("gitlab.vulnerability.report_type")
	VulnerabilityLevelKey  = attribute.Key("gitlab.vulnerability.severity")
)

func WebhookEventType(val string) attribute.KeyValue {
	return WebhookEventTypeKey.String(val)
}

func WebhookObjectKind(val string) attribute.KeyValue {
	return WebhookObjectKindKey.String(val)
}

func WebhookEventName(val string) attribute.KeyValue {
	return WebhookEventNameKey.String(val)
}

func WebhookAction(val string) attribute.KeyValue {
	return WebhookActionKey.String(val)
}

func WebhookStatus(val string) attribute.KeyValue {
	return WebhookStatusKey.String(val)
}

func WebhookResult(val string) attribute.KeyValue {
	return WebhookResultKey.String(val)
}

func WebhookGoType(val string) attribute.KeyValue {
	return WebhookGoTypeKey.String(val)
}

func ProjectID(val int64) attribute.KeyValue {
	return ProjectIDKey.Int64(val)
}

func ProjectName(val string) attribute.KeyValue {
	return ProjectNameKey.String(val)
}

func ProjectPath(val string) attribute.KeyValue {
	return ProjectPathKey.String(val)
}

func ProjectNamespaceID(val int64) attribute.KeyValue {
	return ProjectNamespaceIDKey.Int64(val)
}

func ProjectVisibility(val string) attribute.KeyValue {
	return ProjectVisibilityKey.String(val)
}

func GroupID(val int64) attribute.KeyValue {
	return GroupIDKey.Int64(val)
}

func GroupName(val string) attribute.KeyValue {
	return GroupNameKey.String(val)
}

func GroupPath(val string) attribute.KeyValue {
	return GroupPathKey.String(val)
}

func ParentGroupID(val int64) attribute.KeyValue {
	return ParentGroupIDKey.Int64(val)
}

func ParentGroupPath(val string) attribute.KeyValue {
	return ParentGroupPathKey.String(val)
}

func Ref(val string) attribute.KeyValue {
	return RefKey.String(val)
}

func SHA(val string) attribute.KeyValue {
	return SHAKey.String(val)
}

func EmojiAwardableType(val string) attribute.KeyValue {
	return EmojiAwardableTypeKey.String(val)
}

func FeatureFlagID(val int64) attribute.KeyValue {
	return FeatureFlagIDKey.Int64(val)
}

func FeatureFlagName(val string) attribute.KeyValue {
	return FeatureFlagNameKey.String(val)
}

func FeatureFlagActive(val bool) attribute.KeyValue {
	return FeatureFlagActiveKey.Bool(val)
}

func IssueIID(val int64) attribute.KeyValue {
	return IssueIIDKey.Int64(val)
}

func IssueState(val string) attribute.KeyValue {
	return IssueStateKey.String(val)
}

func JobID(val int64) attribute.KeyValue {
	return JobIDKey.Int64(val)
}

func JobName(val string) attribute.KeyValue {
	return JobNameKey.String(val)
}

func JobStage(val string) attribute.KeyValue {
	return JobStageKey.String(val)
}

func MergeRequestIID(val int64) attribute.KeyValue {
	return MergeRequestIIDKey.Int64(val)
}

func MergeRequestState(val string) attribute.KeyValue {
	return MergeRequestStateKey.String(val)
}

func MergeRequestMergeStatus(val string) attribute.KeyValue {
	return MergeRequestMergeStatusKey.String(val)
}

func MergeRequestDetailedMergeStatus(val string) attribute.KeyValue {
	return MergeRequestDetailedMergeStatusKey.String(val)
}

func MilestoneIID(val int64) attribute.KeyValue {
	return MilestoneIIDKey.Int64(val)
}

func MilestoneState(val string) attribute.KeyValue {
	return MilestoneStateKey.String(val)
}

func NoteNoteableType(val string) attribute.KeyValue {
	return NoteNoteableTypeKey.String(val)
}

func PipelineID(val int64) attribute.KeyValue {
	return PipelineIDKey.Int64(val)
}

func PipelineIID(val int64) attribute.KeyValue {
	return PipelineIIDKey.Int64(val)
}

func PipelineSource(val string) attribute.KeyValue {
	return PipelineSourceKey.String(val)
}

func PushTotalCommitsCount(val int64) attribute.KeyValue {
	return PushTotalCommitsCountKey.Int64(val)
}

func ReleaseTag(val string) attribute.KeyValue {
	return ReleaseTagKey.String(val)
}

func ResourceAccessTokenID(val int64) attribute.KeyValue {
	return ResourceAccessTokenIDKey.Int64(val)
}

func VulnerabilityState(val string) attribute.KeyValue {
	return VulnerabilityStateKey.String(val)
}

func VulnerabilityReport(val string) attribute.KeyValue {
	return VulnerabilityReportKey.String(val)
}

func VulnerabilityLevel(val string) attribute.KeyValue {
	return VulnerabilityLevelKey.String(val)
}
