package components

import (
	"github.com/AnatoleLucet/loom"
	"github.com/AnatoleLucet/sig"
)

// Show conditionally renders the given node if the when function returns true.
func Show(when func() bool, node loom.Node) loom.Node {
	o := sig.NewOwner()

	return Bind(func() loom.Node {
		if when() {
			return loom.NodeFunc(func(ctx *loom.RenderContext) error {
				return o.Run(func() error {
					return node.Render(ctx)
				})
			})
		}

		o.Dispose()
		return Fragment()
	})
}
