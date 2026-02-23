package components

import (
	"fmt"

	"github.com/AnatoleLucet/loom"
)

// Ref creates a Node that calls the provided function with the parent node as an argument.
func Ref[T any](fn func(T)) loom.Node {
	return &refNode[T]{fn}
}

type refNode[T any] struct {
	fn func(T)
}

func (n *refNode[T]) ID() string {
	return "loom.Ref"
}

func (n *refNode[T]) Mount(slot *loom.Slot) error {
	parent := slot.Parent()

	casted, ok := parent.(T)
	if !ok {
		return fmt.Errorf("Ref: %w: the given ref type did not match the node type", ErrNodeRefMissMatch)
	}

	n.fn(casted)
	return nil
}

func (n *refNode[T]) Update(slot *loom.Slot) error {
	return nil
}

func (n *refNode[T]) Unmount(slot *loom.Slot) error {
	return nil
}
