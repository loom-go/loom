package components

import (
	"testing"

	"github.com/AnatoleLucet/loom"
	"github.com/AnatoleLucet/loom/test"
	"github.com/stretchr/testify/assert"
)

func TestBind(t *testing.T) {
	t.Run("renders children", func(t *testing.T) {
		child1 := test.NewMockNode("child1")
		child2 := test.NewMockNode("child2")

		bind := Bind(func() loom.Node {
			return Fragment(child1, child2)
		})
		_, err := loom.Render("parent", bind)
		assert.NoError(t, err)

		assert.Equal(t, 1, child1.MountCalls(), "child1 should be mounted once")
		assert.Equal(t, 1, child2.MountCalls(), "child2 should be mounted once")
	})

	t.Run("updates children", func(t *testing.T) {
		count, setCount := Signal(0)
		child := test.NewMockNode("child")

		bind := Bind(func() loom.Node {
			count()
			return child
		})
		_, err := loom.Render("parent", bind)
		assert.NoError(t, err)

		assert.Equal(t, 1, child.MountCalls(), "child should be mounted once")
		setCount(1)
		assert.Equal(t, 1, child.UpdateCalls(), "child should be updated once")
	})

	t.Run("updates nested children", func(t *testing.T) {
		count, setCount := Signal(0)
		child := test.NewMockNode("child")

		bind := Bind(func() loom.Node {
			count()
			return Fragment(Fragment(), child, Fragment())
		})

		_, err := loom.Render("parent", bind)
		assert.NoError(t, err)

		assert.Equal(t, 1, child.MountCalls(), "child should be mounted once")
		setCount(1)
		assert.Equal(t, 1, child.UpdateCalls(), "child should be updated once")
	})

	t.Run("unmounts children", func(t *testing.T) {
		count, setCount := Signal(0)
		child := test.NewMockNode("child")

		bind := Bind(func() loom.Node {
			if count()%2 == 0 {
				return Fragment(Fragment(), child, Fragment())
			}
			return Fragment()
		})
		_, err := loom.Render("parent", bind)
		assert.NoError(t, err)

		assert.Equal(t, 1, child.MountCalls(), "child should be mounted once")
		setCount(1)
		assert.Equal(t, 1, child.UnmountCalls(), "child should be unmounted once")
		setCount(2)
		assert.Equal(t, 2, child.MountCalls(), "child should be mounted twice")
	})

	t.Run("cleanups the render function", func(t *testing.T) {
		cleanupCalls := 0
		count, setCount := Signal(0)

		bind := Bind(func() loom.Node {
			OnCleanup(func() { cleanupCalls++ })
			count()
			return Fragment()
		})
		_, err := loom.Render("parent", bind)
		assert.NoError(t, err)

		assert.Equal(t, 0, cleanupCalls, "cleanup should not be called yet")
		setCount(1)
		assert.Equal(t, 1, cleanupCalls, "cleanup should be called once")
		setCount(2)
		assert.Equal(t, 2, cleanupCalls, "cleanup should be called twice")
	})

	t.Run("nested Bind calls", func(t *testing.T) {
		countA, setCountA := Signal(0)
		countB, setCountB := Signal(0)
		childA := test.NewMockNode("childA")
		childB := test.NewMockNode("childB")

		bind := Bind(func() loom.Node {
			countA()
			return Fragment(childA, Bind(func() loom.Node {
				countB()
				return childB
			}))
		})
		_, err := loom.Render("parent", bind)
		assert.NoError(t, err)

		assert.Equal(t, 1, childA.MountCalls(), "childA should be mounted once")
		assert.Equal(t, 1, childB.MountCalls(), "childB should be mounted once")

		setCountB(1)
		assert.Equal(t, 0, childA.UpdateCalls(), "childA should not be updated")
		assert.Equal(t, 1, childB.UpdateCalls(), "childB should be updated once")
		setCountA(1)
		assert.Equal(t, 1, childA.UpdateCalls(), "childA should be updated once")
		assert.Equal(t, 1, childB.UpdateCalls(), "childB should not be updated again")
	})

	t.Run("conditional nested Bind calls", func(t *testing.T) {
		count, setCount := Signal(0)
		visible, setVisible := Signal(true)

		child := test.NewMockNode("child")

		rootBindCalls := 0
		innerBindCalls := 0
		innerCleanupCalls := 0

		bind := Bind(func() loom.Node {
			rootBindCalls++

			if !visible() {
				return Fragment()
			}

			return Fragment(Bind(func() loom.Node {
				innerBindCalls++
				OnCleanup(func() { innerCleanupCalls++ })

				count()
				return child
			}))
		})
		_, err := loom.Render("parent", bind)
		assert.NoError(t, err)

		assert.Equal(t, 1, rootBindCalls, "root bind should be called once")
		assert.Equal(t, 1, innerBindCalls, "inner bind should be called once")
		assert.Equal(t, 1, child.MountCalls(), "child should be mounted once")
		assert.Equal(t, 0, innerCleanupCalls, "inner cleanup should not be called yet")

		setVisible(false)
		assert.Equal(t, 2, rootBindCalls, "root bind should be called twice")
		assert.Equal(t, 1, innerBindCalls, "inner bind should not be called again")
		assert.Equal(t, 1, child.UnmountCalls(), "child should be unmounted once")
		assert.Equal(t, 1, innerCleanupCalls, "inner cleanup should be called once")

		setCount(1)
		assert.Equal(t, 2, rootBindCalls, "root bind should not be called again")
		assert.Equal(t, 1, innerBindCalls, "inner bind should not be called again")
		assert.Equal(t, 1, child.UnmountCalls(), "child should still be unmounted")
		assert.Equal(t, 1, innerCleanupCalls, "inner cleanup should still be called once")
	})
}
