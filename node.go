package loom

import (
	"context"
	"sync"
)

type Node interface {
	Render(ctx *RenderContext) error
}

type NodeFunc func(ctx *RenderContext) error

func (n NodeFunc) Render(ctx *RenderContext) error {
	return n(ctx)
}

// todo: rename to RenderContext
type RenderContext struct {
	mu  sync.RWMutex
	ctx context.Context
}

func NewRenderContext() *RenderContext {
	return &RenderContext{ctx: context.Background()}
}

func (c *RenderContext) Get(key string) any {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.ctx.Value(key)
}

func (c *RenderContext) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.ctx = context.WithValue(c.ctx, key, value)
}

func (c *RenderContext) Clone() *RenderContext {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return &RenderContext{ctx: c.ctx}
}
