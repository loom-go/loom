package test

import "github.com/loom-go/loom"

type MockNode struct {
	id string

	mountCalls     int
	mountListeners []func()

	updateCalls     int
	updateListeners []func()

	unmountCalls     int
	unmountListeners []func()

	children []loom.Node
}

func NewMockNode(id string, children ...loom.Node) *MockNode {
	return &MockNode{id: id, children: children}
}

func (n *MockNode) ID() string {
	return n.id
}

func (n *MockNode) Mount(slot *loom.Slot) error {
	n.mountCalls++

	for _, listener := range n.mountListeners {
		listener()
	}

	return slot.RenderChildren(n.children...)
}

func (n *MockNode) Update(slot *loom.Slot) error {
	n.updateCalls++

	for _, listener := range n.updateListeners {
		listener()
	}

	return slot.RenderChildren(n.children...)
}

func (n *MockNode) Unmount(slot *loom.Slot) error {
	n.unmountCalls++

	for _, listener := range n.unmountListeners {
		listener()
	}

	return nil
}

func (n *MockNode) MountCalls() int {
	return n.mountCalls
}

func (n *MockNode) OnMount(listener func()) {
	n.mountListeners = append(n.mountListeners, listener)
}

func (n *MockNode) UpdateCalls() int {
	return n.updateCalls
}

func (n *MockNode) OnUpdate(listener func()) {
	n.updateListeners = append(n.updateListeners, listener)
}

func (n *MockNode) UnmountCalls() int {
	return n.unmountCalls
}

func (n *MockNode) OnUnmount(listener func()) {
	n.unmountListeners = append(n.unmountListeners, listener)
}
