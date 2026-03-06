package signals

import "github.com/AnatoleLucet/sig"

type Accessor[T any] func() T

type SignalOptions[T any] = sig.SignalOptions[T]

func Signal[T any](initial T, options ...SignalOptions[T]) (Accessor[T], func(T)) {
	s := sig.NewSignal(initial, options...)
	return s.Read, s.Write
}

func Memo[T any](fn func() T) Accessor[T] {
	c := sig.NewComputed(fn)
	return c.Read
}

func Effect(effect func()) {
	sig.NewEffect(effect)
}

func RenderEffect(effect func()) {
	sig.NewRenderEffect(effect)
}

func Untrack[T any](fn func() T) T {
	return sig.Untrack(fn)
}

func Batch(fn func()) {
	sig.NewBatch(fn)
}

func OnCleanup(fn func()) {
	sig.OnCleanup(fn)
}

func OnSettled(fn func()) {
	sig.OnSettled(fn)
}

func OnUserSettled(fn func()) {
	sig.OnUserSettled(fn)
}

func OnRenderSettled(fn func()) {
	sig.OnRenderSettled(fn)
}
