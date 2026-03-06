package signals

import "github.com/AnatoleLucet/sig"

type Context[T any] struct {
	ctx *sig.Context[Accessor[T]]
}

func NewContext[T any](defaultValue T) (Accessor[T], *Context[T]) {
	ctx := sig.NewContext[Accessor[T]](func() T {
		return defaultValue
	})

	return func() T { return ctx.Value()() }, &Context[T]{ctx}
}

func (c *Context[T]) Get() T {
	return Untrack(c.ctx.Value())
}

func (c *Context[T]) Set(value T) {
	c.ctx.Set(func() T { return value })
}

func (c *Context[T]) Provider(value T, fn func()) {
	c.BindProvider(func() T { return value }, fn)
}

func (c *Context[T]) BindProvider(value Accessor[T], fn func()) {
	NewOwner().Run(func() error {
		c.ctx.Set(value)
		fn()
		return nil
	})
}
