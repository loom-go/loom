package components

import (
	"context"

	"github.com/AnatoleLucet/loom"
)

type component struct {
	ctx context.Context
}

func Self() loom.Component {
	ctx, cancel := context.WithCancel(context.Background())
	OnCleanup(cancel)

	return &component{ctx}
}

func (c *component) Context() context.Context {
	return c.ctx
}

func (c *component) IsDisposed() bool {
	return c.ctx.Err() != nil
}

func (c *component) Disposed() <-chan struct{} {
	return c.ctx.Done()
}
