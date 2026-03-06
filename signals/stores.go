package signals

type Writable[T any] struct {
	get func() T
	set func(T)
}

func NewWritable[T any](initial T) *Writable[T] {
	get, set := Signal(initial)

	return &Writable[T]{get: get, set: set}
}

func (s *Writable[T]) Get() T {
	return s.get()
}

func (s *Writable[T]) Set(value T) {
	s.set(value)
}

func (s *Writable[T]) Update(fn func(T) T) {
	s.set(fn(s.get()))
}

type Mutable[T any] struct {
	get func() T
	set func(T)
}

func NewMutable[T any](initial T) *Mutable[T] {
	opts := SignalOptions[T]{
		// always consider the value changed
		Predicate: func(a, b T) bool { return false },
	}

	get, set := Signal(initial, opts)

	return &Mutable[T]{get: get, set: set}
}

func (s *Mutable[T]) Get() T {
	return s.get()
}

func (s *Mutable[T]) Set(value T) {
	s.set(value)
}

func (s *Mutable[T]) Mutate(fn func(*T)) {
	v := s.get()
	fn(&v)
	s.set(v)
}
