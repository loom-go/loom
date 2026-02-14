package loom

// all this should be context bounded.
// but does that even make any sens?
// maybe it should only block reactive updates. BUT if they are deffered, this probably means the reactive ctx and ownership
// is going to be all wrong. maybe we could fix that by manually controlling the ownership in deffered updates.

// owner := NewOwner()
// Effect(func(){
//   OnCleanup(func() {
//		 owner.Dispose()
//   })
//   node := fn()
//   err := owner.Run(func() error { return slot.RenderChildren(node) })
// })

func LockTree() {
	// todo: should make sure the tree cannot be changed while this is locked. updated should be non blocking and deffered until unlock.
}

func UnlockTree() {
	// ... allow updates again, and run deffered updates
}

func IsTreeLocked() bool {
	// ... return true if the tree is currently locked
	return false
}

func DeferUpdate(fn func()) {
	if IsTreeLocked() {
		// ... add fn to a queue of deferred updates to be run when the tree is unlocked
		return
	}

	fn()
}
