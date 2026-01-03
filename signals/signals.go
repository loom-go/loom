package signals

import "github.com/AnatoleLucet/sig"

type Accessor[T any] func() T

func Signal[T any](initial T) (Accessor[T], func(T)) {
	s := sig.NewSignal(initial)
	return s.Read, s.Write
}

func Memo[T any](fn func() T) Accessor[T] {
	c := sig.NewComputed(fn)
	return c.Read
}

func Effect(effect func()) {
	sig.NewEffect(effect)
}

func Batch(fn func()) {
	sig.NewBatch(fn)
}

func OnCleanup(fn func()) {
	sig.OnCleanup(fn)
}

func Untrack[T any](fn func() T) T {
	return sig.Untrack(fn)
}
