package eventmeta

import (
	"reflect"
	"strings"

	"github.com/flc1125/go-gitlab-webhook/middleware/otel/v3/internal/semconv"
	"gitlab.com/gitlab-org/api/client-go/v2"
	"go.opentelemetry.io/otel/attribute"
)

// Metadata describes how a GitLab webhook event should be represented in a span.
type Metadata struct {
	SpanName   string
	Attributes []attribute.KeyValue
}

// Extract reads supported GitLab webhook event structs and returns stable span
// metadata. Span names only use low-cardinality values; high-cardinality data
// such as project paths, refs, IDs, and names are kept as attributes.
func Extract(event any) Metadata { //nolint:cyclop
	switch e := event.(type) {
	case *gitlab.BuildEvent:
		return newBuilder("build").
			objectKind(e.ObjectKind).
			project(e.ProjectID, "", e.ProjectName).
			ref(e.Ref).
			sha(e.SHA).
			job(e.BuildID, e.BuildName, e.BuildStage).
			status(e.BuildStatus).
			build()
	case *gitlab.CommitCommentEvent:
		return newBuilder("note").
			objectKind(e.ObjectKind).
			project(e.ProjectID, e.Project.PathWithNamespace, e.Project.Name).
			note(e.ObjectAttributes.NoteableType, string(e.ObjectAttributes.Action)).
			build()
	case *gitlab.DeploymentEvent:
		return newBuilder("deployment").
			objectKind(e.ObjectKind).
			project(e.Project.ID, e.Project.PathWithNamespace, e.Project.Name).
			ref(e.Ref).
			status(e.Status).
			build()
	case *gitlab.EmojiEvent:
		return newBuilder("emoji").
			objectKind(e.ObjectKind).
			project(e.ProjectID, e.Project.PathWithNamespace, e.Project.Name).
			action(e.ObjectAttributes.Action).
			attrString(semconv.EmojiAwardableType, e.ObjectAttributes.AwardableType).
			build()
	case *gitlab.FeatureFlagEvent:
		return newBuilder("feature_flag").
			objectKind(e.ObjectKind).
			project(e.Project.ID, e.Project.PathWithNamespace, e.Project.Name).
			attrInt64(semconv.FeatureFlagID, e.ObjectAttributes.ID).
			attrString(semconv.FeatureFlagName, e.ObjectAttributes.Name).
			attrBool(semconv.FeatureFlagActive, e.ObjectAttributes.Active).
			build()
	case *gitlab.GroupResourceAccessTokenEvent:
		return newBuilder("resource_access_token").
			objectKind(e.ObjectKind).
			eventName(e.EventName).
			operation(e.EventName).
			group(e.Group.GroupID, e.Group.FullPath, e.Group.GroupName).
			attrInt64(semconv.ResourceAccessTokenID, e.ObjectAttributes.ID).
			build()
	case *gitlab.IssueCommentEvent:
		return newBuilder("note").
			objectKind(e.ObjectKind).
			project(e.ProjectID, e.Project.PathWithNamespace, e.Project.Name).
			note(e.ObjectAttributes.NoteableType, string(e.ObjectAttributes.Action)).
			issue(e.Issue.IID, e.Issue.State).
			build()
	case *gitlab.IssueEvent:
		return newBuilder("issue").
			objectKind(e.ObjectKind).
			project(e.Project.ID, e.Project.PathWithNamespace, e.Project.Name).
			action(e.ObjectAttributes.Action).
			issue(e.ObjectAttributes.IID, e.ObjectAttributes.State).
			build()
	case *gitlab.JobEvent:
		return newBuilder("job").
			objectKind(e.ObjectKind).
			project(e.ProjectID, "", e.ProjectName).
			ref(e.Ref).
			sha(e.SHA).
			job(e.BuildID, e.BuildName, e.BuildStage).
			status(e.BuildStatus).
			build()
	case *gitlab.MemberEvent:
		return newBuilder("member").
			eventName(e.EventName).
			operation(e.EventName).
			group(e.GroupID, e.GroupPath, e.GroupName).
			build()
	case *gitlab.MilestoneWebhookEvent:
		b := newBuilder("milestone").
			objectKind(e.ObjectKind).
			project(e.Project.ID, e.Project.PathWithNamespace, e.Project.Name).
			action(e.Action).
			attrInt64(semconv.MilestoneIID, e.ObjectAttributes.IID).
			attrString(semconv.MilestoneState, e.ObjectAttributes.State)
		if e.Group != nil {
			b.group(e.Group.GroupID, e.Group.FullPath, e.Group.GroupName)
		}
		return b.build()
	case *gitlab.MergeCommentEvent:
		return newBuilder("note").
			objectKind(e.ObjectKind).
			project(e.ProjectID, e.Project.PathWithNamespace, e.Project.Name).
			note(e.ObjectAttributes.NoteableType, string(e.ObjectAttributes.Action)).
			mergeRequest(e.MergeRequest.IID, e.MergeRequest.State).
			build()
	case *gitlab.MergeEvent:
		return newBuilder("merge_request").
			objectKind(e.ObjectKind).
			project(e.Project.ID, e.Project.PathWithNamespace, e.Project.Name).
			action(e.ObjectAttributes.Action).
			mergeRequest(e.ObjectAttributes.IID, e.ObjectAttributes.State).
			attrString(semconv.MergeRequestMergeStatus, e.ObjectAttributes.MergeStatus).
			attrString(semconv.MergeRequestDetailedMergeStatus, e.ObjectAttributes.DetailedMergeStatus).
			build()
	case *gitlab.PipelineEvent:
		return newBuilder("pipeline").
			objectKind(e.ObjectKind).
			project(e.Project.ID, e.Project.PathWithNamespace, e.Project.Name).
			ref(e.ObjectAttributes.Ref).
			sha(e.ObjectAttributes.SHA).
			pipeline(e.ObjectAttributes.ID, e.ObjectAttributes.IID, e.ObjectAttributes.Source).
			status(e.ObjectAttributes.Status).
			build()
	case *gitlab.ProjectWebhookEvent:
		return newBuilder("project").
			eventName(e.EventName).
			operation(e.EventName).
			project(e.ProjectID, e.PathWithNamespace, e.Name).
			attrInt64(semconv.ProjectNamespaceID, e.ProjectNamespaceID).
			attrString(semconv.ProjectVisibility, e.ProjectVisibility).
			build()
	case *gitlab.ProjectResourceAccessTokenEvent:
		return newBuilder("resource_access_token").
			objectKind(e.ObjectKind).
			eventName(e.EventName).
			operation(e.EventName).
			project(e.Project.ID, e.Project.PathWithNamespace, e.Project.Name).
			attrInt64(semconv.ResourceAccessTokenID, e.ObjectAttributes.ID).
			build()
	case *gitlab.PushEvent:
		return newBuilder("push").
			objectKind(e.ObjectKind).
			eventName(e.EventName).
			project(e.ProjectID, e.Project.PathWithNamespace, e.Project.Name).
			ref(e.Ref).
			sha(e.CheckoutSHA).
			attrInt64(semconv.PushTotalCommitsCount, e.TotalCommitsCount).
			build()
	case *gitlab.ReleaseEvent:
		return newBuilder("release").
			objectKind(e.ObjectKind).
			project(e.Project.ID, e.Project.PathWithNamespace, e.Project.Name).
			action(e.Action).
			attrString(semconv.ReleaseTag, e.Tag).
			build()
	case *gitlab.SnippetCommentEvent:
		return newBuilder("note").
			objectKind(e.ObjectKind).
			project(e.ProjectID, e.Project.PathWithNamespace, e.Project.Name).
			note(e.ObjectAttributes.NoteableType, string(e.ObjectAttributes.Action)).
			build()
	case *gitlab.SubGroupEvent:
		return newBuilder("subgroup").
			eventName(e.EventName).
			operation(e.EventName).
			group(e.GroupID, e.FullPath, e.Name).
			attrInt64(semconv.ParentGroupID, e.ParentGroupID).
			attrString(semconv.ParentGroupPath, e.ParentFullPath).
			build()
	case *gitlab.TagEvent:
		return newBuilder("tag_push").
			objectKind(e.ObjectKind).
			eventName(e.EventName).
			project(e.ProjectID, e.Project.PathWithNamespace, e.Project.Name).
			ref(e.Ref).
			sha(e.CheckoutSHA).
			attrInt64(semconv.PushTotalCommitsCount, e.TotalCommitsCount).
			build()
	case *gitlab.VulnerabilityEvent:
		return newBuilder("vulnerability").
			objectKind(e.ObjectKind).
			project(e.ObjectAttributes.ProjectID, "", "").
			status(e.ObjectAttributes.State).
			attrString(semconv.VulnerabilityLevel, e.ObjectAttributes.Severity).
			attrString(semconv.VulnerabilityReport, e.ObjectAttributes.ReportType).
			attrString(semconv.VulnerabilityState, e.ObjectAttributes.State).
			build()
	case *gitlab.WikiPageEvent:
		return newBuilder("wiki_page").
			objectKind(e.ObjectKind).
			project(0, e.Project.PathWithNamespace, e.Project.Name).
			action(e.ObjectAttributes.Action).
			build()
	default:
		return fallback(event)
	}
}

type builder struct {
	eventType string
	spanName  []string
	attrs     []attribute.KeyValue
}

func newBuilder(eventType string) *builder {
	return &builder{
		eventType: eventType,
		attrs: []attribute.KeyValue{
			semconv.WebhookEventType(eventType),
		},
	}
}

func (b *builder) objectKind(value string) *builder {
	return b.attrString(semconv.WebhookObjectKind, value)
}

func (b *builder) eventName(value string) *builder {
	return b.attrString(semconv.WebhookEventName, value)
}

func (b *builder) project(id int64, path, name string) *builder {
	return b.
		attrInt64(semconv.ProjectID, id).
		attrString(semconv.ProjectPath, path).
		attrString(semconv.ProjectName, name)
}

func (b *builder) group(id int64, path, name string) *builder {
	return b.
		attrInt64(semconv.GroupID, id).
		attrString(semconv.GroupPath, path).
		attrString(semconv.GroupName, name)
}

func (b *builder) ref(value string) *builder {
	return b.attrString(semconv.Ref, value)
}

func (b *builder) sha(value string) *builder {
	return b.attrString(semconv.SHA, value)
}

func (b *builder) action(value string) *builder {
	b.operation(value)
	return b.attrString(semconv.WebhookAction, value)
}

func (b *builder) status(value string) *builder {
	b.operation(value)
	return b.attrString(semconv.WebhookStatus, value)
}

func (b *builder) note(noteableType, action string) *builder {
	b.operation(noteableType)
	b.operation(action)
	return b.attrString(semconv.NoteNoteableType, noteableType).attrString(semconv.WebhookAction, action)
}

func (b *builder) issue(iid int64, state string) *builder {
	return b.attrInt64(semconv.IssueIID, iid).attrString(semconv.IssueState, state)
}

func (b *builder) job(id int64, name, stage string) *builder {
	return b.
		attrInt64(semconv.JobID, id).
		attrString(semconv.JobName, name).
		attrString(semconv.JobStage, stage)
}

func (b *builder) mergeRequest(iid int64, state string) *builder {
	return b.attrInt64(semconv.MergeRequestIID, iid).attrString(semconv.MergeRequestState, state)
}

func (b *builder) pipeline(id, iid int64, source string) *builder {
	return b.
		attrInt64(semconv.PipelineID, id).
		attrInt64(semconv.PipelineIID, iid).
		attrString(semconv.PipelineSource, source)
}

func (b *builder) attrString(keyValue func(string) attribute.KeyValue, value string) *builder {
	value = strings.TrimSpace(value)
	if value != "" {
		b.attrs = append(b.attrs, keyValue(value))
	}

	return b
}

func (b *builder) attrInt64(keyValue func(int64) attribute.KeyValue, value int64) *builder {
	if value != 0 {
		b.attrs = append(b.attrs, keyValue(value))
	}

	return b
}

func (b *builder) attrBool(keyValue func(bool) attribute.KeyValue, value bool) *builder {
	b.attrs = append(b.attrs, keyValue(value))

	return b
}

func (b *builder) operation(value string) *builder {
	value = normalizeOperation(value)
	if value != "" {
		b.spanName = append(b.spanName, value)
	}

	return b
}

func normalizeOperation(value string) string {
	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(value, " ", "_")
	value = strings.ReplaceAll(value, "-", "_")
	value = strings.ReplaceAll(value, "/", "_")

	var normalized strings.Builder
	for i := 0; i < len(value); i++ {
		c := value[i]
		if c >= 'A' && c <= 'Z' {
			if i > 0 && value[i-1] != '_' {
				normalized.WriteByte('_')
			}
			normalized.WriteByte(c + ('a' - 'A'))
			continue
		}
		normalized.WriteByte(c)
	}

	return normalized.String()
}

func (b *builder) build() Metadata {
	name := "gitlab.webhook." + b.eventType
	if len(b.spanName) > 0 {
		name += " " + strings.Join(b.spanName, " ")
	}

	return Metadata{
		SpanName:   name,
		Attributes: b.attrs,
	}
}

func fallback(event any) Metadata {
	b := newBuilder("unknown")
	if event != nil {
		b.attrString(semconv.WebhookGoType, goTypeName(event))
	}

	return b.build()
}

func goTypeName(event any) string {
	t := reflect.TypeOf(event)
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	if t.Name() != "" {
		return t.PkgPath() + "." + t.Name()
	}

	return t.String()
}
