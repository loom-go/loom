package loom_test

import (
	"sync"
	"testing"
	"time"

	"github.com/AnatoleLucet/loom"
	"github.com/AnatoleLucet/loom/test"
	"github.com/stretchr/testify/assert"
)

func TestSlot(t *testing.T) {
	t.Run("NewSlot creates empty slot", func(t *testing.T) {
		slot := loom.NewSlot()

		assert.Nil(t, slot.Node(), "node should be nil")
		assert.Nil(t, slot.Parent(), "parent should be nil")
		assert.Nil(t, slot.Self(), "self should be nil")
		assert.False(t, slot.Mounted(), "should not be mounted")
	})

	t.Run("Mounted reflects node state", func(t *testing.T) {
		slot := loom.NewSlot()

		assert.False(t, slot.Mounted(), "should not be mounted initially")

		slot.SetNode(test.NewMockNode("test"))
		assert.True(t, slot.Mounted(), "should be mounted after SetNode")

		slot.SetNode(nil)
		assert.False(t, slot.Mounted(), "should not be mounted after SetNode(nil)")
	})
}

func TestSlot_Child(t *testing.T) {
	t.Run("creates child slots on demand", func(t *testing.T) {
		slot := loom.NewSlot()

		child0 := slot.Child(0)
		child1 := slot.Child(1)
		child2 := slot.Child(2)

		assert.NotNil(t, child0, "child0 should be created")
		assert.NotNil(t, child1, "child1 should be created")
		assert.NotNil(t, child2, "child2 should be created")
		assert.False(t, child0.Mounted(), "child0 should not be mounted")
		assert.False(t, child1.Mounted(), "child1 should not be mounted")
		assert.False(t, child2.Mounted(), "child2 should not be mounted")
	})

	t.Run("reuses existing child slots", func(t *testing.T) {
		slot := loom.NewSlot()

		child0First := slot.Child(0)
		child0First.SetNode(test.NewMockNode("test"))

		child0Second := slot.Child(0)

		assert.Same(t, child0First, child0Second, "should return same slot")
		assert.True(t, child0Second.Mounted(), "slot should still be mounted")
	})

	t.Run("creates intermediate slots when accessing higher index", func(t *testing.T) {
		slot := loom.NewSlot()

		child5 := slot.Child(5)
		assert.NotNil(t, child5, "child5 should be created")

		// intermediate slots should also exist now
		child2 := slot.Child(2)
		assert.NotNil(t, child2, "child2 should exist")
	})
}

func TestSlot_RenderChildren(t *testing.T) {
	t.Run("mounts new children", func(t *testing.T) {
		slot := loom.NewSlot()
		slot.SetParent("parent")

		child1 := test.NewMockNode("child1")
		child2 := test.NewMockNode("child2")
		child3 := test.NewMockNode("child3")

		err := slot.RenderChildren(child1, child2, child3)
		assert.NoError(t, err)

		assert.Equal(t, 1, child1.MountCalls(), "child1 should be mounted once")
		assert.Equal(t, 1, child2.MountCalls(), "child2 should be mounted once")
		assert.Equal(t, 1, child3.MountCalls(), "child3 should be mounted once")
		assert.Equal(t, 0, child1.UpdateCalls(), "child1 should not be updated")
		assert.Equal(t, 0, child2.UpdateCalls(), "child2 should not be updated")
		assert.Equal(t, 0, child3.UpdateCalls(), "child3 should not be updated")
	})

	t.Run("updates children with matching ID", func(t *testing.T) {
		slot := loom.NewSlot()
		slot.SetParent("parent")

		child := test.NewMockNode("child")

		err := slot.RenderChildren(child)
		assert.NoError(t, err)
		assert.Equal(t, 1, child.MountCalls(), "child should be mounted once")
		assert.Equal(t, 0, child.UpdateCalls(), "child should not be updated yet")

		err = slot.RenderChildren(child)
		assert.NoError(t, err)
		assert.Equal(t, 1, child.MountCalls(), "child should still be mounted once")
		assert.Equal(t, 1, child.UpdateCalls(), "child should be updated once")

		err = slot.RenderChildren(child)
		assert.NoError(t, err)
		assert.Equal(t, 1, child.MountCalls(), "child should still be mounted once")
		assert.Equal(t, 2, child.UpdateCalls(), "child should be updated twice")
	})

	t.Run("replaces children with different ID", func(t *testing.T) {
		slot := loom.NewSlot()
		slot.SetParent("parent")

		child1 := test.NewMockNode("child1")
		child2 := test.NewMockNode("child2")

		err := slot.RenderChildren(child1)
		assert.NoError(t, err)
		assert.Equal(t, 1, child1.MountCalls(), "child1 should be mounted")

		err = slot.RenderChildren(child2)
		assert.NoError(t, err)
		assert.Equal(t, 1, child1.UnmountCalls(), "child1 should be unmounted")
		assert.Equal(t, 1, child2.MountCalls(), "child2 should be mounted")
		assert.Equal(t, 0, child2.UpdateCalls(), "child2 should not be updated")
	})

	t.Run("unmounts extra children when list shrinks", func(t *testing.T) {
		slot := loom.NewSlot()
		slot.SetParent("parent")

		child1 := test.NewMockNode("child1")
		child2 := test.NewMockNode("child2")
		child3 := test.NewMockNode("child3")
		child4 := test.NewMockNode("child4")
		child5 := test.NewMockNode("child5")

		err := slot.RenderChildren(child1, child2, child3, child4, child5)
		assert.NoError(t, err)

		err = slot.RenderChildren(child1, child2)
		assert.NoError(t, err)

		assert.Equal(t, 0, child1.UnmountCalls(), "child1 should not be unmounted")
		assert.Equal(t, 0, child2.UnmountCalls(), "child2 should not be unmounted")
		assert.Equal(t, 1, child3.UnmountCalls(), "child3 should be unmounted")
		assert.Equal(t, 1, child4.UnmountCalls(), "child4 should be unmounted")
		assert.Equal(t, 1, child5.UnmountCalls(), "child5 should be unmounted")
	})

	t.Run("handles nil children by unmounting", func(t *testing.T) {
		slot := loom.NewSlot()
		slot.SetParent("parent")

		child := test.NewMockNode("child")

		err := slot.RenderChildren(child)
		assert.NoError(t, err)
		assert.Equal(t, 1, child.MountCalls(), "child should be mounted")

		err = slot.RenderChildren(nil)
		assert.NoError(t, err)
		assert.Equal(t, 1, child.UnmountCalls(), "child should be unmounted")

		// slot should now be unmounted but still exist
		childSlot := slot.Child(0)
		assert.False(t, childSlot.Mounted(), "child slot should not be mounted")

		// re-render with same child should mount again
		err = slot.RenderChildren(child)
		assert.NoError(t, err)
		assert.Equal(t, 2, child.MountCalls(), "child should be mounted again")
	})

	t.Run("uses self as parent for children", func(t *testing.T) {
		slot := loom.NewSlot()
		slot.SetParent("parent")
		slot.SetSelf("self")

		child := test.NewMockNode("child")

		err := slot.RenderChildren(child)
		assert.NoError(t, err)

		childSlot := slot.Child(0)
		assert.Equal(t, "self", childSlot.Parent(), "child should have self as parent")
	})

	t.Run("falls back to parent when self is nil (transparent nodes)", func(t *testing.T) {
		slot := loom.NewSlot()
		slot.SetParent("parent")
		// self is nil - simulating transparent node like Fragment, Bind

		child := test.NewMockNode("child")

		err := slot.RenderChildren(child)
		assert.NoError(t, err)

		childSlot := slot.Child(0)
		assert.Equal(t, "parent", childSlot.Parent(), "child should have parent as parent")
	})

	t.Run("handles empty children list", func(t *testing.T) {
		slot := loom.NewSlot()
		slot.SetParent("parent")

		child := test.NewMockNode("child")

		err := slot.RenderChildren(child)
		assert.NoError(t, err)

		err = slot.RenderChildren()
		assert.NoError(t, err)

		assert.Equal(t, 1, child.UnmountCalls(), "child should be unmounted")
	})

	t.Run("renders nested children through node's Mount", func(t *testing.T) {
		slot := loom.NewSlot()
		slot.SetParent("parent")
		slot.SetSelf("root")

		grandchild := test.NewMockNode("grandchild")
		child := test.NewMockNode("child", grandchild)

		err := slot.RenderChildren(child)
		assert.NoError(t, err)

		assert.Equal(t, 1, child.MountCalls(), "child should be mounted")
		assert.Equal(t, 1, grandchild.MountCalls(), "grandchild should be mounted")
	})

	t.Run("mixed mount update and replace in same call", func(t *testing.T) {
		slot := loom.NewSlot()
		slot.SetParent("parent")

		child1 := test.NewMockNode("child1")
		child2 := test.NewMockNode("child2")
		child3 := test.NewMockNode("child3")

		err := slot.RenderChildren(child1, child2, child3)
		assert.NoError(t, err)

		newChild := test.NewMockNode("newChild")
		err = slot.RenderChildren(child1, newChild, child3)
		assert.NoError(t, err)

		assert.Equal(t, 1, child1.MountCalls(), "child1 should be mounted once")
		assert.Equal(t, 1, child1.UpdateCalls(), "child1 should be updated once")
		assert.Equal(t, 0, child1.UnmountCalls(), "child1 should not be unmounted")

		assert.Equal(t, 1, child2.UnmountCalls(), "child2 should be unmounted")
		assert.Equal(t, 1, newChild.MountCalls(), "newChild should be mounted")

		assert.Equal(t, 1, child3.MountCalls(), "child3 should be mounted once")
		assert.Equal(t, 1, child3.UpdateCalls(), "child3 should be updated once")
		assert.Equal(t, 0, child3.UnmountCalls(), "child3 should not be unmounted")
	})
}

func TestSlot_AppendChildren(t *testing.T) {
	t.Run("appends to existing children", func(t *testing.T) {
		slot := loom.NewSlot()
		slot.SetParent("parent")

		child1 := test.NewMockNode("child1")
		child2 := test.NewMockNode("child2")

		err := slot.RenderChildren(child1, child2)
		assert.NoError(t, err)

		child3 := test.NewMockNode("child3")
		child4 := test.NewMockNode("child4")

		err = slot.AppendChildren(child3, child4)
		assert.NoError(t, err)

		assert.Equal(t, 1, child1.MountCalls(), "child1 should be mounted once")
		assert.Equal(t, 0, child1.UpdateCalls(), "child1 should not be updated")
		assert.Equal(t, 1, child2.MountCalls(), "child2 should be mounted once")
		assert.Equal(t, 0, child2.UpdateCalls(), "child2 should not be updated")
		assert.Equal(t, 1, child3.MountCalls(), "child3 should be mounted")
		assert.Equal(t, 1, child4.MountCalls(), "child4 should be mounted")
	})

	t.Run("skips nil children", func(t *testing.T) {
		slot := loom.NewSlot()
		slot.SetParent("parent")

		child1 := test.NewMockNode("child1")
		child2 := test.NewMockNode("child2")

		err := slot.AppendChildren(nil, child1, nil, child2, nil)
		assert.NoError(t, err)

		assert.Equal(t, 1, child1.MountCalls(), "child1 should be mounted")
		assert.Equal(t, 1, child2.MountCalls(), "child2 should be mounted")
	})

	t.Run("uses self as parent", func(t *testing.T) {
		slot := loom.NewSlot()
		slot.SetParent("parent")
		slot.SetSelf("self")

		child := test.NewMockNode("child")

		err := slot.AppendChildren(child)
		assert.NoError(t, err)

		childSlot := slot.Child(0)
		assert.Equal(t, "self", childSlot.Parent(), "child should have self as parent")
	})

	t.Run("falls back to parent when self is nil", func(t *testing.T) {
		slot := loom.NewSlot()
		slot.SetParent("parent")

		child := test.NewMockNode("child")

		err := slot.AppendChildren(child)
		assert.NoError(t, err)

		childSlot := slot.Child(0)
		assert.Equal(t, "parent", childSlot.Parent(), "child should have parent as parent")
	})
}

func TestSlot_Unmount(t *testing.T) {
	t.Run("calls node Unmount", func(t *testing.T) {
		slot := loom.NewSlot()
		slot.SetParent("parent")

		node := test.NewMockNode("node")
		slot.SetNode(node)

		err := slot.Unmount()
		assert.NoError(t, err)

		assert.Equal(t, 1, node.UnmountCalls(), "node should be unmounted")
	})

	t.Run("clears node and self", func(t *testing.T) {
		slot := loom.NewSlot()
		slot.SetParent("parent")
		slot.SetSelf("self")
		slot.SetNode(test.NewMockNode("node"))

		err := slot.Unmount()
		assert.NoError(t, err)

		assert.Nil(t, slot.Node(), "node should be nil")
		assert.Nil(t, slot.Self(), "self should be nil")
		assert.False(t, slot.Mounted(), "should not be mounted")
	})

	t.Run("unmounts children first", func(t *testing.T) {
		slot := loom.NewSlot()
		slot.SetParent("parent")

		child := test.NewMockNode("child")
		parent := test.NewMockNode("parent", child)

		slot.SetNode(parent)
		err := parent.Mount(slot)
		assert.NoError(t, err)

		err = slot.Unmount()
		assert.NoError(t, err)

		assert.Equal(t, 1, child.UnmountCalls(), "child should be unmounted")
		assert.Equal(t, 1, parent.UnmountCalls(), "parent should be unmounted")
	})

	t.Run("handles unmounting when not mounted", func(t *testing.T) {
		slot := loom.NewSlot()

		err := slot.Unmount()
		assert.NoError(t, err)
	})

	t.Run("clears children slice", func(t *testing.T) {
		slot := loom.NewSlot()
		slot.SetParent("parent")

		child := test.NewMockNode("child")
		err := slot.RenderChildren(child)
		assert.NoError(t, err)

		slot.SetNode(test.NewMockNode("node"))
		err = slot.Unmount()
		assert.NoError(t, err)

		newChildSlot := slot.Child(0)
		assert.False(t, newChildSlot.Mounted(), "new child slot should not be mounted")
	})
}

func TestSlot_UnmountChildren(t *testing.T) {
	t.Run("unmounts all children", func(t *testing.T) {
		slot := loom.NewSlot()
		slot.SetParent("parent")

		child1 := test.NewMockNode("child1")
		child2 := test.NewMockNode("child2")
		child3 := test.NewMockNode("child3")

		err := slot.RenderChildren(child1, child2, child3)
		assert.NoError(t, err)

		err = slot.UnmountChildren()
		assert.NoError(t, err)

		assert.Equal(t, 1, child1.UnmountCalls(), "child1 should be unmounted")
		assert.Equal(t, 1, child2.UnmountCalls(), "child2 should be unmounted")
		assert.Equal(t, 1, child3.UnmountCalls(), "child3 should be unmounted")
	})

	t.Run("clears children slice", func(t *testing.T) {
		slot := loom.NewSlot()
		slot.SetParent("parent")

		child := test.NewMockNode("child")
		err := slot.RenderChildren(child)
		assert.NoError(t, err)

		err = slot.UnmountChildren()
		assert.NoError(t, err)

		newChildSlot := slot.Child(0)
		assert.False(t, newChildSlot.Mounted(), "new child slot should not be mounted")
	})
}

func TestSlot_UnmountChild(t *testing.T) {
	t.Run("unmounts specific child", func(t *testing.T) {
		slot := loom.NewSlot()
		slot.SetParent("parent")

		child1 := test.NewMockNode("child1")
		child2 := test.NewMockNode("child2")
		child3 := test.NewMockNode("child3")

		err := slot.RenderChildren(child1, child2, child3)
		assert.NoError(t, err)

		err = slot.UnmountChild(1)
		assert.NoError(t, err)

		assert.Equal(t, 0, child1.UnmountCalls(), "child1 should not be unmounted")
		assert.Equal(t, 1, child2.UnmountCalls(), "child2 should be unmounted")
		assert.Equal(t, 0, child3.UnmountCalls(), "child3 should not be unmounted")
	})

	t.Run("removes child from children slice", func(t *testing.T) {
		slot := loom.NewSlot()
		slot.SetParent("parent")

		child1 := test.NewMockNode("child1")
		child2 := test.NewMockNode("child2")
		child3 := test.NewMockNode("child3")

		err := slot.RenderChildren(child1, child2, child3)
		assert.NoError(t, err)

		err = slot.UnmountChild(1)
		assert.NoError(t, err)

		assert.Equal(t, child3, slot.Child(1).Node(), "child3 should be at index 1")
		assert.Equal(t, child1, slot.Child(0).Node(), "child1 should still be at index 0")
	})

	t.Run("unmounts from beginning", func(t *testing.T) {
		slot := loom.NewSlot()
		slot.SetParent("parent")

		child1 := test.NewMockNode("child1")
		child2 := test.NewMockNode("child2")
		child3 := test.NewMockNode("child3")

		err := slot.RenderChildren(child1, child2, child3)
		assert.NoError(t, err)

		err = slot.UnmountChild(0)
		assert.NoError(t, err)

		assert.Equal(t, 1, child1.UnmountCalls(), "child1 should be unmounted")
		assert.Equal(t, child2, slot.Child(0).Node(), "child2 should now be at index 0")
		assert.Equal(t, child3, slot.Child(1).Node(), "child3 should now be at index 1")
	})

	t.Run("unmounts from end", func(t *testing.T) {
		slot := loom.NewSlot()
		slot.SetParent("parent")

		child1 := test.NewMockNode("child1")
		child2 := test.NewMockNode("child2")
		child3 := test.NewMockNode("child3")

		err := slot.RenderChildren(child1, child2, child3)
		assert.NoError(t, err)

		err = slot.UnmountChild(2)
		assert.NoError(t, err)

		assert.Equal(t, 1, child3.UnmountCalls(), "child3 should be unmounted")
		assert.Equal(t, child1, slot.Child(0).Node(), "child1 should still be at index 0")
		assert.Equal(t, child2, slot.Child(1).Node(), "child2 should still be at index 1")
	})
}

func TestSlot_ReplaceWith(t *testing.T) {
	t.Run("unmounts old and mounts new", func(t *testing.T) {
		slot := loom.NewSlot()
		slot.SetParent("parent")

		oldNode := test.NewMockNode("old")
		newNode := test.NewMockNode("new")

		slot.SetNode(oldNode)

		err := slot.ReplaceWith(newNode)
		assert.NoError(t, err)

		assert.Equal(t, 1, oldNode.UnmountCalls(), "old should be unmounted")
		assert.Equal(t, 1, newNode.MountCalls(), "new should be mounted")
		assert.Equal(t, newNode, slot.Node(), "slot should have new node")
	})

	t.Run("replaces with nil", func(t *testing.T) {
		slot := loom.NewSlot()
		slot.SetParent("parent")

		node := test.NewMockNode("node")
		slot.SetNode(node)

		err := slot.ReplaceWith(nil)
		assert.NoError(t, err)

		assert.Equal(t, 1, node.UnmountCalls(), "node should be unmounted")
		assert.Nil(t, slot.Node(), "node should be nil")
		assert.False(t, slot.Mounted(), "should not be mounted")
	})

	t.Run("replaces when not mounted", func(t *testing.T) {
		slot := loom.NewSlot()
		slot.SetParent("parent")

		newNode := test.NewMockNode("new")

		err := slot.ReplaceWith(newNode)
		assert.NoError(t, err)

		assert.Equal(t, 1, newNode.MountCalls(), "new should be mounted")
		assert.Equal(t, newNode, slot.Node(), "slot should have new node")
	})

	t.Run("clears self when unmounting", func(t *testing.T) {
		slot := loom.NewSlot()
		slot.SetParent("parent")
		slot.SetSelf("old-self")
		slot.SetNode(test.NewMockNode("old"))

		newNode := test.NewMockNode("new")

		err := slot.ReplaceWith(newNode)
		assert.NoError(t, err)

		assert.Nil(t, slot.Self(), "self should be cleared")
	})

	t.Run("unmounts children when replacing", func(t *testing.T) {
		slot := loom.NewSlot()
		slot.SetParent("parent")

		grandchild := test.NewMockNode("grandchild")
		oldNode := test.NewMockNode("old", grandchild)

		slot.SetNode(oldNode)
		err := oldNode.Mount(slot)
		assert.NoError(t, err)

		newNode := test.NewMockNode("new")

		err = slot.ReplaceWith(newNode)
		assert.NoError(t, err)

		assert.Equal(t, 1, grandchild.UnmountCalls(), "grandchild should be unmounted")
		assert.Equal(t, 1, oldNode.UnmountCalls(), "old should be unmounted")
	})
}

func TestSlot_Freestyle(t *testing.T) {
	t.Run("deep nesting with unmount cascade", func(t *testing.T) {
		slot := loom.NewSlot()
		slot.SetParent("root")
		slot.SetSelf("level0")

		level3 := test.NewMockNode("level3")
		level2 := test.NewMockNode("level2", level3)
		level1 := test.NewMockNode("level1", level2)

		err := slot.RenderChildren(level1)
		assert.NoError(t, err)

		assert.Equal(t, 1, level1.MountCalls(), "level1 should be mounted")
		assert.Equal(t, 1, level2.MountCalls(), "level2 should be mounted")
		assert.Equal(t, 1, level3.MountCalls(), "level3 should be mounted")

		err = slot.UnmountChildren()
		assert.NoError(t, err)

		assert.Equal(t, 1, level1.UnmountCalls(), "level1 should be unmounted")
		assert.Equal(t, 1, level2.UnmountCalls(), "level2 should be unmounted")
		assert.Equal(t, 1, level3.UnmountCalls(), "level3 should be unmounted")
	})

	t.Run("sequential append and remove operations", func(t *testing.T) {
		slot := loom.NewSlot()
		slot.SetParent("parent")

		child1 := test.NewMockNode("child1")
		child2 := test.NewMockNode("child2")
		child3 := test.NewMockNode("child3")
		child4 := test.NewMockNode("child4")

		err := slot.RenderChildren(child1, child2)
		assert.NoError(t, err)

		err = slot.AppendChildren(child3, child4)
		assert.NoError(t, err)

		assert.Equal(t, 1, child1.MountCalls(), "child1 should be mounted once")
		assert.Equal(t, 1, child2.MountCalls(), "child2 should be mounted once")
		assert.Equal(t, 1, child3.MountCalls(), "child3 should be mounted once")
		assert.Equal(t, 1, child4.MountCalls(), "child4 should be mounted once")

		err = slot.UnmountChild(1)
		assert.NoError(t, err)

		assert.Equal(t, 1, child2.UnmountCalls(), "child2 should be unmounted")
		assert.Equal(t, child1, slot.Child(0).Node(), "child1 should be at 0")
		assert.Equal(t, child3, slot.Child(1).Node(), "child3 should be at 1")
		assert.Equal(t, child4, slot.Child(2).Node(), "child4 should be at 2")

		child5 := test.NewMockNode("child5")
		err = slot.Child(0).ReplaceWith(child5)
		assert.NoError(t, err)

		assert.Equal(t, 1, child1.UnmountCalls(), "child1 should be unmounted")
		assert.Equal(t, 1, child5.MountCalls(), "child5 should be mounted")
		assert.Equal(t, child5, slot.Child(0).Node(), "child5 should be at 0")
	})

	t.Run("shrink and grow list", func(t *testing.T) {
		slot := loom.NewSlot()
		slot.SetParent("parent")

		child1 := test.NewMockNode("child1")
		child2 := test.NewMockNode("child2")
		child3 := test.NewMockNode("child3")
		child4 := test.NewMockNode("child4")
		child5 := test.NewMockNode("child5")

		err := slot.RenderChildren(child1, child2, child3, child4, child5)
		assert.NoError(t, err)

		err = slot.RenderChildren(child1, child2)
		assert.NoError(t, err)

		assert.Equal(t, 1, child3.UnmountCalls(), "child3 should be unmounted")
		assert.Equal(t, 1, child4.UnmountCalls(), "child4 should be unmounted")
		assert.Equal(t, 1, child5.UnmountCalls(), "child5 should be unmounted")
		assert.Equal(t, 1, child1.UpdateCalls(), "child1 should be updated")
		assert.Equal(t, 1, child2.UpdateCalls(), "child2 should be updated")

		child6 := test.NewMockNode("child6")
		child7 := test.NewMockNode("child7")
		err = slot.RenderChildren(child1, child2, child6, child7)
		assert.NoError(t, err)

		assert.Equal(t, 2, child1.UpdateCalls(), "child1 should be updated again")
		assert.Equal(t, 2, child2.UpdateCalls(), "child2 should be updated again")
		assert.Equal(t, 1, child6.MountCalls(), "child6 should be mounted")
		assert.Equal(t, 1, child7.MountCalls(), "child7 should be mounted")
	})

	t.Run("updating node while its mounting", func(t *testing.T) {
		slot := loom.NewSlot()
		slot.SetParent("parent")

		child := test.NewMockNode("child")
		child.OnMount(func() {
			time.Sleep(10 * time.Millisecond) // simulate work on mount
			slot.SetSelf("mounted")
		})
		child.OnUpdate(func() {
			assert.Equal(t, "mounted", slot.Self(), "self should be set during update")
		})

		var wg sync.WaitGroup

		wg.Go(func() {
			err := slot.RenderChildren(child)
			assert.NoError(t, err)
		})

		wg.Go(func() {
			time.Sleep(1 * time.Millisecond) // ensure this runs during mount
			err := slot.RenderChildren(child)
			assert.NoError(t, err)
		})

		wg.Wait()
	})
}
