package components

import (
	"github.com/AnatoleLucet/loom"
	"github.com/AnatoleLucet/loom/signals"
)

// Bind creates a reactive Node that re-renders whenever any signal
// accessed within the provided function changes.
func Bind(fn func() loom.Node) loom.Node {
	return loom.NodeFunc(func(ctx *loom.RenderContext) (err error) {
		// create an owner outside of the effect
		// so the node is not disposed on re-renders
		o := signals.NewOwner()

		initial := true
		signals.Effect(func() {
			node := fn()

			err = o.Run(func() error { return node.Render(ctx) })
			if err != nil && !initial {
				panic(err)
			}

			initial = false
		})

		return err
	})
}
