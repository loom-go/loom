package components

import (
	"reflect"
	"sync/atomic"

	"github.com/AnatoleLucet/sig"
	"github.com/loom-go/loom"
)

func Keyed[T, K any](
	items Accessor[[]T],
	keyer func(T) K,
	mapper func(Accessor[T], Accessor[int]) loom.Node,
) loom.Node {
	return &keyedNode[T, K]{
		input:  items,
		keyer:  keyer,
		mapper: mapper,

		renderOwner: NewOwner(),
	}
}

type keyedNode[T, K any] struct {
	input  Accessor[[]T]
	keyer  func(T) K
	mapper func(Accessor[T], Accessor[int]) loom.Node

	// currently rendered items
	mapped []loom.Node
	// owners of the rendered items
	owners []*Owner
	// signals for each item and index
	items   []*sig.Signal[T]
	indices []*sig.Signal[int]

	// owns the rendered children
	renderOwner *Owner
}

func (n *keyedNode[T, K]) ID() string {
	return "loom.Keyed"
}

func (n *keyedNode[T, K]) Mount(slot *loom.Slot) (err error) {
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

func (n *keyedNode[T, K]) Update(slot *loom.Slot) error {
	return nil
}

func (n *keyedNode[T, K]) Unmount(slot *loom.Slot) error {
	n.disposeAll()
	return nil
}

func (n *keyedNode[T, K]) initItems(items []T) {
	n.mapped = make([]loom.Node, len(items))
	n.owners = make([]*Owner, len(items))
	n.items = make([]*sig.Signal[T], len(items))
	n.indices = make([]*sig.Signal[int], len(items))

	for i, item := range items {
		n.initItem(i, &item)
	}
}

func (n *keyedNode[T, K]) initItem(index int, item *T) {
	itemSignal := sig.NewSignal(*item)
	indexSignal := sig.NewSignal(index)

	var node loom.Node
	owner := NewOwner()
	owner.Run(func() error {
		node = n.mapper(itemSignal.Read, indexSignal.Read)
		return nil
	})

	n.mapped[index] = Own(owner, node)
	n.owners[index] = owner
	n.items[index] = itemSignal
	n.indices[index] = indexSignal
}

func (n *keyedNode[T, K]) reconcile(slot *loom.Slot, newItems []T) error {
	newLen := len(newItems)
	oldLen := len(n.mapped)

	// fast path for empty list
	if newLen == 0 {
		Batch(func() { n.disposeAll() })
		return slot.UnmountChildren()
	}

	// fast path for create
	if oldLen == 0 {
		Batch(func() { n.initItems(newItems) })
		return slot.RenderChildren(n.mapped...)
	}

	n.updateItems(newItems)
	if newLen > oldLen {
		return slot.AppendChildren(n.mapped[oldLen:]...)
	} else if newLen < oldLen {
		var err error
		for i := oldLen - 1; i >= newLen; i-- {
			err = slot.UnmountChild(i)
		}
		return err
	}

	return nil
}

func (n *keyedNode[T, K]) updateItems(newItems []T) {
	start := 0
	end := min(len(n.mapped), len(newItems))

	// skip common prefix
	for start < end {
		newItem := newItems[start]
		currItem := Untrack(n.items[start].Read)

		if !n.compareItems(newItem, currItem) {
			break
		}
		start++
	}

	// skip common suffix
	for end > start {
		newItem := newItems[end-1]
		currItem := Untrack(n.items[end-1].Read)

		if !n.compareItems(newItem, currItem) {
			break
		}
		end--
	}

	Batch(func() {
		// update existing items
		for i := start; i < end; i++ {
			newItem := newItems[i]
			n.items[i].Write(newItem)
			n.indices[i].Write(i)
		}

		// dipose removed items
		for i := end; i < len(n.mapped); i++ {
			n.disposeItem(i)
		}

		// resize slices
		oldLen := len(n.mapped)
		n.resizeItems(len(newItems))

		// create new items
		for i := oldLen; i < len(newItems); i++ {
			n.initItem(i, &newItems[i])
		}
	})
}

func (n *keyedNode[T, K]) resizeItems(newLen int) {
	oldLen := len(n.mapped)
	if newLen > oldLen {
		n.mapped = append(n.mapped, make([]loom.Node, newLen-oldLen)...)
		n.owners = append(n.owners, make([]*Owner, newLen-oldLen)...)
		n.items = append(n.items, make([]*sig.Signal[T], newLen-oldLen)...)
		n.indices = append(n.indices, make([]*sig.Signal[int], newLen-oldLen)...)
	} else {
		n.mapped = n.mapped[:newLen]
		n.owners = n.owners[:newLen]
		n.items = n.items[:newLen]
		n.indices = n.indices[:newLen]
	}
}

func (n *keyedNode[T, K]) disposeItem(index int) {
	if n.owners[index] != nil {
		n.owners[index].Dispose()
		n.owners[index] = nil
	}
	n.mapped[index] = nil
	n.items[index] = nil
	n.indices[index] = nil
}

func (n *keyedNode[T, K]) disposeAll() {
	n.renderOwner.Dispose()

	n.mapped = nil
	n.owners = nil
	n.items = nil
	n.indices = nil
}

func (n *keyedNode[T, K]) compareItems(a, b T) bool {
	if n.keyer == nil {
		return reflect.DeepEqual(a, b)
	}

	return reflect.DeepEqual(n.keyer(a), n.keyer(b))
}
