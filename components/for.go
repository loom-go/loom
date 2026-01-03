package components

import (
	"reflect"

	"github.com/AnatoleLucet/loom"
	"github.com/AnatoleLucet/loom/signals"
	"github.com/AnatoleLucet/sig"
)

// ForMapper represents a function that maps each item of a For() to a Node, optionally returning a key.
type ForMapper[T any] interface {
	func(signals.Accessor[T], signals.Accessor[int]) loom.Node | func(signals.Accessor[T], signals.Accessor[int]) (any, loom.Node)
}

// type forRendered[T any] struct {
// 	item  T
// 	node  Node
// 	owner interface{ Dispose() } // todo: export sig.owner
// }

type ForNode[T any] struct {
	input  signals.Accessor[[]T]
	mapper func(signals.Accessor[T], signals.Accessor[int]) (any, loom.Node)

	// whether the user provide a key for each items
	keyed bool

	// map of keys to track existing items
	keys map[any]any

	// currently rendered items
	mapped []loom.Node
	// owners of the rendered items
	owners []interface{ Dispose() } // todo: export sig.owner
	// signals for each item and index
	items   []*sig.Signal[T]
	indexes []*sig.Signal[int]
}

func For[T any, M ForMapper[T]](items signals.Accessor[[]T], mapper M) loom.Node {
	keyed, mappern := normalizeForMapper[T](mapper)

	return &ForNode[T]{
		input:  items,
		mapper: mappern,

		keyed: keyed,
	}
}

func (n *ForNode[T]) Render(ctx *loom.RenderContext) (err error) {
	initial := true
	sig.NewEffect(func() {
		defer func() {
			if err != nil && !initial {
				panic(err)
			}
			initial = false
		}()

		newItems := n.input()

		// fast path for empty list
		if len(newItems) == 0 {
			n.dispose()
			n.render(ctx)
			return
		}

		// fast path for create
		if len(n.mapped) == 0 {
			n.init(newItems)
			n.render(ctx)
			return
		}

	})

	return err
}

func (n *ForNode[T]) init(items []T) {
	n.mapped = make([]loom.Node, len(items))
	n.owners = make([]interface{ Dispose() }, len(items))
	n.items = make([]*sig.Signal[T], len(items))
	n.indexes = make([]*sig.Signal[int], len(items))
	n.clearKeys()

	for i, item := range items {
		itemSignal := sig.NewSignal(item)
		indexSignal := sig.NewSignal(i)

		var key any
		var node loom.Node
		owner := sig.NewOwner()
		owner.Run(func() error {
			key, node = n.mapper(itemSignal.Read, indexSignal.Read)
			return nil
		})

		n.mapped[i] = node
		n.owners[i] = owner
		n.items[i] = itemSignal
		n.indexes[i] = indexSignal

		if n.keyed {
			n.addKey(&items[i], key)
		}
	}
}

// func (n *ForNode[T]) render(ctx *Context, oldItems, newItems []T) error {
func (n *ForNode[T]) render(ctx *loom.RenderContext) error {
	for _, node := range n.mapped {
		// if n.compare(&oldItems[i], &newItems[i]) {
		// 	continue
		// }

		err := node.Render(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (n *ForNode[T]) dispose() {
	for _, owner := range n.owners {
		owner.Dispose()
	}

	n.mapped = nil
	n.owners = nil
	n.items = nil
	n.indexes = nil
	n.clearKeys()
}

func (n *ForNode[T]) compare(a, b *T) bool {
	if !n.keyed {
		return reflect.DeepEqual(a, b)
	}

	// todo: should render child to get the key, if !hasKey(a|b)

	if !n.hasKey(a) || !n.hasKey(b) {
		return false
	}

	return reflect.DeepEqual(n.getKey(a), n.getKey(b))
}

func (n *ForNode[T]) getKey(item *T) any      { return n.keys[item] }
func (n *ForNode[T]) addKey(item *T, key any) { n.keys[item] = key }
func (n *ForNode[T]) removeKey(item *T)       { delete(n.keys, item) }
func (n *ForNode[T]) hasKey(item *T) bool     { _, exists := n.keys[item]; return exists }
func (n *ForNode[T]) clearKeys()              { n.keys = make(map[any]any) }

// ugly function to map the ForChild union into a regular func type that always returns (any, Node)
func normalizeForMapper[T any, C ForMapper[T]](child C) (bool, func(signals.Accessor[T], signals.Accessor[int]) (any, loom.Node)) {
	switch render := any(child).(type) {
	case func(signals.Accessor[T], signals.Accessor[int]) loom.Node:
		return false, func(item signals.Accessor[T], index signals.Accessor[int]) (any, loom.Node) {
			return nil, render(item, index)
		}
	case func(signals.Accessor[T], signals.Accessor[int]) (any, loom.Node):
		return true, render
	default:
		panic("For: child must be of type func(Accessor[T], Accessor[int]) Node or func(Accessor[T], Accessor[int]) (any, Node)")
	}
}
