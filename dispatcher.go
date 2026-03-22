package gitlabwebhook

import (
	"context"
	"errors"
	"sync"

	gitlab "gitlab.com/gitlab-org/api/client-go/v2"
)

var (
	ErrUnsupportedEvent = errors.New("gitlab-webhook: unsupported event type")
	ErrInvalidToken     = errors.New("gitlab-webhook: invalid token")
)

type Dispatcher struct {
	buildListeners                      []BuildListener
	commitCommentListeners              []CommitCommentListener
	deploymentListeners                 []DeploymentListener
	emojiListeners                      []EmojiListener
	featureFlagListeners                []FeatureFlagListener
	groupResourceAccessTokenListeners   []GroupResourceAccessTokenListener
	issueCommentListeners               []IssueCommentListener
	issueListeners                      []IssueListener
	jobListeners                        []JobListener
	memberListeners                     []MemberListener
	milestoneListeners                  []MilestoneListener
	mergeCommentListeners               []MergeCommentListener
	mergeListeners                      []MergeListener
	pipelineListeners                   []PipelineListener
	projectListeners                    []ProjectListener
	projectResourceAccessTokenListeners []ProjectResourceAccessTokenListener
	pushListeners                       []PushListener
	releaseListeners                    []ReleaseListener
	snippetCommentListeners             []SnippetCommentListener
	subGroupListeners                   []SubGroupListener
	tagListeners                        []TagListener
	vulnerabilityListeners              []VulnerabilityListener
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
		for _, register := range dispatcherListenerRegistrations {
			register(d, listener)
		}
	}
}

func (d *Dispatcher) RegisterBuildListener(listeners ...BuildListener) {
	d.buildListeners = append(d.buildListeners, listeners...)
}

func (d *Dispatcher) RegisterCommitCommentListener(listeners ...CommitCommentListener) {
	d.commitCommentListeners = append(d.commitCommentListeners, listeners...)
}

func (d *Dispatcher) RegisterDeploymentListener(listeners ...DeploymentListener) {
	d.deploymentListeners = append(d.deploymentListeners, listeners...)
}

func (d *Dispatcher) RegisterEmojiListener(listeners ...EmojiListener) {
	d.emojiListeners = append(d.emojiListeners, listeners...)
}

func (d *Dispatcher) RegisterFeatureFlagListener(listeners ...FeatureFlagListener) {
	d.featureFlagListeners = append(d.featureFlagListeners, listeners...)
}

func (d *Dispatcher) RegisterGroupResourceAccessTokenListener(listeners ...GroupResourceAccessTokenListener) {
	d.groupResourceAccessTokenListeners = append(d.groupResourceAccessTokenListeners, listeners...)
}

func (d *Dispatcher) RegisterIssueCommentListener(listeners ...IssueCommentListener) {
	d.issueCommentListeners = append(d.issueCommentListeners, listeners...)
}

func (d *Dispatcher) RegisterIssueListener(listeners ...IssueListener) {
	d.issueListeners = append(d.issueListeners, listeners...)
}

func (d *Dispatcher) RegisterJobListener(listeners ...JobListener) {
	d.jobListeners = append(d.jobListeners, listeners...)
}

func (d *Dispatcher) RegisterMemberListener(listeners ...MemberListener) {
	d.memberListeners = append(d.memberListeners, listeners...)
}

func (d *Dispatcher) RegisterMilestoneListener(listeners ...MilestoneListener) {
	d.milestoneListeners = append(d.milestoneListeners, listeners...)
}

func (d *Dispatcher) RegisterMergeCommentListener(listeners ...MergeCommentListener) {
	d.mergeCommentListeners = append(d.mergeCommentListeners, listeners...)
}

func (d *Dispatcher) RegisterMergeListener(listeners ...MergeListener) {
	d.mergeListeners = append(d.mergeListeners, listeners...)
}

func (d *Dispatcher) RegisterPipelineListener(listeners ...PipelineListener) {
	d.pipelineListeners = append(d.pipelineListeners, listeners...)
}

func (d *Dispatcher) RegisterProjectListener(listeners ...ProjectListener) {
	d.projectListeners = append(d.projectListeners, listeners...)
}

func (d *Dispatcher) RegisterProjectResourceAccessTokenListener(listeners ...ProjectResourceAccessTokenListener) {
	d.projectResourceAccessTokenListeners = append(d.projectResourceAccessTokenListeners, listeners...)
}

func (d *Dispatcher) RegisterPushListener(listeners ...PushListener) {
	d.pushListeners = append(d.pushListeners, listeners...)
}

func (d *Dispatcher) RegisterReleaseListener(listeners ...ReleaseListener) {
	d.releaseListeners = append(d.releaseListeners, listeners...)
}

func (d *Dispatcher) RegisterSnippetCommentListener(listeners ...SnippetCommentListener) {
	d.snippetCommentListeners = append(d.snippetCommentListeners, listeners...)
}

func (d *Dispatcher) RegisterSubGroupListener(listeners ...SubGroupListener) {
	d.subGroupListeners = append(d.subGroupListeners, listeners...)
}

func (d *Dispatcher) RegisterTagListener(listeners ...TagListener) {
	d.tagListeners = append(d.tagListeners, listeners...)
}

func (d *Dispatcher) RegisterVulnerabilityListener(listeners ...VulnerabilityListener) {
	d.vulnerabilityListeners = append(d.vulnerabilityListeners, listeners...)
}

func (d *Dispatcher) RegisterWikiPageListener(listeners ...WikiPageListener) {
	d.wikiPageListeners = append(d.wikiPageListeners, listeners...)
}

func (d *Dispatcher) Dispatch(ctx context.Context, event any) error {
	for _, dispatch := range dispatcherEventDispatchers {
		handled, err := dispatch(d, ctx, event)
		if handled {
			return err
		}
	}

	return ErrUnsupportedEvent
}

func (d *Dispatcher) DispatchWebhook(ctx context.Context, eventType gitlab.EventType, payload []byte) error {
	event, err := gitlab.ParseWebhook(eventType, payload)
	if err != nil {
		return err
	}

	return d.Dispatch(ctx, event)
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
