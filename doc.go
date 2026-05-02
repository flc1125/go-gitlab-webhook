// Package gitlabwebhook dispatches GitLab webhook events to registered listeners.
//
// Create a [Dispatcher] with [NewDispatcher], register listeners, and pass
// incoming HTTP webhook requests to [Dispatcher.DispatchRequest]. Middleware can
// be added with [WithMiddlewares] or [Dispatcher.Use] before dispatching starts.
package gitlabwebhook
