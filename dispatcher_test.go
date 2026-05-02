package gitlabwebhook

import (
	"bytes"
	"context"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	gitlab "gitlab.com/gitlab-org/api/client-go/v2"
)

type dispatcherContextKey struct{}

func newDispatcherContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, dispatcherContextKey{}, "DispatcherContext")
}

func testDispatcherContext(ctx context.Context, t *testing.T) {
	v, ok := ctx.Value(dispatcherContextKey{}).(string)
	assert.True(t, ok)
	assert.Equal(t, "DispatcherContext", v)
}

func TestDispatcher_Dispatch(t *testing.T) {
	dispatcher := NewDispatcher(
		RegisterListeners(&testListener{t: t}),
	)
	dispatcher.RegisterListeners(&testListener{t: t})

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/webhook", r.URL.Path)
		assert.NoError(t,
			dispatcher.DispatchRequest(r,
				DispatchRequestWithContext(newDispatcherContext(r.Context())),
			),
		)

		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))
	defer srv.Close()

	tests := []struct {
		name      string
		eventType gitlab.EventType
		filepath  string
	}{
		{"build", gitlab.EventTypeBuild, "webhooks/build.json"},                //nolint:lll
		{"commit comment", gitlab.EventTypeNote, "webhooks/note_commit.json"},  //nolint:lll
		{"deployment", gitlab.EventTypeDeployment, "webhooks/deployment.json"}, //nolint:lll
		{"emoji", gitlab.EventTypeEmoji, "webhooks/emoji.json"},
		{"feature flag", gitlab.EventTypeFeatureFlag, "webhooks/feature_flag.json"},                                           //nolint:lll
		{"group resource access token", gitlab.EventTypeResourceAccessToken, "webhooks/resource_access_token_group.json"},     //nolint:lll
		{"project resource access token", gitlab.EventTypeResourceAccessToken, "webhooks/resource_access_token_project.json"}, //nolint:lll
		{"issue comment", gitlab.EventTypeNote, "webhooks/note_issue.json"},                                                   //nolint:lll
		{"issue", gitlab.EventTypeIssue, "webhooks/issue.json"},                                                               //nolint:lll
		{"job", gitlab.EventTypeJob, "webhooks/job.json"},
		{"member", gitlab.EventTypeMember, "webhooks/member.json"},
		{"milestone", gitlab.EventTypeMilestone, "webhooks/milestone.json"},
		{"merge comment", gitlab.EventTypeNote, "webhooks/note_merge_request.json"}, //nolint:lll
		{"merge", gitlab.EventTypeMergeRequest, "webhooks/merge_request.json"},      //nolint:lll
		{"pipeline", gitlab.EventTypePipeline, "webhooks/pipeline.json"},            //nolint:lll
		{"project", gitlab.EventTypeProject, "webhooks/project.json"},
		{"push", gitlab.EventTypePush, "webhooks/push.json"},
		{"release", gitlab.EventTypeRelease, "webhooks/release.json"},           //nolint:lll
		{"snippet comment", gitlab.EventTypeNote, "webhooks/note_snippet.json"}, //nolint:lll
		{"subgroup", gitlab.EventTypeSubGroup, "webhooks/subgroup.json"},        //nolint:lll
		{"tag", gitlab.EventTypeTagPush, "webhooks/tag_push.json"},
		{"vulnerability", gitlab.EventTypeVulnerability, "webhooks/vulnerability.json"}, //nolint:lll
		{"wiki page", gitlab.EventTypeWikiPage, "webhooks/wiki_page.json"},              //nolint:lll
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, srv.URL+"/webhook", bytes.NewReader(loadFixture(tt.filepath)))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Gitlab-Event", string(tt.eventType))

			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close() //nolint:errcheck

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var buf bytes.Buffer
			_, _ = buf.ReadFrom(resp.Body)
			assert.Equal(t, `{"status":"ok"}`, buf.String())
		})
	}
}

type testListener struct {
	t *testing.T
}

var (
	_ BuildListener                      = (*testListener)(nil)
	_ CommitCommentListener              = (*testListener)(nil)
	_ DeploymentListener                 = (*testListener)(nil)
	_ EmojiListener                      = (*testListener)(nil)
	_ FeatureFlagListener                = (*testListener)(nil)
	_ GroupResourceAccessTokenListener   = (*testListener)(nil)
	_ IssueCommentListener               = (*testListener)(nil)
	_ IssueListener                      = (*testListener)(nil)
	_ JobListener                        = (*testListener)(nil)
	_ MemberListener                     = (*testListener)(nil)
	_ MilestoneListener                  = (*testListener)(nil)
	_ MergeCommentListener               = (*testListener)(nil)
	_ MergeListener                      = (*testListener)(nil)
	_ PipelineListener                   = (*testListener)(nil)
	_ ProjectListener                    = (*testListener)(nil)
	_ ProjectResourceAccessTokenListener = (*testListener)(nil)
	_ PushListener                       = (*testListener)(nil)
	_ ReleaseListener                    = (*testListener)(nil)
	_ SnippetCommentListener             = (*testListener)(nil)
	_ SubGroupListener                   = (*testListener)(nil)
	_ TagListener                        = (*testListener)(nil)
	_ VulnerabilityListener              = (*testListener)(nil)
	_ WikiPageListener                   = (*testListener)(nil)
)

func (t *testListener) OnBuild(ctx context.Context, event *gitlab.BuildEvent) error {
	testDispatcherContext(ctx, t.t)
	assert.Equal(t.t, "gitlab-org/gitlab-test", event.ProjectName)
	return nil
}

func (t *testListener) OnCommitComment(ctx context.Context, event *gitlab.CommitCommentEvent) error {
	testDispatcherContext(ctx, t.t)
	assert.Equal(t.t, "Gitlab Test", event.Project.Name)
	return nil
}

func (t *testListener) OnDeployment(ctx context.Context, event *gitlab.DeploymentEvent) error {
	testDispatcherContext(ctx, t.t)
	assert.Equal(t.t, "test-deployment-webhooks", event.Project.Name)
	return nil
}

func (t *testListener) OnEmoji(ctx context.Context, event *gitlab.EmojiEvent) error {
	testDispatcherContext(ctx, t.t)
	assert.Equal(t.t, "awesome-project", event.Project.Name)
	assert.Equal(t.t, "thumbsup", event.ObjectAttributes.Name)
	return nil
}

func (t *testListener) OnFeatureFlag(ctx context.Context, event *gitlab.FeatureFlagEvent) error {
	testDispatcherContext(ctx, t.t)
	assert.Equal(t.t, "gitlabhq/gitlab-test", event.Project.PathWithNamespace)
	return nil
}

func (t *testListener) OnGroupResourceAccessToken(ctx context.Context, event *gitlab.GroupResourceAccessTokenEvent) error { //nolint:lll
	testDispatcherContext(ctx, t.t)
	assert.Equal(t.t, "expiring_access_token", event.EventName)
	return nil
}

func (t *testListener) OnIssueComment(ctx context.Context, event *gitlab.IssueCommentEvent) error {
	testDispatcherContext(ctx, t.t)
	assert.Equal(t.t, "Gitlab Test", event.Project.Name)
	return nil
}

func (t *testListener) OnIssue(ctx context.Context, event *gitlab.IssueEvent) error {
	testDispatcherContext(ctx, t.t)
	assert.Equal(t.t, "Gitlab Test", event.Project.Name)
	return nil
}

func (t *testListener) OnJob(ctx context.Context, event *gitlab.JobEvent) error {
	testDispatcherContext(ctx, t.t)
	assert.Equal(t.t, "auto_deploy:start", event.BuildName)
	return nil
}

func (t *testListener) OnMember(ctx context.Context, event *gitlab.MemberEvent) error {
	testDispatcherContext(ctx, t.t)
	assert.Equal(t.t, "User1", event.UserName)
	return nil
}

func (t *testListener) OnMilestone(ctx context.Context, event *gitlab.MilestoneWebhookEvent) error {
	testDispatcherContext(ctx, t.t)
	assert.Equal(t.t, "Gitlab Test", event.Project.Name)
	return nil
}

func (t *testListener) OnMergeComment(ctx context.Context, event *gitlab.MergeCommentEvent) error {
	testDispatcherContext(ctx, t.t)
	assert.Equal(t.t, "Gitlab Test", event.Project.Name)
	return nil
}

func (t *testListener) OnMerge(ctx context.Context, event *gitlab.MergeEvent) error {
	testDispatcherContext(ctx, t.t)
	assert.Equal(t.t, "Gitlab Test", event.Project.Name)
	return nil
}

func (t *testListener) OnPipeline(ctx context.Context, event *gitlab.PipelineEvent) error {
	testDispatcherContext(ctx, t.t)
	assert.Equal(t.t, "Gitlab Test", event.Project.Name)
	return nil
}

func (t *testListener) OnProject(ctx context.Context, event *gitlab.ProjectWebhookEvent) error {
	testDispatcherContext(ctx, t.t)
	assert.Equal(t.t, "Flight", event.Name)
	assert.Equal(t.t, "flightjs/flight", event.PathWithNamespace)
	return nil
}

func (t *testListener) OnProjectResourceAccessToken(ctx context.Context, event *gitlab.ProjectResourceAccessTokenEvent) error { //nolint:lll
	testDispatcherContext(ctx, t.t)
	assert.Equal(t.t, "expiring_access_token", event.EventName)
	return nil
}

func (t *testListener) OnPush(ctx context.Context, event *gitlab.PushEvent) error {
	testDispatcherContext(ctx, t.t)
	assert.Equal(t.t, "mike/diaspora", event.Project.PathWithNamespace)
	return nil
}

func (t *testListener) OnRelease(ctx context.Context, event *gitlab.ReleaseEvent) error {
	testDispatcherContext(ctx, t.t)
	assert.Equal(t.t, "Project Name", event.Project.Name)
	return nil
}

func (t *testListener) OnSnippetComment(ctx context.Context, event *gitlab.SnippetCommentEvent) error {
	testDispatcherContext(ctx, t.t)
	assert.Equal(t.t, "Gitlab Test", event.Project.Name)
	return nil
}

func (t *testListener) OnSubGroup(ctx context.Context, event *gitlab.SubGroupEvent) error {
	testDispatcherContext(ctx, t.t)
	assert.Equal(t.t, "SubGroup 1", event.Name)
	return nil
}

func (t *testListener) OnTag(ctx context.Context, event *gitlab.TagEvent) error {
	testDispatcherContext(ctx, t.t)
	assert.Equal(t.t, "Example", event.Project.Name)
	return nil
}

func (t *testListener) OnVulnerability(ctx context.Context, event *gitlab.VulnerabilityEvent) error {
	testDispatcherContext(ctx, t.t)
	assert.Equal(t.t, "Potential SQL Injection", event.ObjectAttributes.Title)
	return nil
}

func (t *testListener) OnWikiPage(ctx context.Context, event *gitlab.WikiPageEvent) error {
	testDispatcherContext(ctx, t.t)
	assert.Equal(t.t, "awesome-project", event.Project.Name)
	return nil
}

func TestDispatcher_DispatchRequestWithToken(t *testing.T) {
	simpleListener := &simpleTestListener{}
	dispatcher := NewDispatcher(
		RegisterListeners(simpleListener),
	)

	validToken := "test-secret-token"
	invalidToken := "wrong-token"

	tests := []struct {
		name           string
		token          string
		headerToken    string
		expectedError  error
		shouldDispatch bool
	}{
		{
			name:           "valid token should dispatch successfully",
			token:          validToken,
			headerToken:    validToken,
			expectedError:  nil,
			shouldDispatch: true,
		},
		{
			name:           "invalid token should return ErrInvalidToken",
			token:          validToken,
			headerToken:    invalidToken,
			expectedError:  ErrInvalidToken,
			shouldDispatch: false,
		},
		{
			name:           "missing token header should return ErrInvalidToken",
			token:          validToken,
			headerToken:    "",
			expectedError:  ErrInvalidToken,
			shouldDispatch: false,
		},
		{
			name:           "no token provided should dispatch successfully",
			token:          "",
			headerToken:    "",
			expectedError:  nil,
			shouldDispatch: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(loadFixture("webhooks/push.json")))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Gitlab-Event", string(gitlab.EventTypePush))

			if tt.headerToken != "" {
				req.Header.Set("X-Gitlab-Token", tt.headerToken)
			}

			var opts []DispatchRequestOption
			if tt.token != "" {
				opts = append(opts, DispatchRequestWithToken(tt.token))
			}

			// Reset listener call count
			simpleListener.called = false

			err = dispatcher.DispatchRequest(req, opts...)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedError)
				assert.False(t, simpleListener.called, "Listener should not be called when token validation fails")
			} else {
				assert.NoError(t, err)
				assert.True(t, simpleListener.called, "Listener should be called when token validation passes")
			}
		})
	}
}

type simpleTestListener struct {
	called bool
}

var _ PushListener = (*simpleTestListener)(nil)

func (s *simpleTestListener) OnPush(ctx context.Context, event *gitlab.PushEvent) error {
	s.called = true
	return nil
}

type middlewareContextKey struct{}

func TestDispatcher_Middleware(t *testing.T) {
	t.Run("wraps listener dispatch", func(t *testing.T) {
		order := make([]string, 0, 3)
		listener := &middlewareTestListener{
			onPush: func(ctx context.Context, event *gitlab.PushEvent) error {
				order = append(order, "listener")
				assert.Equal(t, "middleware", ctx.Value(middlewareContextKey{}))
				return nil
			},
		}

		dispatcher := NewDispatcher(
			RegisterListeners(listener),
			WithMiddlewares(func(next HandlerFunc) HandlerFunc {
				return func(ctx context.Context, event any) error {
					order = append(order, "middleware before")
					ctx = context.WithValue(ctx, middlewareContextKey{}, "middleware")
					err := next(ctx, event)
					order = append(order, "middleware after")
					return err
				}
			}),
		)

		err := dispatcher.Dispatch(t.Context(), &gitlab.PushEvent{})

		assert.NoError(t, err)
		assert.Equal(t, []string{"middleware before", "listener", "middleware after"}, order)
	})

	t.Run("uses registration order", func(t *testing.T) {
		order := make([]string, 0, 4)
		dispatcher := NewDispatcher()
		dispatcher.Use(
			func(next HandlerFunc) HandlerFunc {
				return func(ctx context.Context, event any) error {
					order = append(order, "first before")
					err := next(ctx, event)
					order = append(order, "first after")
					return err
				}
			},
			func(next HandlerFunc) HandlerFunc {
				return func(ctx context.Context, event any) error {
					order = append(order, "second before")
					err := next(ctx, event)
					order = append(order, "second after")
					return err
				}
			},
		)

		err := dispatcher.Dispatch(t.Context(), &gitlab.PushEvent{})

		assert.NoError(t, err)
		assert.Equal(t, []string{"first before", "second before", "second after", "first after"}, order)
	})

	t.Run("can stop dispatch", func(t *testing.T) {
		expectedErr := errors.New("stop dispatch")
		listener := &middlewareTestListener{
			onPush: func(ctx context.Context, event *gitlab.PushEvent) error {
				t.Fatal("listener should not be called")
				return nil
			},
		}

		dispatcher := NewDispatcher(
			RegisterListeners(listener),
			WithMiddlewares(func(next HandlerFunc) HandlerFunc {
				return func(ctx context.Context, event any) error {
					return expectedErr
				}
			}),
		)

		err := dispatcher.Dispatch(t.Context(), &gitlab.PushEvent{})

		assert.ErrorIs(t, err, expectedErr)
	})

	t.Run("for event runs only for matching events", func(t *testing.T) {
		called := false
		dispatcher := NewDispatcher(
			WithMiddlewares(MiddlewareForEvent(func(ctx context.Context, event *gitlab.PushEvent) error {
				called = true
				return nil
			})),
		)

		err := dispatcher.Dispatch(t.Context(), &gitlab.PushEvent{})

		assert.NoError(t, err)
		assert.True(t, called)

		called = false
		err = dispatcher.Dispatch(t.Context(), &gitlab.MergeEvent{})

		assert.NoError(t, err)
		assert.False(t, called)
	})

	t.Run("for event can stop dispatch", func(t *testing.T) {
		expectedErr := errors.New("stop push")
		listener := &middlewareTestListener{
			onPush: func(ctx context.Context, event *gitlab.PushEvent) error {
				t.Fatal("listener should not be called")
				return nil
			},
		}

		dispatcher := NewDispatcher(
			RegisterListeners(listener),
			WithMiddlewares(MiddlewareForEvent(func(ctx context.Context, event *gitlab.PushEvent) error {
				return expectedErr
			})),
		)

		err := dispatcher.Dispatch(t.Context(), &gitlab.PushEvent{})

		assert.ErrorIs(t, err, expectedErr)
	})
}

type middlewareTestListener struct {
	onPush func(context.Context, *gitlab.PushEvent) error
}

var _ PushListener = (*middlewareTestListener)(nil)

func (l *middlewareTestListener) OnPush(ctx context.Context, event *gitlab.PushEvent) error {
	return l.onPush(ctx, event)
}

func loadFixture(filePath string) []byte {
	content, err := os.ReadFile(filepath.Join("internal", "testdata", filePath))
	if err != nil {
		log.Fatal(err)
	}

	return content
}
