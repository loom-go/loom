package loom

import "context"

// Component represents the current component instance, providing access to its context and lifecycle.
// It can be used to sync goroutines with the component's lifecycle,
// ensuring they are properly cleaned up when the component is unmounted.
type Component interface {
	// Context returns a context that is canceled when the component is unmounted.
	Context() context.Context

	// IsDisposed returns true if the component has been unmounted and its context has been canceled.
	IsDisposed() (disposed bool)

	// Disposed returns a channel that is closed when the component is unmounted.
	Disposed() (done <-chan struct{})
}
