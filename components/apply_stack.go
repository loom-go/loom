package components

import "github.com/AnatoleLucet/loom"

// todo: none of that is concurency safe. prob will need to be

var applierStacks = make(map[any]*applierStack)

func getApplierStack(parent any) *applierStack {
	if stack, ok := applierStacks[parent]; ok {
		return stack
	}

	stack := &applierStack{}
	applierStacks[parent] = stack
	return stack
}

type applierLayer struct {
	id       uint64
	appliers []loom.Applier
	removers []func() error
}

func (l *applierLayer) apply(parent any) error {
	if err := l.remove(parent); err != nil {
		return err
	}

	for _, applier := range l.appliers {
		remove, err := applier.Apply(parent)
		if err != nil {
			return err
		}

		l.removers = append(l.removers, remove)
	}

	return nil
}

func (l *applierLayer) remove(parent any) error {
	for _, remove := range l.removers {
		if remove == nil {
			continue
		}

		if err := remove(); err != nil {
			return err
		}
	}
	l.removers = nil

	return nil
}

type applierStack struct {
	layers []*applierLayer
}

func (s *applierStack) getLayer(id uint64) *applierLayer {
	for i := len(s.layers) - 1; i >= 0; i-- {
		if s.layers[i].id == id {
			return s.layers[i]
		}
	}

	return nil
}

func (s *applierStack) pushLayer(layer *applierLayer) {
	s.layers = append(s.layers, layer)
}

func (s *applierStack) popLayer(id uint64) {
	for i := len(s.layers) - 1; i >= 0; i-- {
		if s.layers[i].id == id {
			s.layers = append(s.layers[:i], s.layers[i+1:]...)
			return
		}
	}
}
