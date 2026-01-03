package components

import (
	"github.com/AnatoleLucet/loom"
	"github.com/AnatoleLucet/sig"
)

// Bind creates a reactive Node that re-renders whenever any signal
// accessed within the provided function changes.
func Bind(fn func() loom.Node) loom.Node {
	return loom.NodeFunc(func(ctx *loom.RenderContext) (err error) {
		o := sig.NewOwner()

		initial := true
		sig.NewEffect(func() {
			node := fn()
			// todo: owner doesn't seem to be needed. check if it can be removed
			err = o.Run(func() error {
				return node.Render(ctx)
			})

			if !initial && err != nil {
				panic(err)
			}

			initial = false
		})

		return err
	})
}
