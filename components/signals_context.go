package components

import (
	"github.com/AnatoleLucet/loom"
	"github.com/AnatoleLucet/loom/signals"
)

type Context[T any] struct {
	*signals.Context[T]
}

func NewContext[T any](defaultValue T) (Accessor[T], *Context[T]) {
	get, context := signals.NewContext(defaultValue)

	return get, &Context[T]{context}
}

func (c *Context[T]) Provider(value T, fn func() loom.Node) loom.Node {
	return &providerNode[T]{
		ctx:   c,
		value: func() T { return value },
		fn:    fn,
	}
}

func (c *Context[T]) BindProvider(value Accessor[T], fn func() loom.Node) loom.Node {
	return &providerNode[T]{
		ctx:   c,
		value: value,
		fn:    fn,
	}
}

type providerNode[T any] struct {
	ctx   *Context[T]
	value func() T
	fn    func() loom.Node
}

func (p *providerNode[T]) ID() string {
	return "loom.Provider"
}

func (p *providerNode[T]) Mount(slot *loom.Slot) error {
	return p.Update(slot)
}

func (p *providerNode[T]) Update(slot *loom.Slot) (err error) {
	p.ctx.Context.BindProvider(p.value, func() {
		err = slot.RenderChildren(p.fn())
	})

	return err
}

func (p *providerNode[T]) Unmount(slot *loom.Slot) error {
	return nil
}
