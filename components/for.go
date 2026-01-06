package components

import (
	"errors"
	"reflect"

	"github.com/AnatoleLucet/loom"
	"github.com/AnatoleLucet/loom/signals"
	"github.com/AnatoleLucet/sig"
)

// ForKeyer represents either a key function or a mapper function for For().
type ForKeyer[T any] interface {
	func(T) any | func(signals.Accessor[T], signals.Accessor[int]) loom.Node
}

func For[
	T any,
	Keyer ForKeyer[T],
](
	items signals.Accessor[[]T],
	keyer Keyer,
	mapper ...func(signals.Accessor[T], signals.Accessor[int]) loom.Node,
) loom.Node {
	keyerFn, mapperFn, err := parseForKeyer(keyer, mapper...)
	if err != nil {
		panic(err)
	}

	return &forNode[T]{
		input: items,

		mapper: mapperFn,
		keyer:  keyerFn,

		renderOwner: signals.NewOwner(),
	}
}

type forNode[T any] struct {
	input signals.Accessor[[]T]

	keyer  func(T) any
	mapper func(signals.Accessor[T], signals.Accessor[int]) loom.Node

	// currently rendered items
	mapped []loom.Node
	// owners of the rendered items
	owners []*signals.Owner
	// signals for each item and index
	items   []*sig.Signal[T]
	indexes []*sig.Signal[int]

	// owns the rendered children
	renderOwner *signals.Owner
}

func (n *forNode[T]) ID() string {
	return "loom.For"
}

func (n *forNode[T]) Mount(slot *loom.Slot) (err error) {
	initial := true

	signals.Effect(func() {
		newItems := n.input()

		err = n.renderOwner.Run(func() error {
			// fast path for empty list
			if len(newItems) == 0 {
				n.disposeItems()
				return slot.UnmountChildren()
			}

			// fast path for create
			if len(n.mapped) == 0 {
				n.initItems(newItems)
				return slot.RenderChildren(n.mapped...)
			}

			// update existing items
			oldLen := len(n.mapped)
			n.updateItems(newItems)
			newLen := len(n.mapped)

			if newLen > oldLen {
				err = slot.AppendChildren(n.mapped[oldLen:]...)
			} else if newLen < oldLen {
				for i := oldLen - 1; i >= newLen; i-- {
					err = slot.UnmountChild(i)
				}
			}

			return err
		})

		if err != nil && !initial {
			panic(err)
		}
		initial = false
	})

	return err
}

func (n *forNode[T]) Update(slot *loom.Slot) error {
	return nil
}

func (n *forNode[T]) Unmount(slot *loom.Slot) error {
	n.disposeItems()
	return nil
}

func (n *forNode[T]) initItems(items []T) {
	n.mapped = make([]loom.Node, len(items))
	n.owners = make([]*signals.Owner, len(items))
	n.items = make([]*sig.Signal[T], len(items))
	n.indexes = make([]*sig.Signal[int], len(items))

	for i, item := range items {
		n.initItem(i, &item)
	}
}

func (n *forNode[T]) initItem(index int, item *T) {
	itemSignal := sig.NewSignal(*item)
	indexSignal := sig.NewSignal(index)

	var node loom.Node
	owner := signals.NewOwner()
	owner.Run(func() error {
		node = n.mapper(itemSignal.Read, indexSignal.Read)
		return nil
	})

	n.mapped[index] = node
	n.owners[index] = owner
	n.items[index] = itemSignal
	n.indexes[index] = indexSignal
}

func (n *forNode[T]) updateItems(newItems []T) {
	start := 0
	end := min(len(n.mapped), len(newItems))

	// skip common prefix
	for start < end {
		newItem := newItems[start]
		currItem := signals.Untrack(n.items[start].Read)

		if !n.compareItems(newItem, currItem) {
			break
		}
		start++
	}

	// skip common suffix
	for end > start {
		newItem := newItems[end-1]
		currItem := signals.Untrack(n.items[end-1].Read)

		if !n.compareItems(newItem, currItem) {
			break
		}
		end--
	}

	signals.Batch(func() {
		// update existing items
		for i := start; i < end; i++ {
			newItem := newItems[i]
			n.items[i].Write(newItem)
			n.indexes[i].Write(i)
		}

		// dispose removed items
		for i := end; i < len(n.mapped); i++ {
			n.owners[i].Dispose()
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

func (n *forNode[T]) resizeItems(newLen int) {
	oldLen := len(n.mapped)
	if newLen > oldLen {
		n.mapped = append(n.mapped, make([]loom.Node, newLen-oldLen)...)
		n.owners = append(n.owners, make([]*signals.Owner, newLen-oldLen)...)
		n.items = append(n.items, make([]*sig.Signal[T], newLen-oldLen)...)
		n.indexes = append(n.indexes, make([]*sig.Signal[int], newLen-oldLen)...)
	} else {
		n.mapped = n.mapped[:newLen]
		n.owners = n.owners[:newLen]
		n.items = n.items[:newLen]
		n.indexes = n.indexes[:newLen]
	}
}

func (n *forNode[T]) disposeItems() {
	n.renderOwner.Dispose()

	n.mapped = nil
	n.owners = nil
	n.items = nil
	n.indexes = nil
}

func (n *forNode[T]) compareItems(a, b T) bool {
	if n.keyer == nil {
		return reflect.DeepEqual(a, b)
	}

	return reflect.DeepEqual((n.keyer)(a), (n.keyer)(b))
}

func parseForKeyer[
	T any,
	Keyer ForKeyer[T],
](
	keyer Keyer,
	mapper ...func(signals.Accessor[T], signals.Accessor[int]) loom.Node,
) (
	func(T) any,
	func(signals.Accessor[T], signals.Accessor[int]) loom.Node,
	error,
) {
	mapperFn, ok := any(keyer).(func(signals.Accessor[T], signals.Accessor[int]) loom.Node)
	if ok {
		if len(mapper) > 0 {
			return nil, nil, errors.New("For: expected at most one mapper function")
		}

		return nil, mapperFn, nil
	}

	keyerFn, keyed := any(keyer).(func(T) any)
	if !keyed {
		return nil, nil, errors.New("For: expected keyer to be either a key function or a mapper function")
	}
	if len(mapper) == 0 {
		return nil, nil, errors.New("For: expected mapper function")
	}
	if len(mapper) > 1 {
		return nil, nil, errors.New("For: expected at most one mapper function")
	}

	return keyerFn, mapper[0], nil
}
