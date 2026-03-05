package components

import (
	"fmt"

	"github.com/AnatoleLucet/loom"
)

// Ref creates a Node that calls the provided function with the parent node as an argument.
func Ref[T any](ref *T) loom.Node {
	return &refNode[T]{ref}
}

type refNode[T any] struct {
	ref *T
}

func (n *refNode[T]) ID() string {
	return "loom.Ref"
}

func (n *refNode[T]) Mount(slot *loom.Slot) error {
	parent := slot.Parent()

	casted, ok := parent.(T)
	if !ok {
		return fmt.Errorf("Ref: %w: the given ref type (%T) does not match the parent node type (%T)", ErrNodeRefMissMatch, *new(T), parent)
	}

	*n.ref = casted
	return nil
}

func (n *refNode[T]) Update(slot *loom.Slot) error {
	return nil
}

func (n *refNode[T]) Unmount(slot *loom.Slot) error {
	return nil
}

// OnRef creates a Node that calls the provided function with the parent node as an argument.
func OnRef[T any](fn func(T)) loom.Node {
	return &onRefNode[T]{fn}
}

type onRefNode[T any] struct {
	fn func(T)
}

func (n *onRefNode[T]) ID() string {
	return "loom.OnRef"
}

func (n *onRefNode[T]) Mount(slot *loom.Slot) error {
	parent := slot.Parent()

	casted, ok := parent.(T)
	if !ok {
		return fmt.Errorf("OnRef: %w: the given ref type (%T) does not match the parent node type (%T)", ErrNodeRefMissMatch, *new(T), parent)
	}

	n.fn(casted)
	return nil
}

func (n *onRefNode[T]) Update(slot *loom.Slot) error {
	return nil
}

func (n *onRefNode[T]) Unmount(slot *loom.Slot) error {
	return nil
}
