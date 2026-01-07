package components

import (
	"testing"

	"github.com/AnatoleLucet/loom"
	"github.com/AnatoleLucet/loom/signals"
	"github.com/AnatoleLucet/loom/test"
	"github.com/stretchr/testify/assert"
)

func TestShow(t *testing.T) {
	t.Run("renders children", func(t *testing.T) {
		display, setDisplay := signals.Signal(false)
		child := test.NewMockNode("child")

		show := Show(display, func() loom.Node {
			return Fragment(Fragment(), child, Fragment())
		})
		err := loom.Render("parent", show)
		assert.NoError(t, err)

		assert.Equal(t, 0, child.MountCalls(), "child should not be mounted")
		setDisplay(true)
		assert.Equal(t, 1, child.MountCalls(), "child should be mounted once")
	})

	t.Run("remounts children", func(t *testing.T) {
		display, setDisplay := signals.Signal(false)
		child := test.NewMockNode("child")

		show := Show(display, func() loom.Node {
			return child
		})
		err := loom.Render("parent", show)
		assert.NoError(t, err)

		assert.Equal(t, 0, child.MountCalls(), "child should not be mounted")
		setDisplay(true)
		assert.Equal(t, 1, child.MountCalls(), "child should be mounted once")
		setDisplay(false)
		assert.Equal(t, 1, child.UnmountCalls(), "child should be unmounted once")
		setDisplay(true)
		assert.Equal(t, 2, child.MountCalls(), "child should be mounted twice")
	})

	t.Run("dont rerender when already displayed", func(t *testing.T) {
		display, setDisplay := signals.Signal(false)
		child := test.NewMockNode("child")

		renderFnCalls := 0
		show := Show(display, func() loom.Node {
			renderFnCalls++
			return child
		})
		err := loom.Render("parent", show)
		assert.NoError(t, err)

		assert.Equal(t, 0, child.MountCalls(), "child should not be mounted")
		assert.Equal(t, 0, renderFnCalls, "render function should not be called yet")
		setDisplay(true)
		assert.Equal(t, 1, child.MountCalls(), "child should be mounted once")
		assert.Equal(t, 1, renderFnCalls, "render function should be called once")
		setDisplay(true)
		assert.Equal(t, 1, child.MountCalls(), "child should still be mounted once")
		assert.Equal(t, 1, renderFnCalls, "render function should still be called once")
	})

	t.Run("cleans up render scope", func(t *testing.T) {
		display, setDisplay := signals.Signal(false)
		child := test.NewMockNode("child")

		cleanupCalls := 0
		show := Show(display, func() loom.Node {
			signals.OnCleanup(func() { cleanupCalls++ })
			return child
		})
		err := loom.Render("parent", show)
		assert.NoError(t, err)

		assert.Equal(t, 0, child.MountCalls(), "child should not be mounted")
		assert.Equal(t, 0, cleanupCalls, "cleanup should not be called yet")
		setDisplay(true)
		assert.Equal(t, 1, child.MountCalls(), "child should be mounted once")
		assert.Equal(t, 0, cleanupCalls, "cleanup should not be called yet")
		setDisplay(false)
		assert.Equal(t, 1, cleanupCalls, "cleanup should be called once")
	})
}
