package components

import (
	"github.com/AnatoleLucet/loom"
	"github.com/AnatoleLucet/loom/signals"
)

func Provider[T any](ctx *signals.Context[T], value T, fn func() loom.Node) loom.Node {
	return &providerNode[T]{
		owner: signals.NewOwner(),
		ctx:   ctx,
		value: value,
		fn:    fn,
	}
}

func ProviderBind[T any](ctx *signals.Context[T], value signals.Accessor[T], fn func() loom.Node) loom.Node {
	owner := signals.NewOwner()

	return Bind(func() loom.Node {
		return &providerNode[T]{
			owner: owner,
			ctx:   ctx,
			value: value(),
			fn:    fn,
		}
	})
}

type providerNode[T any] struct {
	owner *signals.Owner
	ctx   *signals.Context[T]
	value T
	fn    func() loom.Node
}

func (p *providerNode[T]) ID() string {
	return "loom.Provider"
}

func (p *providerNode[T]) Mount(slot *loom.Slot) error {
	return p.Update(slot)
}

func (p *providerNode[T]) Update(slot *loom.Slot) error {
	return p.owner.Run(func() error {
		p.ctx.Set(p.value)
		return slot.RenderChildren(p.fn())
	})
}

func (p *providerNode[T]) Unmount(slot *loom.Slot) error {
	return nil
}
