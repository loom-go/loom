package components

import (
	"fmt"
	"slices"
	"testing"

	"github.com/loom-go/loom"
	"github.com/loom-go/loom/test"
	"github.com/stretchr/testify/assert"
)

func TestFor(t *testing.T) {
	removeChild := func(children []*test.MockNode, child *test.MockNode) []*test.MockNode {
		index := slices.Index(children, child)
		if index == -1 {
			return children
		}
		return slices.Delete(children, index, index+1)
	}

	t.Run("renders items", func(t *testing.T) {
		items, _ := Signal([]string{"A", "B", "C"})

		var children []*test.MockNode
		forNode := For(
			items,
			func(item string, index Accessor[int]) loom.Node {
				child := test.NewMockNode(fmt.Sprintf("item-%s-%d", item, index()))
				children = append(children, child)
				return child
			},
		)
		_, err := loom.Render("parent", forNode)
		assert.NoError(t, err)

		assert.Len(t, children, 3, "should have created three children")
		for i, item := range items() {
			assert.Equal(t, 1, children[i].MountCalls(), "child should be mounted once")
			assert.Equal(t, fmt.Sprintf("item-%s-%d", item, i), children[i].ID(), "child should have correct name")
		}
	})

	t.Run("renders empty list", func(t *testing.T) {
		items, _ := Signal([]string{})

		var children []*test.MockNode
		forNode := For(
			items,
			func(item string, index Accessor[int]) loom.Node {
				child := test.NewMockNode(fmt.Sprintf("item-%s-%d", item, index()))
				children = append(children, child)
				return child
			},
		)
		_, err := loom.Render("parent", forNode)
		assert.NoError(t, err)

		assert.Len(t, children, 0, "should have created no children")
	})

	t.Run("appends new items", func(t *testing.T) {
		items, setItems := Signal([]string{"A", "B", "C"})

		renderCalls := 0
		var children []*test.MockNode
		forNode := For(
			items,
			func(item string, index Accessor[int]) loom.Node {
				renderCalls++
				child := test.NewMockNode(fmt.Sprintf("item-%s-%d", item, index()))
				children = append(children, child)
				return child
			},
		)

		_, err := loom.Render("parent", forNode)
		assert.NoError(t, err)

		assert.Len(t, children, 3, "should have created three children")

		setItems([]string{"A", "B", "C", "D", "E"})

		assert.Len(t, children, 5, "should have created five children")
		assert.Equal(t, 5, renderCalls, "render function should be called for each item")

		for i, item := range items() {
			assert.Equal(t, 1, children[i].MountCalls(), "child should be mounted once")
			assert.Equal(t, 0, children[i].UpdateCalls(), "child should not be updated")
			assert.Equal(t, 0, children[i].UnmountCalls(), "child should not be unmounted")
			assert.Equal(t, fmt.Sprintf("item-%s-%d", item, i), children[i].ID(), "child should have correct name")
		}
	})

	t.Run("prepends new items", func(t *testing.T) {
		items, setItems := Signal([]string{"C", "D", "E"})

		renderCalls := 0
		var children []*test.MockNode
		forNode := For(
			items,
			func(item string, index Accessor[int]) loom.Node {
				renderCalls++
				child := test.NewMockNode(fmt.Sprintf("item-%s-%d", item, index()))
				children = append(children, child)

				OnCleanup(func() {
					children = removeChild(children, child)
				})

				return child
			},
		)

		_, err := loom.Render("parent", forNode)
		assert.NoError(t, err)

		assert.Len(t, children, 3, "should have created three children")

		setItems([]string{"A", "B", "C", "D", "E"})

		assert.Len(t, children, 5, "should have five children")
		assert.Equal(t, 5, renderCalls, "render function should be called for each item")

		assert.Equal(t, "item-C-0", children[0].ID(), "first child should be original C (reused)")
		assert.Equal(t, "item-D-1", children[1].ID(), "second child should be original D (reused)")
		assert.Equal(t, "item-E-2", children[2].ID(), "third child should be original E (reused)")
	})

	t.Run("inserts items in middle", func(t *testing.T) {
		items, setItems := Signal([]string{"A", "B", "D", "E"})

		renderCalls := 0
		var children []*test.MockNode
		forNode := For(
			items,
			func(item string, index Accessor[int]) loom.Node {
				renderCalls++
				child := test.NewMockNode(fmt.Sprintf("item-%s-%d", item, index()))
				children = append(children, child)

				OnCleanup(func() {
					children = removeChild(children, child)
				})

				return child
			},
		)

		_, err := loom.Render("parent", forNode)
		assert.NoError(t, err)

		assert.Len(t, children, 4, "should have created four children")

		setItems([]string{"A", "B", "C", "D", "E"})

		assert.Len(t, children, 5, "should have five children")
		assert.Equal(t, 5, renderCalls, "render function should be called for each item")

		assert.Equal(t, "item-A-0", children[0].ID(), "first child should be original A (reused)")
		assert.Equal(t, "item-B-1", children[1].ID(), "second child should be original B (reused)")
	})

	t.Run("removes items from end", func(t *testing.T) {
		items, setItems := Signal([]string{"A", "B", "C", "D", "E"})

		renderCalls := 0
		cleanupCalls := 0
		var children []*test.MockNode

		forNode := For(
			items,
			func(item string, index Accessor[int]) loom.Node {
				renderCalls++
				child := test.NewMockNode(fmt.Sprintf("item-%s-%d", item, index()))
				children = append(children, child)

				OnCleanup(func() {
					cleanupCalls++
					children = removeChild(children, child)
				})

				return child
			},
		)

		_, err := loom.Render("parent", forNode)
		assert.NoError(t, err)

		assert.Len(t, children, 5, "should have created five children")

		setItems([]string{"A", "B", "C"})

		assert.Len(t, children, 3, "should have three children after removal")
		assert.Equal(t, 2, cleanupCalls, "cleanup should be called for each removed item")
		assert.Equal(t, 5, renderCalls, "render function should be called only for initial items")

		for i, item := range items() {
			assert.Equal(t, 1, children[i].MountCalls(), "child should be mounted once")
			assert.Equal(t, 0, children[i].UnmountCalls(), "child should not be unmounted")
			assert.Equal(t, fmt.Sprintf("item-%s-%d", item, i), children[i].ID(), "child should have correct name")
		}
	})

	t.Run("removes items from beginning", func(t *testing.T) {
		items, setItems := Signal([]string{"A", "B", "C", "D", "E"})

		renderCalls := 0
		cleanupCalls := 0
		var children []*test.MockNode

		forNode := For(
			items,
			func(item string, index Accessor[int]) loom.Node {
				renderCalls++
				child := test.NewMockNode(fmt.Sprintf("item-%s-%d", item, index()))
				children = append(children, child)

				OnCleanup(func() {
					cleanupCalls++
					children = removeChild(children, child)
				})

				return child
			},
		)

		_, err := loom.Render("parent", forNode)
		assert.NoError(t, err)

		assert.Len(t, children, 5, "should have created five children")

		setItems([]string{"C", "D", "E"})

		assert.Len(t, children, 3, "should have three children after removal")
		assert.Equal(t, 2, cleanupCalls, "cleanup should be called for each removed item")
		assert.Equal(t, 5, renderCalls, "render function should be called only for initial items")

		assert.Equal(t, "item-C-2", children[0].ID(), "first child should be original C (reused)")
		assert.Equal(t, "item-D-3", children[1].ID(), "second child should be original D (reused)")
		assert.Equal(t, "item-E-4", children[2].ID(), "third child should be original E (reused)")

		// With the SolidJS-inspired approach, items that move get remounted
		assert.Equal(t, 2, children[0].MountCalls(), "C should be remounted when moved")
		assert.Equal(t, 1, children[0].UnmountCalls(), "C should be unmounted when moved")
		assert.Equal(t, 2, children[1].MountCalls(), "D should be remounted when moved")
		assert.Equal(t, 1, children[1].UnmountCalls(), "D should be unmounted when moved")
		assert.Equal(t, 2, children[2].MountCalls(), "E should be remounted when moved")
		assert.Equal(t, 1, children[2].UnmountCalls(), "E should be unmounted when moved")
	})

	t.Run("removes items from middle", func(t *testing.T) {
		items, setItems := Signal([]string{"A", "B", "C", "D", "E"})

		renderCalls := 0
		cleanupCalls := 0
		var children []*test.MockNode

		forNode := For(
			items,
			func(item string, index Accessor[int]) loom.Node {
				renderCalls++
				child := test.NewMockNode(fmt.Sprintf("item-%s-%d", item, index()))
				children = append(children, child)

				OnCleanup(func() {
					cleanupCalls++
					children = removeChild(children, child)
				})

				return child
			},
		)

		_, err := loom.Render("parent", forNode)
		assert.NoError(t, err)

		assert.Len(t, children, 5, "should have created five children")

		setItems([]string{"A", "C", "E"})

		assert.Len(t, children, 3, "should have three children after removal")
		assert.Equal(t, 2, cleanupCalls, "cleanup should be called for each removed item")
		assert.Equal(t, 5, renderCalls, "render function should be called only for initial items")

		assert.Equal(t, "item-A-0", children[0].ID(), "first child should be original A (reused)")
		assert.Equal(t, "item-C-2", children[1].ID(), "second child should be original C (reused)")
		assert.Equal(t, "item-E-4", children[2].ID(), "third child should be original E (reused)")

		// A is in the common prefix so it stays in place
		assert.Equal(t, 1, children[0].MountCalls(), "A should be mounted once")
		assert.Equal(t, 0, children[0].UnmountCalls(), "A should not be unmounted")

		// C and E are in the window (moved) so they get remounted
		assert.Equal(t, 2, children[1].MountCalls(), "C should be remounted when moved")
		assert.Equal(t, 1, children[1].UnmountCalls(), "C should be unmounted when moved")
		assert.Equal(t, 2, children[2].MountCalls(), "E should be remounted when moved")
		assert.Equal(t, 1, children[2].UnmountCalls(), "E should be unmounted when moved")
	})

	t.Run("clears all items", func(t *testing.T) {
		items, setItems := Signal([]string{"A", "B", "C"})

		cleanupCalls := 0
		var children []*test.MockNode

		forNode := For(
			items,
			func(item string, index Accessor[int]) loom.Node {
				child := test.NewMockNode(fmt.Sprintf("item-%s-%d", item, index()))
				children = append(children, child)

				OnCleanup(func() {
					cleanupCalls++
				})

				return child
			},
		)

		_, err := loom.Render("parent", forNode)
		assert.NoError(t, err)

		assert.Len(t, children, 3, "should have created three children")

		setItems([]string{})

		assert.Equal(t, 3, cleanupCalls, "cleanup should be called for each removed item")

		for _, child := range children {
			assert.Equal(t, 1, child.MountCalls(), "child should be mounted once")
			assert.Equal(t, 1, child.UnmountCalls(), "child should be unmounted once")
		}
	})

	t.Run("handles duplicate values", func(t *testing.T) {
		items, setItems := Signal([]string{"A", "A", "B"})

		renderCalls := 0
		var children []*test.MockNode

		forNode := For(
			items,
			func(item string, index Accessor[int]) loom.Node {
				renderCalls++
				child := test.NewMockNode(fmt.Sprintf("item-%s-%d", item, index()))
				children = append(children, child)
				return child
			},
		)

		_, err := loom.Render("parent", forNode)
		assert.NoError(t, err)

		assert.Len(t, children, 3, "should have created three children")
		assert.Equal(t, "item-A-0", children[0].ID(), "first child should be A-0")
		assert.Equal(t, "item-A-1", children[1].ID(), "second child should be A-1")
		assert.Equal(t, "item-B-2", children[2].ID(), "third child should be B-2")

		setItems([]string{"A", "B", "A"})

		assert.Equal(t, 3, renderCalls, "render function should not be called for reused items")
		assert.Len(t, children, 3, "should have three children")
	})

	t.Run("disposes reactive children", func(t *testing.T) {
		items, setItems := Signal([]string{"A", "B"})

		var children []*test.MockNode
		var setters []func(int)

		forNode := For(
			items,
			func(item string, index Accessor[int]) loom.Node {
				get, set := Signal(0)
				child := test.NewMockNode(fmt.Sprintf("item-%s-%d", item, index()))
				children = append(children, child)
				setters = append(setters, set)

				return Bind(func() loom.Node {
					get()
					return child
				})
			},
		)

		_, err := loom.Render("parent", forNode)
		assert.NoError(t, err)

		assert.Len(t, children, 2, "should have created two children")

		setItems([]string{"B"})

		assert.Equal(t, 1, children[0].MountCalls(), "removed child should be mounted once")
		assert.Equal(t, 0, children[0].UpdateCalls(), "removed child should not be updated after removal")
		assert.Equal(t, 1, children[0].UnmountCalls(), "removed child should be unmounted once after removal")

		setters[0](1)
		assert.Equal(t, 1, children[0].MountCalls(), "removed child should be mounted once")
		assert.Equal(t, 0, children[0].UpdateCalls(), "removed child should not be updated after removal")
		assert.Equal(t, 1, children[0].UnmountCalls(), "removed child should be unmounted once after removal")
	})
}
