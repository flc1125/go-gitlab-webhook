package gitlabwebhook

import (
	"context"
	"crypto/subtle"
	"io"
	"net/http"

	gitlab "gitlab.com/gitlab-org/api/client-go/v2"
)

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

func newDispatchRequestOptions(req *http.Request, opts []DispatchRequestOption) *dispatchRequestOptions {
	options := &dispatchRequestOptions{
		ctx: req.Context(),
	}

	for _, opt := range opts {
		opt(options)
	}

	return options
}

func validateRequestToken(req *http.Request, token string) error {
	if token == "" {
		return nil
	}

	requestToken := req.Header.Get("X-Gitlab-Token")
	// Constant-time comparison avoids leaking token length match timing.
	if subtle.ConstantTimeCompare([]byte(requestToken), []byte(token)) != 1 {
		return ErrInvalidToken
	}

	return nil
}

func readRequestPayload(req *http.Request) ([]byte, error) {
	return io.ReadAll(req.Body)
}

func (d *Dispatcher) DispatchRequest(req *http.Request, opts ...DispatchRequestOption) error {
	options := newDispatchRequestOptions(req, opts)

	if err := validateRequestToken(req, options.token); err != nil {
		return err
	}

	payload, err := readRequestPayload(req)
	if err != nil {
		return err
	}

	return d.DispatchWebhook(options.ctx, gitlab.HookEventType(req), payload)
}
