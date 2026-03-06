package components

import "sync/atomic"

var nextID atomic.Uint64

func newID() uint64 {
	return nextID.Add(1)
}
