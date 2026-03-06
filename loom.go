package loom

import (
	"fmt"
	"runtime/debug"

	"github.com/AnatoleLucet/loom/signals"
)

func Render(parent any, node Node) (*signals.Owner, error) {
	owner := signals.NewOwner()

	slot := NewSlot()
	slot.SetNode(node)
	slot.SetParent(parent)

	err := owner.Run(func() (err error) {
		defer func() {
			// should not happend during initial render, but just in case
			if r := recover(); r != nil {
				err = fmt.Errorf("%v:\n%s", r, debug.Stack())
			}
		}()

		return node.Mount(slot)
	})

	return owner, err
}
