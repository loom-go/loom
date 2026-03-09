package components

import (
	"github.com/loom-go/loom"
)

// Own wraps a node so that it is managed by the given signals.Owner.
// Note that disposing the owner will NOT unmount the node;
// Own is only responsible assigning ownership of signals used within the node.
func Own(owner *Owner, node loom.Node) loom.Node {
	return &ownNode{owner, node}
}

type ownNode struct {
	owner *Owner
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
