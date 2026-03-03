package components

import (
	"context"

	"github.com/AnatoleLucet/loom"
	"github.com/AnatoleLucet/loom/signals"
)

// mainly used to re-export what's in loom/signals for users to import with loom/components

type Accessor[T any] = signals.Accessor[T]

func Signal[T any](initial T) (Accessor[T], func(T)) {
	return signals.Signal(initial)
}

func Memo[T any](fn func() T) Accessor[T] {
	return signals.Memo(fn)
}

func Effect(effect func()) {
	signals.Effect(effect)
}

func RenderEffect(effect func()) {
	signals.RenderEffect(effect)
}

func Batch(fn func()) {
	signals.Batch(fn)
}

func OnCleanup(fn func()) {
	signals.OnCleanup(fn)
}

func OnSettled(fn func()) {
	signals.OnSettled(fn)
}

func OnUserSettled(fn func()) {
	signals.OnUserSettled(fn)
}

func OnRenderSettled(fn func()) {
	signals.OnRenderSettled(fn)
}

func Untrack[T any](fn func() T) T {
	return signals.Untrack(fn)
}

type Context[T any] = signals.Context[T]

func NewContext[T any](defaultValue T) *Context[T] {
	return signals.NewContext(defaultValue)
}

type Owner = signals.Owner

func NewOwner() *Owner {
	return signals.NewOwner()
}

type component struct {
	ctx context.Context
}

func Self() loom.Component {
	ctx, cancel := context.WithCancel(context.Background())

	OnCleanup(cancel)

	return &component{ctx}
}

func (c *component) Context() context.Context {
	return c.ctx
}

func (c *component) IsDisposed() bool {
	return c.ctx.Err() != nil
}

func (c *component) Disposed() <-chan struct{} {
	return c.ctx.Done()
}
