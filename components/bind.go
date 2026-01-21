package components

import (
	"sync/atomic"

	"github.com/AnatoleLucet/loom"
	"github.com/AnatoleLucet/loom/signals"
)

// Bind creates a reactive Node that re-renders whenever any signal
// accessed within the provided function changes.
func Bind(fn func() loom.Node) loom.Node {
	return &bindNode{
		fn:          fn,
		renderOwner: signals.NewOwner(),
	}
}

type bindNode struct {
	fn          func() loom.Node
	renderOwner *signals.Owner // owns the rendered children
}

func (n *bindNode) ID() string {
	return "loom.Bind"
}

func (n *bindNode) Mount(slot *loom.Slot) (err error) {
	var initial atomic.Bool
	initial.Store(true)

	signals.RenderEffect(func() {
		node := n.fn()

		err := n.renderOwner.Run(func() error {
			return slot.RenderChildren(node)
		})

		if err != nil && !initial.Load() {
			panic(err)
		}
		initial.Store(false)
	})

	return err
}

func (n *bindNode) Update(slot *loom.Slot) error {
	return nil
}

func (n *bindNode) Unmount(slot *loom.Slot) error {
	n.renderOwner.Dispose()
	return nil
}
