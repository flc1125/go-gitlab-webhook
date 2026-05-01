package gitlabwebhook

import "context"

// HandlerFunc handles a parsed GitLab webhook event.
//
// The event argument is one of the concrete event pointer types returned by
// gitlab.ParseWebhook, such as *gitlab.PushEvent or *gitlab.MergeEvent.
type HandlerFunc func(ctx context.Context, event any) error

// Middleware wraps a HandlerFunc to run code before or after event dispatch.
//
// A middleware can stop dispatch by returning an error without calling next.
type Middleware func(next HandlerFunc) HandlerFunc

// MiddlewareForEvent returns middleware that runs fn only for events assignable to E.
//
// If fn returns an error, dispatch stops and that error is returned. Events of
// other types skip fn and continue to the next handler.
func MiddlewareForEvent[E any](fn func(context.Context, E) error) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, event any) error {
			e, ok := event.(E)
			if !ok {
				return next(ctx, event)
			}

			if err := fn(ctx, e); err != nil {
				return err
			}

			return next(ctx, event)
		}
	}
}

// WithMiddlewares registers middleware during dispatcher construction.
//
// Middleware runs once per parsed webhook event, before the event is dispatched
// to registered listeners.
func WithMiddlewares(middlewares ...Middleware) Option {
	return func(d *Dispatcher) {
		d.Use(middlewares...)
	}
}

// Use appends middleware to the dispatcher.
//
// Middleware is applied in registration order: the first middleware wraps the
// second, and so on, with listener dispatch as the final handler.
func (d *Dispatcher) Use(middlewares ...Middleware) {
	d.middlewares = append(d.middlewares, middlewares...)
}
