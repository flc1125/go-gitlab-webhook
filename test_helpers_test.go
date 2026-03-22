package gitlabwebhook

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync/atomic"
	"testing"

	gitlab "gitlab.com/gitlab-org/api/client-go/v2"
)

type dispatcherContextKey struct{}

func newDispatcherContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, dispatcherContextKey{}, "DispatcherContext")
}

func loadFixture(t *testing.T, filePath string) []byte {
	t.Helper()

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("read fixture %q: %v", filePath, err)
	}

	return content
}

func newWebhookRequest(t *testing.T, eventType gitlab.EventType, payload []byte) *http.Request {
	t.Helper()

	req, err := http.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(payload))
	if err != nil {
		t.Fatalf("create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Gitlab-Event", string(eventType))

	return req
}

func expectString(field, got, want string) error {
	if got != want {
		return fmt.Errorf("%s: got %q, want %q", field, got, want)
	}

	return nil
}

type eventAssertingListener struct {
	calls atomic.Int32
}

var (
	_ BuildListener                      = (*eventAssertingListener)(nil)
	_ CommitCommentListener              = (*eventAssertingListener)(nil)
	_ DeploymentListener                 = (*eventAssertingListener)(nil)
	_ EmojiListener                      = (*eventAssertingListener)(nil)
	_ FeatureFlagListener                = (*eventAssertingListener)(nil)
	_ GroupResourceAccessTokenListener   = (*eventAssertingListener)(nil)
	_ IssueCommentListener               = (*eventAssertingListener)(nil)
	_ IssueListener                      = (*eventAssertingListener)(nil)
	_ JobListener                        = (*eventAssertingListener)(nil)
	_ MemberListener                     = (*eventAssertingListener)(nil)
	_ MilestoneListener                  = (*eventAssertingListener)(nil)
	_ MergeCommentListener               = (*eventAssertingListener)(nil)
	_ MergeListener                      = (*eventAssertingListener)(nil)
	_ PipelineListener                   = (*eventAssertingListener)(nil)
	_ ProjectListener                    = (*eventAssertingListener)(nil)
	_ ProjectResourceAccessTokenListener = (*eventAssertingListener)(nil)
	_ PushListener                       = (*eventAssertingListener)(nil)
	_ ReleaseListener                    = (*eventAssertingListener)(nil)
	_ SnippetCommentListener             = (*eventAssertingListener)(nil)
	_ SubGroupListener                   = (*eventAssertingListener)(nil)
	_ TagListener                        = (*eventAssertingListener)(nil)
	_ VulnerabilityListener              = (*eventAssertingListener)(nil)
	_ WikiPageListener                   = (*eventAssertingListener)(nil)
)

func (l *eventAssertingListener) validateContext(ctx context.Context) error {
	v, ok := ctx.Value(dispatcherContextKey{}).(string)
	if !ok {
		return fmt.Errorf("dispatcher context missing")
	}

	return expectString("dispatcher context", v, "DispatcherContext")
}

func (l *eventAssertingListener) record(ctx context.Context) error {
	l.calls.Add(1)

	return l.validateContext(ctx)
}

func (l *eventAssertingListener) OnBuild(ctx context.Context, event *gitlab.BuildEvent) error {
	if err := l.record(ctx); err != nil {
		return err
	}

	return expectString("project name", event.ProjectName, "gitlab-org/gitlab-test")
}

func (l *eventAssertingListener) OnCommitComment(ctx context.Context, event *gitlab.CommitCommentEvent) error {
	if err := l.record(ctx); err != nil {
		return err
	}

	return expectString("project name", event.Project.Name, "Gitlab Test")
}

func (l *eventAssertingListener) OnDeployment(ctx context.Context, event *gitlab.DeploymentEvent) error {
	if err := l.record(ctx); err != nil {
		return err
	}

	return expectString("project name", event.Project.Name, "test-deployment-webhooks")
}

func (l *eventAssertingListener) OnEmoji(ctx context.Context, event *gitlab.EmojiEvent) error {
	if err := l.record(ctx); err != nil {
		return err
	}

	if err := expectString("project name", event.Project.Name, "awesome-project"); err != nil {
		return err
	}

	return expectString("emoji name", event.ObjectAttributes.Name, "thumbsup")
}

func (l *eventAssertingListener) OnFeatureFlag(ctx context.Context, event *gitlab.FeatureFlagEvent) error {
	if err := l.record(ctx); err != nil {
		return err
	}

	return expectString("path with namespace", event.Project.PathWithNamespace, "gitlabhq/gitlab-test")
}

func (l *eventAssertingListener) OnGroupResourceAccessToken(ctx context.Context, event *gitlab.GroupResourceAccessTokenEvent) error {
	if err := l.record(ctx); err != nil {
		return err
	}

	return expectString("event name", event.EventName, "expiring_access_token")
}

func (l *eventAssertingListener) OnIssueComment(ctx context.Context, event *gitlab.IssueCommentEvent) error {
	if err := l.record(ctx); err != nil {
		return err
	}

	return expectString("project name", event.Project.Name, "Gitlab Test")
}

func (l *eventAssertingListener) OnIssue(ctx context.Context, event *gitlab.IssueEvent) error {
	if err := l.record(ctx); err != nil {
		return err
	}

	return expectString("project name", event.Project.Name, "Gitlab Test")
}

func (l *eventAssertingListener) OnJob(ctx context.Context, event *gitlab.JobEvent) error {
	if err := l.record(ctx); err != nil {
		return err
	}

	return expectString("build name", event.BuildName, "auto_deploy:start")
}

func (l *eventAssertingListener) OnMember(ctx context.Context, event *gitlab.MemberEvent) error {
	if err := l.record(ctx); err != nil {
		return err
	}

	return expectString("user name", event.UserName, "User1")
}

func (l *eventAssertingListener) OnMilestone(ctx context.Context, event *gitlab.MilestoneWebhookEvent) error {
	if err := l.record(ctx); err != nil {
		return err
	}

	return expectString("project name", event.Project.Name, "Gitlab Test")
}

func (l *eventAssertingListener) OnMergeComment(ctx context.Context, event *gitlab.MergeCommentEvent) error {
	if err := l.record(ctx); err != nil {
		return err
	}

	return expectString("project name", event.Project.Name, "Gitlab Test")
}

func (l *eventAssertingListener) OnMerge(ctx context.Context, event *gitlab.MergeEvent) error {
	if err := l.record(ctx); err != nil {
		return err
	}

	return expectString("project name", event.Project.Name, "Gitlab Test")
}

func (l *eventAssertingListener) OnPipeline(ctx context.Context, event *gitlab.PipelineEvent) error {
	if err := l.record(ctx); err != nil {
		return err
	}

	return expectString("project name", event.Project.Name, "Gitlab Test")
}

func (l *eventAssertingListener) OnProject(ctx context.Context, event *gitlab.ProjectWebhookEvent) error {
	if err := l.record(ctx); err != nil {
		return err
	}

	if err := expectString("project name", event.Name, "Flight"); err != nil {
		return err
	}

	return expectString("path with namespace", event.PathWithNamespace, "flightjs/flight")
}

func (l *eventAssertingListener) OnProjectResourceAccessToken(ctx context.Context, event *gitlab.ProjectResourceAccessTokenEvent) error {
	if err := l.record(ctx); err != nil {
		return err
	}

	return expectString("event name", event.EventName, "expiring_access_token")
}

func (l *eventAssertingListener) OnPush(ctx context.Context, event *gitlab.PushEvent) error {
	if err := l.record(ctx); err != nil {
		return err
	}

	return expectString("path with namespace", event.Project.PathWithNamespace, "mike/diaspora")
}

func (l *eventAssertingListener) OnRelease(ctx context.Context, event *gitlab.ReleaseEvent) error {
	if err := l.record(ctx); err != nil {
		return err
	}

	return expectString("project name", event.Project.Name, "Project Name")
}

func (l *eventAssertingListener) OnSnippetComment(ctx context.Context, event *gitlab.SnippetCommentEvent) error {
	if err := l.record(ctx); err != nil {
		return err
	}

	return expectString("project name", event.Project.Name, "Gitlab Test")
}

func (l *eventAssertingListener) OnSubGroup(ctx context.Context, event *gitlab.SubGroupEvent) error {
	if err := l.record(ctx); err != nil {
		return err
	}

	return expectString("name", event.Name, "SubGroup 1")
}

func (l *eventAssertingListener) OnTag(ctx context.Context, event *gitlab.TagEvent) error {
	if err := l.record(ctx); err != nil {
		return err
	}

	return expectString("project name", event.Project.Name, "Example")
}

func (l *eventAssertingListener) OnVulnerability(ctx context.Context, event *gitlab.VulnerabilityEvent) error {
	if err := l.record(ctx); err != nil {
		return err
	}

	return expectString("title", event.ObjectAttributes.Title, "Potential SQL Injection")
}

func (l *eventAssertingListener) OnWikiPage(ctx context.Context, event *gitlab.WikiPageEvent) error {
	if err := l.record(ctx); err != nil {
		return err
	}

	return expectString("project name", event.Project.Name, "awesome-project")
}

type pushTrackingListener struct {
	called atomic.Int32
}

var _ PushListener = (*pushTrackingListener)(nil)

func (l *pushTrackingListener) OnPush(context.Context, *gitlab.PushEvent) error {
	l.called.Add(1)
	return nil
}

type errorPushListener struct {
	err error
}

var _ PushListener = (*errorPushListener)(nil)

func (l *errorPushListener) OnPush(context.Context, *gitlab.PushEvent) error {
	return l.err
}

type errReader struct{}

func (errReader) Read([]byte) (int, error) {
	return 0, fmt.Errorf("read error")
}

func (errReader) Close() error {
	return nil
}

var _ io.ReadCloser = errReader{}
