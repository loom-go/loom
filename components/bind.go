package components

import (
	"sync/atomic"

	"github.com/AnatoleLucet/loom"
)

// Bind creates a reactive Node that re-renders whenever any signal
// accessed within the provided function changes.
func Bind(fn func() loom.Node) loom.Node {
	return &bindNode{
		fn: fn,
	}
}

type bindNode struct {
	fn func() loom.Node
}

func (n *bindNode) ID() string {
	return "loom.Bind"
}

func (n *bindNode) Mount(slot *loom.Slot) (err error) {
	var initial atomic.Bool
	initial.Store(true)

	RenderEffect(func() {
		err = slot.RenderChildren(n.fn())

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
	return nil
}
