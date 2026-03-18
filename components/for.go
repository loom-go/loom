package components

import (
	"sync/atomic"

	"github.com/AnatoleLucet/sig"
	"github.com/loom-go/loom"
)

func For[T comparable](
	items Accessor[[]T],
	mapper func(T, Accessor[int]) loom.Node,
) loom.Node {
	return &forNode[T]{
		input:       items,
		mapper:      mapper,
		renderOwner: NewOwner(),
	}
}

type forNode[T comparable] struct {
	input  Accessor[[]T]
	mapper func(T, Accessor[int]) loom.Node

	// currently rendered items
	items []*forItem[T]
	// owns the rendered children
	renderOwner *Owner
}

func (n *forNode[T]) ID() string {
	return "loom.For"
}

func (n *forNode[T]) Mount(slot *loom.Slot) (err error) {
	var initial atomic.Bool
	initial.Store(true)

	RenderEffect(func() {
		newItems := n.input()

		err = n.renderOwner.Run(func() error {
			return n.reconcile(slot, newItems)
		})

		if err != nil && !initial.Load() {
			panic(err)
		}
		initial.Store(false)
	})

	return err
}

func (n *forNode[T]) Update(slot *loom.Slot) error {
	return nil
}

func (n *forNode[T]) Unmount(slot *loom.Slot) error {
	n.disposeAll()
	return nil
}

func (n *forNode[T]) reconcile(slot *loom.Slot, newItems []T) error {
	oldLen := len(n.items)
	newLen := len(newItems)

	// fast path for empty list
	if newLen == 0 {
		n.disposeAll()
		return slot.UnmountChildren()
	}

	// fast path for create
	if oldLen == 0 {
		n.initItems(newItems)
		for i, item := range n.items {
			if err := slot.RenderChild(i, item); err != nil {
				return err
			}
		}
		return nil
	}

	result := make([]*forItem[T], newLen)

	// common prefix
	start := 0
	for start < oldLen && start < newLen && n.items[start].value == newItems[start] {
		result[start] = n.items[start]
		start++
	}

	// common suffix
	oldEnd := oldLen - 1
	newEnd := newLen - 1
	for oldEnd >= start && newEnd >= start && n.items[oldEnd].value == newItems[newEnd] {
		result[newEnd] = n.items[oldEnd]
		oldEnd--
		newEnd--
	}

	// index map for new window [start...newEnd]
	indices := make(map[T]int, newEnd-start+1)
	for i := start; i <= newEnd; i++ {
		indices[newItems[i]] = i
	}

	moved := make([]*forItem[T], newLen)

	// walk old window [start...oldEnd]
	// collect in moved if item is in new window, dispose if not
	for i := start; i <= oldEnd; i++ {
		item := n.items[i]
		if j, ok := indices[item.value]; ok {
			moved[j] = item
			delete(indices, item.value)
		} else {
			item.owner.Dispose()
		}
	}

	// unmount all items in window [start...oldEnd]
	// potentially moved, or removed
	for i := oldLen - 1; i >= start; i-- {
		if err := slot.UnmountChild(i); err != nil {
			return err
		}
	}

	// render all items at their new positions
	for i := start; i < newLen; i++ {
		if i <= newEnd {
			// create or reuse from moved
			if moved[i] != nil {
				result[i] = moved[i]
				result[i].index.Write(i)
			} else {
				result[i] = n.initItem(i, newItems[i])
			}
		}

		if err := slot.RenderChild(i, result[i]); err != nil {
			return err
		}
	}

	n.items = result
	return nil
}

func (n *forNode[T]) initItems(items []T) {
	for i, item := range items {
		n.items = append(n.items, n.initItem(i, item))
	}
}

func (n *forNode[T]) initItem(index int, item T) *forItem[T] {
	indexSignal := sig.NewSignal(index)

	var node loom.Node
	owner := NewOwner()
	owner.Run(func() error {
		node = n.mapper(item, indexSignal.Read)
		return nil
	})

	return &forItem[T]{
		value: item,
		node:  node,
		owner: owner,
		index: indexSignal,
	}
}

func (n *forNode[T]) disposeAll() {
	n.renderOwner.Dispose()
	n.items = nil
}

// used to make sure children are owned by the owner
// and the reactive scope can stay active when items are moved in the list
type forItem[T comparable] struct {
	value T
	node  loom.Node
	owner *Owner
	index *sig.Signal[int]
}

func (n *forItem[T]) ID() string {
	return "loom.forItem"
}

func (n *forItem[T]) Mount(slot *loom.Slot) error {
	return n.Update(slot)
}

func (n *forItem[T]) Update(slot *loom.Slot) error {
	return n.owner.Run(func() error {
		return slot.RenderChildren(n.node)
	})
}

func (n *forItem[T]) Unmount(slot *loom.Slot) error {
	return nil
}
