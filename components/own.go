package components

import (
	"github.com/AnatoleLucet/loom"
	"github.com/AnatoleLucet/loom/signals"
)

// Own wraps a node so that it is managed by the given signals.Owner.
func Own(owner *signals.Owner, node loom.Node) loom.Node {
	return &ownNode{owner, node}
}

type ownNode struct {
	owner *signals.Owner
	node  loom.Node
}

func (n *ownNode) ID() string {
	return "loom.Own"
}

func (n *ownNode) Mount(slot *loom.Slot) error {
	return n.Update(slot)
}

func (n *ownNode) Update(slot *loom.Slot) error {
	return n.owner.Run(func() error {
		return slot.RenderChildren(n.node)
	})
}

func (n *ownNode) Unmount(slot *loom.Slot) error {
	n.owner.Dispose()
	return nil
}
