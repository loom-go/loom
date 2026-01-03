package components

import (
	"github.com/AnatoleLucet/loom"
)

// Fragment groups multiple children without adding extra nodes to the output.
func Fragment(children ...loom.Node) loom.Node {
	return loom.NodeFunc(func(ctx *loom.RenderContext) error {
		for _, child := range children {
			err := child.Render(ctx)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
