package gitlabwebhook

import (
	"context"
	"crypto/subtle"
	"errors"
	"io"
	"net/http"
	"sync"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

var (
	ErrUnsupportedEvent = errors.New("gitlab-webhook: unsupported event type")
	ErrInvalidToken     = errors.New("gitlab-webhook: invalid token")
)

type Dispatcher struct {
	mu                                  sync.RWMutex
	buildListeners                      []BuildListener
	commitCommentListeners              []CommitCommentListener
	deploymentListeners                 []DeploymentListener
	featureFlagListeners                []FeatureFlagListener
	groupResourceAccessTokenListeners   []GroupResourceAccessTokenListener
	issueCommentListeners               []IssueCommentListener
	issueListeners                      []IssueListener
	jobListeners                        []JobListener
	memberListeners                     []MemberListener
	mergeCommentListeners               []MergeCommentListener
	mergeListeners                      []MergeListener
	pipelineListeners                   []PipelineListener
	projectResourceAccessTokenListeners []ProjectResourceAccessTokenListener
	pushListeners                       []PushListener
	releaseListeners                    []ReleaseListener
	snippetCommentListeners             []SnippetCommentListener
	subGroupListeners                   []SubGroupListener
	tagListeners                        []TagListener
	wikiPageListeners                   []WikiPageListener
}

type Option func(*Dispatcher)

func RegisterListeners(listeners ...any) Option {
	return func(d *Dispatcher) {
		d.RegisterListeners(listeners...)
	}
}

func NewDispatcher(opts ...Option) *Dispatcher {
	dispatcher := &Dispatcher{}
	for _, opt := range opts {
		opt(dispatcher)
	}
	return dispatcher
}

func (d *Dispatcher) RegisterListeners(listeners ...any) {
	for _, listener := range listeners {
		if l, ok := listener.(BuildListener); ok {
			d.RegisterBuildListener(l)
		}

		if l, ok := listener.(CommitCommentListener); ok {
			d.RegisterCommitCommentListener(l)
		}

		if l, ok := listener.(DeploymentListener); ok {
			d.RegisterDeploymentListener(l)
		}

		if l, ok := listener.(FeatureFlagListener); ok {
			d.RegisterFeatureFlagListener(l)
		}

		if l, ok := listener.(GroupResourceAccessTokenListener); ok {
			d.RegisterGroupResourceAccessTokenListener(l)
		}

		if l, ok := listener.(IssueCommentListener); ok {
			d.RegisterIssueCommentListener(l)
		}

		if l, ok := listener.(IssueListener); ok {
			d.RegisterIssueListener(l)
		}

		if l, ok := listener.(JobListener); ok {
			d.RegisterJobListener(l)
		}

		if l, ok := listener.(MemberListener); ok {
			d.RegisterMemberListener(l)
		}

		if l, ok := listener.(MergeCommentListener); ok {
			d.RegisterMergeCommentListener(l)
		}

		if l, ok := listener.(MergeListener); ok {
			d.RegisterMergeListener(l)
		}

		if l, ok := listener.(PipelineListener); ok {
			d.RegisterPipelineListener(l)
		}

		if l, ok := listener.(ProjectResourceAccessTokenListener); ok {
			d.RegisterProjectResourceAccessTokenListener(l)
		}

		if l, ok := listener.(PushListener); ok {
			d.RegisterPushListener(l)
		}

		if l, ok := listener.(ReleaseListener); ok {
			d.RegisterReleaseListener(l)
		}

		if l, ok := listener.(SnippetCommentListener); ok {
			d.RegisterSnippetCommentListener(l)
		}

		if l, ok := listener.(SubGroupListener); ok {
			d.RegisterSubGroupListener(l)
		}

		if l, ok := listener.(TagListener); ok {
			d.RegisterTagListener(l)
		}

		if l, ok := listener.(WikiPageListener); ok {
			d.RegisterWikiPageListener(l)
		}
	}
}

func (d *Dispatcher) RegisterBuildListener(listeners ...BuildListener) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.buildListeners = append(d.buildListeners, listeners...)
}

func (d *Dispatcher) RegisterCommitCommentListener(listeners ...CommitCommentListener) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.commitCommentListeners = append(d.commitCommentListeners, listeners...)
}

func (d *Dispatcher) RegisterDeploymentListener(listeners ...DeploymentListener) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.deploymentListeners = append(d.deploymentListeners, listeners...)
}

func (d *Dispatcher) RegisterFeatureFlagListener(listeners ...FeatureFlagListener) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.featureFlagListeners = append(d.featureFlagListeners, listeners...)
}

func (d *Dispatcher) RegisterGroupResourceAccessTokenListener(listeners ...GroupResourceAccessTokenListener) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.groupResourceAccessTokenListeners = append(d.groupResourceAccessTokenListeners, listeners...)
}

func (d *Dispatcher) RegisterIssueCommentListener(listeners ...IssueCommentListener) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.issueCommentListeners = append(d.issueCommentListeners, listeners...)
}

func (d *Dispatcher) RegisterIssueListener(listeners ...IssueListener) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.issueListeners = append(d.issueListeners, listeners...)
}

func (d *Dispatcher) RegisterJobListener(listeners ...JobListener) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.jobListeners = append(d.jobListeners, listeners...)
}

func (d *Dispatcher) RegisterMemberListener(listeners ...MemberListener) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.memberListeners = append(d.memberListeners, listeners...)
}

func (d *Dispatcher) RegisterMergeCommentListener(listeners ...MergeCommentListener) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.mergeCommentListeners = append(d.mergeCommentListeners, listeners...)
}

func (d *Dispatcher) RegisterMergeListener(listeners ...MergeListener) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.mergeListeners = append(d.mergeListeners, listeners...)
}

func (d *Dispatcher) RegisterPipelineListener(listeners ...PipelineListener) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.pipelineListeners = append(d.pipelineListeners, listeners...)
}

func (d *Dispatcher) RegisterProjectResourceAccessTokenListener(listeners ...ProjectResourceAccessTokenListener) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.projectResourceAccessTokenListeners = append(d.projectResourceAccessTokenListeners, listeners...)
}

func (d *Dispatcher) RegisterPushListener(listeners ...PushListener) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.pushListeners = append(d.pushListeners, listeners...)
}

func (d *Dispatcher) RegisterReleaseListener(listeners ...ReleaseListener) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.releaseListeners = append(d.releaseListeners, listeners...)
}

func (d *Dispatcher) RegisterSnippetCommentListener(listeners ...SnippetCommentListener) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.snippetCommentListeners = append(d.snippetCommentListeners, listeners...)
}

func (d *Dispatcher) RegisterSubGroupListener(listeners ...SubGroupListener) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.subGroupListeners = append(d.subGroupListeners, listeners...)
}

func (d *Dispatcher) RegisterTagListener(listeners ...TagListener) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.tagListeners = append(d.tagListeners, listeners...)
}

func (d *Dispatcher) RegisterWikiPageListener(listeners ...WikiPageListener) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.wikiPageListeners = append(d.wikiPageListeners, listeners...)
}

func (d *Dispatcher) Dispatch(ctx context.Context, event any) error {
	switch e := event.(type) {
	case *gitlab.BuildEvent:
		return d.processBuildEvent(ctx, e)
	case *gitlab.CommitCommentEvent:
		return d.processCommitCommentEvent(ctx, e)
	case *gitlab.DeploymentEvent:
		return d.processDeploymentEvent(ctx, e)
	case *gitlab.FeatureFlagEvent:
		return d.processFeatureFlagEvent(ctx, e)
	case *gitlab.GroupResourceAccessTokenEvent:
		return d.processGroupResourceAccessTokenEvent(ctx, e)
	case *gitlab.IssueCommentEvent:
		return d.processIssueCommentEvent(ctx, e)
	case *gitlab.IssueEvent:
		return d.processIssueEvent(ctx, e)
	case *gitlab.JobEvent:
		return d.processJobEvent(ctx, e)
	case *gitlab.MemberEvent:
		return d.processMemberEvent(ctx, e)
	case *gitlab.MergeCommentEvent:
		return d.processMergeCommentEvent(ctx, e)
	case *gitlab.MergeEvent:
		return d.processMergeEvent(ctx, e)
	case *gitlab.PipelineEvent:
		return d.processPipelineEvent(ctx, e)
	case *gitlab.ProjectResourceAccessTokenEvent:
		return d.processProjectResourceAccessTokenEvent(ctx, e)
	case *gitlab.PushEvent:
		return d.processPushEvent(ctx, e)
	case *gitlab.ReleaseEvent:
		return d.processReleaseEvent(ctx, e)
	case *gitlab.SnippetCommentEvent:
		return d.processSnippetCommentEvent(ctx, e)
	case *gitlab.SubGroupEvent:
		return d.processSubGroupEvent(ctx, e)
	case *gitlab.TagEvent:
		return d.processTagEvent(ctx, e)
	case *gitlab.WikiPageEvent:
		return d.processWikiPageEvent(ctx, e)
	default:
		return ErrUnsupportedEvent
	}
}

func (d *Dispatcher) DispatchWebhook(ctx context.Context, eventType gitlab.EventType, payload []byte) error {
	event, err := gitlab.ParseWebhook(eventType, payload)
	if err != nil {
		return err
	}
	return d.Dispatch(ctx, event)
}

type dispatchRequestOptions struct {
	ctx   context.Context
	token string
}

type DispatchRequestOption func(*dispatchRequestOptions)

func DispatchRequestWithContext(ctx context.Context) DispatchRequestOption {
	return func(o *dispatchRequestOptions) {
		o.ctx = ctx
	}
}

func DispatchRequestWithToken(token string) DispatchRequestOption {
	return func(o *dispatchRequestOptions) {
		o.token = token
	}
}

func (d *Dispatcher) DispatchRequest(req *http.Request, opts ...DispatchRequestOption) error {
	o := &dispatchRequestOptions{
		ctx: req.Context(),
	}
	for _, opt := range opts {
		opt(o)
	}

	// check token if provided
	if o.token != "" {
		token := req.Header.Get("X-Gitlab-Token")
		// constant time compare to prevent timing attacks on token comparison
		if subtle.ConstantTimeCompare([]byte(token), []byte(o.token)) != 1 {
			return ErrInvalidToken
		}
	}

	// read payload
	payload, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}

	// dispatch webhook
	return d.DispatchWebhook(o.ctx, gitlab.HookEventType(req), payload)
}

func (d *Dispatcher) processBuildEvent(ctx context.Context, event *gitlab.BuildEvent) error {
	d.mu.RLock()
	listeners := d.buildListeners
	d.mu.RUnlock()
	return processEvent(ctx, listeners, BuildListener.OnBuild, event)
}

func (d *Dispatcher) processCommitCommentEvent(ctx context.Context, event *gitlab.CommitCommentEvent) error {
	d.mu.RLock()
	listeners := d.commitCommentListeners
	d.mu.RUnlock()
	return processEvent(ctx, listeners, CommitCommentListener.OnCommitComment, event)
}

func (d *Dispatcher) processDeploymentEvent(ctx context.Context, event *gitlab.DeploymentEvent) error {
	d.mu.RLock()
	listeners := d.deploymentListeners
	d.mu.RUnlock()
	return processEvent(ctx, listeners, DeploymentListener.OnDeployment, event)
}

func (d *Dispatcher) processFeatureFlagEvent(ctx context.Context, event *gitlab.FeatureFlagEvent) error {
	d.mu.RLock()
	listeners := d.featureFlagListeners
	d.mu.RUnlock()
	return processEvent(ctx, listeners, FeatureFlagListener.OnFeatureFlag, event)
}

func (d *Dispatcher) processGroupResourceAccessTokenEvent(ctx context.Context, event *gitlab.GroupResourceAccessTokenEvent) error { //nolint:lll
	d.mu.RLock()
	listeners := d.groupResourceAccessTokenListeners
	d.mu.RUnlock()
	return processEvent(ctx, listeners, GroupResourceAccessTokenListener.OnGroupResourceAccessToken, event)
}

func (d *Dispatcher) processIssueCommentEvent(ctx context.Context, event *gitlab.IssueCommentEvent) error {
	d.mu.RLock()
	listeners := d.issueCommentListeners
	d.mu.RUnlock()
	return processEvent(ctx, listeners, IssueCommentListener.OnIssueComment, event)
}

func (d *Dispatcher) processIssueEvent(ctx context.Context, event *gitlab.IssueEvent) error {
	d.mu.RLock()
	listeners := d.issueListeners
	d.mu.RUnlock()
	return processEvent(ctx, listeners, IssueListener.OnIssue, event)
}

func (d *Dispatcher) processJobEvent(ctx context.Context, event *gitlab.JobEvent) error {
	d.mu.RLock()
	listeners := d.jobListeners
	d.mu.RUnlock()
	return processEvent(ctx, listeners, JobListener.OnJob, event)
}

func (d *Dispatcher) processMemberEvent(ctx context.Context, event *gitlab.MemberEvent) error {
	d.mu.RLock()
	listeners := d.memberListeners
	d.mu.RUnlock()
	return processEvent(ctx, listeners, MemberListener.OnMember, event)
}

func (d *Dispatcher) processMergeCommentEvent(ctx context.Context, event *gitlab.MergeCommentEvent) error {
	d.mu.RLock()
	listeners := d.mergeCommentListeners
	d.mu.RUnlock()
	return processEvent(ctx, listeners, MergeCommentListener.OnMergeComment, event)
}

func (d *Dispatcher) processMergeEvent(ctx context.Context, event *gitlab.MergeEvent) error {
	d.mu.RLock()
	listeners := d.mergeListeners
	d.mu.RUnlock()
	return processEvent(ctx, listeners, MergeListener.OnMerge, event)
}

func (d *Dispatcher) processPipelineEvent(ctx context.Context, event *gitlab.PipelineEvent) error {
	d.mu.RLock()
	listeners := d.pipelineListeners
	d.mu.RUnlock()
	return processEvent(ctx, listeners, PipelineListener.OnPipeline, event)
}

func (d *Dispatcher) processProjectResourceAccessTokenEvent(ctx context.Context, event *gitlab.ProjectResourceAccessTokenEvent) error { //nolint:lll
	d.mu.RLock()
	listeners := d.projectResourceAccessTokenListeners
	d.mu.RUnlock()
	return processEvent(ctx, listeners, ProjectResourceAccessTokenListener.OnProjectResourceAccessToken, event)
}

func (d *Dispatcher) processPushEvent(ctx context.Context, event *gitlab.PushEvent) error {
	d.mu.RLock()
	listeners := d.pushListeners
	d.mu.RUnlock()
	return processEvent(ctx, listeners, PushListener.OnPush, event)
}

func (d *Dispatcher) processReleaseEvent(ctx context.Context, event *gitlab.ReleaseEvent) error {
	d.mu.RLock()
	listeners := d.releaseListeners
	d.mu.RUnlock()
	return processEvent(ctx, listeners, ReleaseListener.OnRelease, event)
}

func (d *Dispatcher) processSnippetCommentEvent(ctx context.Context, event *gitlab.SnippetCommentEvent) error {
	d.mu.RLock()
	listeners := d.snippetCommentListeners
	d.mu.RUnlock()
	return processEvent(ctx, listeners, SnippetCommentListener.OnSnippetComment, event)
}

func (d *Dispatcher) processSubGroupEvent(ctx context.Context, event *gitlab.SubGroupEvent) error {
	d.mu.RLock()
	listeners := d.subGroupListeners
	d.mu.RUnlock()
	return processEvent(ctx, listeners, SubGroupListener.OnSubGroup, event)
}

func (d *Dispatcher) processTagEvent(ctx context.Context, event *gitlab.TagEvent) error {
	d.mu.RLock()
	listeners := d.tagListeners
	d.mu.RUnlock()
	return processEvent(ctx, listeners, TagListener.OnTag, event)
}

func (d *Dispatcher) processWikiPageEvent(ctx context.Context, event *gitlab.WikiPageEvent) error {
	d.mu.RLock()
	listeners := d.wikiPageListeners
	d.mu.RUnlock()
	return processEvent(ctx, listeners, WikiPageListener.OnWikiPage, event)
}

func processEvent[E any, L any](ctx context.Context, listeners []L, handler func(L, context.Context, E) error, event E) error {
	if len(listeners) == 0 {
		return nil
	}

	wg := sync.WaitGroup{}
	wg.Add(len(listeners))

	errCh := make(chan error, len(listeners))
	for _, listener := range listeners {
		go func() {
			defer wg.Done()
			if err := handler(listener, ctx, event); err != nil {
				errCh <- err
			}
		}()
	}
	wg.Wait()

	close(errCh)
	var err error
	for e := range errCh {
		err = errors.Join(err, e)
	}
	return err
}
