---
title: Self()
weight: 1
---

```go {style=tokyonight-moon}
func Self() Component
```

```go {style=tokyonight-moon}
type Component interface {
	// Context returns a context that is canceled when the component is unmounted.
	Context() context.Context

	// IsDisposed returns true if the component has been unmounted and its context has been canceled.
	IsDisposed() (disposed bool)

	// Disposed returns a channel that is closed when the component is unmounted.
	Disposed() (done <-chan struct{})
}
```
