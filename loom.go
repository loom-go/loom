package loom

func Render(parent any, node Node) (*Slot, error) {
	slot := NewSlot()
	slot.SetNode(node)
	slot.SetParent(parent)

	return slot, node.Mount(slot)
}
