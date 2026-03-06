package components

import (
	"fmt"
	"sync/atomic"

	"github.com/AnatoleLucet/loom"
)

func Apply(appliers ...loom.Applier) loom.Node {
	return &applyNode{appliers: appliers}
}

type applyNode struct {
	appliers []loom.Applier
}

func (n *applyNode) ID() string {
	return "loom.Apply"
}

func (n *applyNode) Mount(slot *loom.Slot) error {
	if err := n.apply(slot, false); err != nil {
		return fmt.Errorf("Apply: %w", err)
	}

	return nil
}

func (n *applyNode) Update(slot *loom.Slot) error {
	if err := n.apply(slot, true); err != nil {
		return fmt.Errorf("Apply: %w", err)
	}

	return nil
}

func (n *applyNode) Unmount(slot *loom.Slot) error {
	if err := n.remove(slot); err != nil {
		return fmt.Errorf("Apply: %w", err)
	}

	if err := n.refreshStack(slot); err != nil {
		return fmt.Errorf("Apply: %w", err)
	}

	return nil
}

func (n *applyNode) apply(slot *loom.Slot, force bool) (err error) {
	// remove old layer if any
	if err := n.remove(slot); err != nil {
		return err
	}

	stack := getApplierStack(slot.Parent())

	// add new layer
	layer := &applierLayer{
		id:       newID(),
		appliers: n.appliers,
	}
	slot.SetSelf(layer)
	stack.pushLayer(layer)

	var initial atomic.Bool
	RenderEffect(func() {
		if !initial.Load() {
			initial.Store(true)
		}

		if force {
			err = n.refreshStack(slot)
			if err != nil && !initial.Load() {
				panic(err) // propagate to the reactive owner
			}
		} else {
			err = layer.apply(slot.Parent())
			force = true // force apply on reactive updates to make sure unset values in the applier get proper fallback
		}
	})

	return err
}

func (n *applyNode) remove(slot *loom.Slot) error {
	self := slot.Self()
	if self == nil {
		return nil
	}

	layer := self.(*applierLayer)
	stack := getApplierStack(slot.Parent())

	stack.popLayer(layer.id)
	return layer.remove(slot.Parent())
}

func (n *applyNode) refreshStack(slot *loom.Slot) error {
	stack := getApplierStack(slot.Parent())

	for _, layer := range stack.layers {
		var err error
		if layer == slot.Self() {
			err = layer.apply(slot.Parent())
		} else {
			// untrack applier layers that's not ours. each apply tracks its own layer
			err = Untrack(func() error {
				return layer.apply(slot.Parent())
			})
		}

		if err != nil {
			return err
		}
	}

	return nil
}
