package components

import (
	"github.com/loom-go/loom/signals"
)

// mainly used to re-export what's in loom/signals for users to import with loom/components

type Accessor[T any] = signals.Accessor[T]

func Signal[T any](initial T) (Accessor[T], func(T)) {
	return signals.Signal(initial)
}

func Memo[T any](fn func() T) Accessor[T] {
	return signals.Memo(fn)
}

func Effect(effect func()) {
	signals.Effect(effect)
}

func RenderEffect(effect func()) {
	signals.RenderEffect(effect)
}

func Batch(fn func()) {
	signals.Batch(fn)
}

func OnCleanup(fn func()) {
	signals.OnCleanup(fn)
}

func OnSettled(fn func()) {
	signals.OnSettled(fn)
}

func OnUserSettled(fn func()) {
	signals.OnUserSettled(fn)
}

func OnRenderSettled(fn func()) {
	signals.OnRenderSettled(fn)
}

func Untrack[T any](fn func() T) T {
	return signals.Untrack(fn)
}

type Owner = signals.Owner

func NewOwner() *Owner {
	return signals.NewOwner()
}

type Writable[T any] = signals.Writable[T]

func NewWritable[T any](initial T) *Writable[T] {
	return signals.NewWritable(initial)
}

type Mutable[T any] = signals.Mutable[T]

func NewMutable[T any](initial T) *Mutable[T] {
	return signals.NewMutable(initial)
}
