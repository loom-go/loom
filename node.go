package loom

type Node interface {
	ID() string
	Mount(slot *Slot) error
	Update(slot *Slot) error
	Unmount(slot *Slot) error
}
