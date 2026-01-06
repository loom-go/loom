package loom

import (
	"sync"
)

// Slot represents a space where a Node can be rendered.
type Slot struct {
	mu sync.RWMutex

	parent any
	self   any

	node     Node
	children []*Slot
}

func NewSlot() *Slot {
	return &Slot{}
}

func (s *Slot) Node() Node {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.node
}

func (s *Slot) SetNode(node Node) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.node = node
}

func (s *Slot) Parent() any {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.parent
}

func (s *Slot) SetParent(parent any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.parent = parent
}

func (s *Slot) Self() any {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.self
}

func (s *Slot) SetSelf(self any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.self = self
}

func (s *Slot) Mounted() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.node != nil
}

func (s *Slot) Child(index int) *Slot {
	s.mu.Lock()
	defer s.mu.Unlock()

	for len(s.children) <= index {
		s.children = append(s.children, NewSlot())
	}
	return s.children[index]
}

func (s *Slot) RenderChildren(children ...Node) error {
	parent := s.self
	if parent == nil {
		// fallback to parent for transparent nodes (fragement, bind, show, etc.)
		parent = s.parent
	}

	for i, child := range children {
		childSlot := s.Child(i)
		childSlot.SetParent(parent)

		var err error
		if child == nil {
			err = childSlot.ReplaceWith(nil)
		} else if !childSlot.Mounted() {
			childSlot.SetNode(child)
			err = child.Mount(childSlot)
		} else if childSlot.Node().ID() == child.ID() {
			err = child.Update(childSlot)
		} else {
			err = childSlot.ReplaceWith(child)
		}

		if err != nil {
			return err
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// unmount extra children
	if extra := len(s.children) - len(children); extra > 0 {
		start := len(children)

		for _, child := range s.children[start:] {
			err := child.Unmount()
			if err != nil {
				return err
			}
		}

		s.children = s.children[:start]
	}

	return nil
}

func (s *Slot) AppendChildren(children ...Node) error {
	parent := s.self
	if parent == nil {
		parent = s.parent
	}

	start := len(s.children)
	for i, child := range children {
		childSlot := s.Child(start + i)
		childSlot.SetParent(parent)

		if child == nil {
			continue
		}

		childSlot.SetNode(child)
		err := child.Mount(childSlot)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Slot) Unmount() error {
	s.UnmountChildren()

	if s.Mounted() {
		err := s.Node().Unmount(s)
		if err != nil {
			return err
		}
	}

	s.SetNode(nil)
	s.SetSelf(nil)

	return nil
}

func (s *Slot) UnmountChildren() error {
	for _, child := range s.children {
		err := child.Unmount()
		if err != nil {
			return err
		}
	}
	s.children = nil
	return nil
}

func (s *Slot) UnmountChild(index int) error {
	childSlot := s.Child(index)
	err := childSlot.Unmount()
	if err != nil {
		return err
	}

	s.mu.Lock()
	s.children = append(s.children[:index], s.children[index+1:]...)
	s.mu.Unlock()

	return nil
}

func (s *Slot) ReplaceWith(node Node) error {
	err := s.Unmount()
	if err != nil {
		return err
	}

	s.SetNode(node)
	if node != nil {
		return node.Mount(s)
	}

	return nil
}
