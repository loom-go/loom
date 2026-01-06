package loom

func Render(parent any, node Node) error {
	slot := NewSlot()
	slot.SetNode(node)
	slot.SetParent(parent)

	return node.Mount(slot)
}
