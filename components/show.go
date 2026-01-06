package components

import (
	"github.com/AnatoleLucet/loom"
	"github.com/AnatoleLucet/loom/signals"
)

// Show conditionally renders the given node if the when function returns true.
func Show(when func() bool, fn func() loom.Node) loom.Node {
	o := signals.NewOwner()

	var node loom.Node
	return Bind(func() loom.Node {
		if !when() {
			o.Dispose()
			node = nil
			return nil
		}

		if node == nil {
			o.Run(func() error {
				node = fn()
				return nil
			})
		}

		return Own(o, node)
	})
}
