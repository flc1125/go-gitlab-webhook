package gitlabwebhook

import (
	"context"
	"crypto/subtle"
	"errors"
	"io"
	"net/http"
	"sync"
	"sync/atomic"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

var (
	ErrUnsupportedEvent = errors.New("gitlab-webhook: unsupported event type")
	ErrInvalidToken     = errors.New("gitlab-webhook: invalid token")
)

type eventListeners struct {
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

type Dispatcher struct {
	listeners atomic.Pointer[eventListeners]
	mu        sync.Mutex // serializes writes
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

func (d *Dispatcher) loadListeners() *eventListeners {
	l := d.listeners.Load()
	if l == nil {
		return &eventListeners{}
	}
	return l
}

func (d *Dispatcher) copyListeners() *eventListeners {
	old := d.loadListeners()
	newListeners := *old
	return &newListeners
}

func (d *Dispatcher) RegisterBuildListener(listeners ...BuildListener) {
	d.mu.Lock()
	defer d.mu.Unlock()
	l := d.copyListeners()
	l.buildListeners = append(l.buildListeners, listeners...)
	d.listeners.Store(l)
}

func (d *Dispatcher) RegisterCommitCommentListener(listeners ...CommitCommentListener) {
	d.mu.Lock()
	defer d.mu.Unlock()
	l := d.copyListeners()
	l.commitCommentListeners = append(l.commitCommentListeners, listeners...)
	d.listeners.Store(l)
}

func (d *Dispatcher) RegisterDeploymentListener(listeners ...DeploymentListener) {
	d.mu.Lock()
	defer d.mu.Unlock()
	l := d.copyListeners()
	l.deploymentListeners = append(l.deploymentListeners, listeners...)
	d.listeners.Store(l)
}

func (d *Dispatcher) RegisterFeatureFlagListener(listeners ...FeatureFlagListener) {
	d.mu.Lock()
	defer d.mu.Unlock()
	l := d.copyListeners()
	l.featureFlagListeners = append(l.featureFlagListeners, listeners...)
	d.listeners.Store(l)
}

func (d *Dispatcher) RegisterGroupResourceAccessTokenListener(listeners ...GroupResourceAccessTokenListener) {
	d.mu.Lock()
	defer d.mu.Unlock()
	l := d.copyListeners()
	l.groupResourceAccessTokenListeners = append(l.groupResourceAccessTokenListeners, listeners...)
	d.listeners.Store(l)
}

func (d *Dispatcher) RegisterIssueCommentListener(listeners ...IssueCommentListener) {
	d.mu.Lock()
	defer d.mu.Unlock()
	l := d.copyListeners()
	l.issueCommentListeners = append(l.issueCommentListeners, listeners...)
	d.listeners.Store(l)
}

func (d *Dispatcher) RegisterIssueListener(listeners ...IssueListener) {
	d.mu.Lock()
	defer d.mu.Unlock()
	l := d.copyListeners()
	l.issueListeners = append(l.issueListeners, listeners...)
	d.listeners.Store(l)
}

func (d *Dispatcher) RegisterJobListener(listeners ...JobListener) {
	d.mu.Lock()
	defer d.mu.Unlock()
	l := d.copyListeners()
	l.jobListeners = append(l.jobListeners, listeners...)
	d.listeners.Store(l)
}

func (d *Dispatcher) RegisterMemberListener(listeners ...MemberListener) {
	d.mu.Lock()
	defer d.mu.Unlock()
	l := d.copyListeners()
	l.memberListeners = append(l.memberListeners, listeners...)
	d.listeners.Store(l)
}

func (d *Dispatcher) RegisterMergeCommentListener(listeners ...MergeCommentListener) {
	d.mu.Lock()
	defer d.mu.Unlock()
	l := d.copyListeners()
	l.mergeCommentListeners = append(l.mergeCommentListeners, listeners...)
	d.listeners.Store(l)
}

func (d *Dispatcher) RegisterMergeListener(listeners ...MergeListener) {
	d.mu.Lock()
	defer d.mu.Unlock()
	l := d.copyListeners()
	l.mergeListeners = append(l.mergeListeners, listeners...)
	d.listeners.Store(l)
}

func (d *Dispatcher) RegisterPipelineListener(listeners ...PipelineListener) {
	d.mu.Lock()
	defer d.mu.Unlock()
	l := d.copyListeners()
	l.pipelineListeners = append(l.pipelineListeners, listeners...)
	d.listeners.Store(l)
}

func (d *Dispatcher) RegisterProjectResourceAccessTokenListener(listeners ...ProjectResourceAccessTokenListener) {
	d.mu.Lock()
	defer d.mu.Unlock()
	l := d.copyListeners()
	l.projectResourceAccessTokenListeners = append(l.projectResourceAccessTokenListeners, listeners...)
	d.listeners.Store(l)
}

func (d *Dispatcher) RegisterPushListener(listeners ...PushListener) {
	d.mu.Lock()
	defer d.mu.Unlock()
	l := d.copyListeners()
	l.pushListeners = append(l.pushListeners, listeners...)
	d.listeners.Store(l)
}

func (d *Dispatcher) RegisterReleaseListener(listeners ...ReleaseListener) {
	d.mu.Lock()
	defer d.mu.Unlock()
	l := d.copyListeners()
	l.releaseListeners = append(l.releaseListeners, listeners...)
	d.listeners.Store(l)
}

func (d *Dispatcher) RegisterSnippetCommentListener(listeners ...SnippetCommentListener) {
	d.mu.Lock()
	defer d.mu.Unlock()
	l := d.copyListeners()
	l.snippetCommentListeners = append(l.snippetCommentListeners, listeners...)
	d.listeners.Store(l)
}

func (d *Dispatcher) RegisterSubGroupListener(listeners ...SubGroupListener) {
	d.mu.Lock()
	defer d.mu.Unlock()
	l := d.copyListeners()
	l.subGroupListeners = append(l.subGroupListeners, listeners...)
	d.listeners.Store(l)
}

func (d *Dispatcher) RegisterTagListener(listeners ...TagListener) {
	d.mu.Lock()
	defer d.mu.Unlock()
	l := d.copyListeners()
	l.tagListeners = append(l.tagListeners, listeners...)
	d.listeners.Store(l)
}

func (d *Dispatcher) RegisterWikiPageListener(listeners ...WikiPageListener) {
	d.mu.Lock()
	defer d.mu.Unlock()
	l := d.copyListeners()
	l.wikiPageListeners = append(l.wikiPageListeners, listeners...)
	d.listeners.Store(l)
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
	return processEvent(ctx, d.loadListeners().buildListeners, BuildListener.OnBuild, event)
}

func (d *Dispatcher) processCommitCommentEvent(ctx context.Context, event *gitlab.CommitCommentEvent) error {
	return processEvent(ctx, d.loadListeners().commitCommentListeners, CommitCommentListener.OnCommitComment, event)
}

func (d *Dispatcher) processDeploymentEvent(ctx context.Context, event *gitlab.DeploymentEvent) error {
	return processEvent(ctx, d.loadListeners().deploymentListeners, DeploymentListener.OnDeployment, event)
}

func (d *Dispatcher) processFeatureFlagEvent(ctx context.Context, event *gitlab.FeatureFlagEvent) error {
	return processEvent(ctx, d.loadListeners().featureFlagListeners, FeatureFlagListener.OnFeatureFlag, event)
}

func (d *Dispatcher) processGroupResourceAccessTokenEvent(ctx context.Context, event *gitlab.GroupResourceAccessTokenEvent) error { //nolint:lll
	return processEvent(ctx, d.loadListeners().groupResourceAccessTokenListeners, GroupResourceAccessTokenListener.OnGroupResourceAccessToken, event)
}

func (d *Dispatcher) processIssueCommentEvent(ctx context.Context, event *gitlab.IssueCommentEvent) error {
	return processEvent(ctx, d.loadListeners().issueCommentListeners, IssueCommentListener.OnIssueComment, event)
}

func (d *Dispatcher) processIssueEvent(ctx context.Context, event *gitlab.IssueEvent) error {
	return processEvent(ctx, d.loadListeners().issueListeners, IssueListener.OnIssue, event)
}

func (d *Dispatcher) processJobEvent(ctx context.Context, event *gitlab.JobEvent) error {
	return processEvent(ctx, d.loadListeners().jobListeners, JobListener.OnJob, event)
}

func (d *Dispatcher) processMemberEvent(ctx context.Context, event *gitlab.MemberEvent) error {
	return processEvent(ctx, d.loadListeners().memberListeners, MemberListener.OnMember, event)
}

func (d *Dispatcher) processMergeCommentEvent(ctx context.Context, event *gitlab.MergeCommentEvent) error {
	return processEvent(ctx, d.loadListeners().mergeCommentListeners, MergeCommentListener.OnMergeComment, event)
}

func (d *Dispatcher) processMergeEvent(ctx context.Context, event *gitlab.MergeEvent) error {
	return processEvent(ctx, d.loadListeners().mergeListeners, MergeListener.OnMerge, event)
}

func (d *Dispatcher) processPipelineEvent(ctx context.Context, event *gitlab.PipelineEvent) error {
	return processEvent(ctx, d.loadListeners().pipelineListeners, PipelineListener.OnPipeline, event)
}

func (d *Dispatcher) processProjectResourceAccessTokenEvent(ctx context.Context, event *gitlab.ProjectResourceAccessTokenEvent) error { //nolint:lll
	return processEvent(ctx, d.loadListeners().projectResourceAccessTokenListeners, ProjectResourceAccessTokenListener.OnProjectResourceAccessToken, event)
}

func (d *Dispatcher) processPushEvent(ctx context.Context, event *gitlab.PushEvent) error {
	return processEvent(ctx, d.loadListeners().pushListeners, PushListener.OnPush, event)
}

func (d *Dispatcher) processReleaseEvent(ctx context.Context, event *gitlab.ReleaseEvent) error {
	return processEvent(ctx, d.loadListeners().releaseListeners, ReleaseListener.OnRelease, event)
}

func (d *Dispatcher) processSnippetCommentEvent(ctx context.Context, event *gitlab.SnippetCommentEvent) error {
	return processEvent(ctx, d.loadListeners().snippetCommentListeners, SnippetCommentListener.OnSnippetComment, event)
}

func (d *Dispatcher) processSubGroupEvent(ctx context.Context, event *gitlab.SubGroupEvent) error {
	return processEvent(ctx, d.loadListeners().subGroupListeners, SubGroupListener.OnSubGroup, event)
}

func (d *Dispatcher) processTagEvent(ctx context.Context, event *gitlab.TagEvent) error {
	return processEvent(ctx, d.loadListeners().tagListeners, TagListener.OnTag, event)
}

func (d *Dispatcher) processWikiPageEvent(ctx context.Context, event *gitlab.WikiPageEvent) error {
	return processEvent(ctx, d.loadListeners().wikiPageListeners, WikiPageListener.OnWikiPage, event)
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
