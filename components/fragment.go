package components

import (
	"github.com/AnatoleLucet/loom"
)

// Fragment groups multiple children without adding extra nodes to the output.
// Children inherit the same parent as the Fragment.
func Fragment(children ...loom.Node) loom.Node {
	return &fragmentNode{children}
}

type fragmentNode struct {
	children []loom.Node
}

func (n *fragmentNode) ID() string {
	return "loom.Fragment"
}

func (n *fragmentNode) Mount(slot *loom.Slot) error {
	return slot.RenderChildren(n.children...)
}

func (n *fragmentNode) Update(slot *loom.Slot) error {
	return slot.RenderChildren(n.children...)
}

func (n *fragmentNode) Unmount(slot *loom.Slot) error {
	return nil
}
