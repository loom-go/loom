package components

import (
	"fmt"
	"slices"
	"testing"

	"github.com/AnatoleLucet/loom"
	"github.com/AnatoleLucet/loom/signals"
	"github.com/AnatoleLucet/loom/test"
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
		items, _ := signals.Signal([]string{"A", "B", "C"})

		var children []*test.MockNode
		forNode := For(items, func(item signals.Accessor[string], index signals.Accessor[int]) loom.Node {
			child := test.NewMockNode(fmt.Sprintf("item-%s-%d", item(), index()))
			children = append(children, child)
			return child
		})
		err := loom.Render("parent", forNode)
		assert.NoError(t, err)

		assert.Len(t, children, 3, "should have created three children")
		for i, item := range items() {
			assert.Equal(t, 1, children[i].MountCalls(), "child should be mounted once")
			assert.Equal(t, fmt.Sprintf("item-%s-%d", item, i), children[i].ID(), "child should have correct name")
		}
	})

	t.Run("renders empty list", func(t *testing.T) {
		items, _ := signals.Signal([]string{})

		var children []*test.MockNode
		forNode := For(items, func(item signals.Accessor[string], index signals.Accessor[int]) loom.Node {
			child := test.NewMockNode(fmt.Sprintf("item-%s-%d", item(), index()))
			children = append(children, child)
			return child
		})
		err := loom.Render("parent", forNode)
		assert.NoError(t, err)

		assert.Len(t, children, 0, "should have created no children")
	})

	t.Run("appends new items", func(t *testing.T) {
		items, setItems := signals.Signal([]string{"A", "B", "C"})

		renderCalls := 0
		var children []*test.MockNode
		var itemAccessors []signals.Accessor[string]
		var indexAccessors []signals.Accessor[int]

		forNode := For(items, func(item signals.Accessor[string], index signals.Accessor[int]) loom.Node {
			renderCalls++

			child := test.NewMockNode(fmt.Sprintf("item-%s-%d", item(), index()))
			children = append(children, child)

			itemAccessors = append(itemAccessors, item)
			indexAccessors = append(indexAccessors, index)

			return child
		})

		err := loom.Render("parent", forNode)
		assert.NoError(t, err)

		assert.Len(t, children, 3, "should have created three children")

		setItems([]string{"A", "B", "C", "D", "E"})

		assert.Len(t, children, 5, "should have created five children")
		assert.Equal(t, 5, renderCalls, "render function should be called for each new item")

		for i, item := range items() {
			assert.Equal(t, 1, children[i].MountCalls(), "child should be mounted once")
			assert.Equal(t, fmt.Sprintf("item-%s-%d", item, i), children[i].ID(), "child should have correct name")
			assert.Equal(t, item, itemAccessors[i](), "item accessor should return correct item")
			assert.Equal(t, i, indexAccessors[i](), "index accessor should return correct index")
		}
	})

	t.Run("prepends new items", func(t *testing.T) {
		items, setItems := signals.Signal([]string{"C", "D", "E"})

		renderCalls := 0
		var children []*test.MockNode
		var itemAccessors []signals.Accessor[string]
		var indexAccessors []signals.Accessor[int]

		forNode := For(items, func(item signals.Accessor[string], index signals.Accessor[int]) loom.Node {
			renderCalls++

			child := test.NewMockNode(fmt.Sprintf("item-%s-%d", item(), index()))
			children = append(children, child)

			itemAccessors = append(itemAccessors, item)
			indexAccessors = append(indexAccessors, index)
			return child
		})

		err := loom.Render("parent", forNode)
		assert.NoError(t, err)

		assert.Len(t, children, 3, "should have created three children")

		setItems([]string{"A", "B", "C", "D", "E"})

		assert.Len(t, children, 5, "should have created five children")
		assert.Equal(t, 5, renderCalls, "render function should be called for each new item")

		assert.Equal(t, "item-C-0", children[0].ID(), "first child should be C-0")  // first child should still be C after prepending. only the accessors are updated
		assert.Equal(t, "item-D-1", children[1].ID(), "second child should be D-1") // second child should still be D after prepending. only the accessors are updated
		assert.Equal(t, "item-E-2", children[2].ID(), "third child should be E-2")
		assert.Equal(t, "item-D-3", children[3].ID(), "fourth child should be D-3")
		assert.Equal(t, "item-E-4", children[4].ID(), "fifth child should be E-4")

		for i, item := range items() {
			assert.Equal(t, 1, children[i].MountCalls(), "child should be mounted once")
			assert.Equal(t, item, itemAccessors[i](), "item accessor should return correct item")
			assert.Equal(t, i, indexAccessors[i](), "index accessor should return correct index")
		}
	})

	t.Run("inserts items in middle", func(t *testing.T) {
		items, setItems := signals.Signal([]string{"A", "B", "D", "E"})

		renderCalls := 0
		var children []*test.MockNode
		var itemAccessors []signals.Accessor[string]
		var indexAccessors []signals.Accessor[int]
		forNode := For(items, func(item signals.Accessor[string], index signals.Accessor[int]) loom.Node {
			renderCalls++

			child := test.NewMockNode(fmt.Sprintf("item-%s-%d", item(), index()))
			children = append(children, child)

			itemAccessors = append(itemAccessors, item)
			indexAccessors = append(indexAccessors, index)
			return child
		})

		err := loom.Render("parent", forNode)
		assert.NoError(t, err)

		assert.Len(t, children, 4, "should have created four children")

		setItems([]string{"A", "B", "C", "D", "E"})

		assert.Len(t, children, 5, "should have created five children")
		assert.Equal(t, 5, renderCalls, "render function should be called for each new item")

		assert.Equal(t, "item-A-0", children[0].ID(), "first child should be A-0")
		assert.Equal(t, "item-B-1", children[1].ID(), "second child should be B-1")
		assert.Equal(t, "item-D-2", children[2].ID(), "third child should be D-2") // third child should still be D after inserting C. only the accessors are updated
		assert.Equal(t, "item-E-3", children[3].ID(), "fourth child should be E-3")
		assert.Equal(t, "item-E-4", children[4].ID(), "fifth child should be E-4") // fifth child should still be E after inserting C. only the accessors are updated

		for i, item := range items() {
			assert.Equal(t, 1, children[i].MountCalls(), "child should be mounted once")
			assert.Equal(t, item, itemAccessors[i](), "item accessor should return correct item")
			assert.Equal(t, i, indexAccessors[i](), "index accessor should return correct index")
		}
	})

	t.Run("removes items from end", func(t *testing.T) {
		items, setItems := signals.Signal([]string{"A", "B", "C", "D", "E"})

		renderCalls := 0
		cleanupCalls := 0
		var children []*test.MockNode
		var itemAccessors []signals.Accessor[string]
		var indexAccessors []signals.Accessor[int]

		forNode := For(items, func(item signals.Accessor[string], index signals.Accessor[int]) loom.Node {
			renderCalls++

			child := test.NewMockNode(fmt.Sprintf("item-%s-%d", item(), index()))
			children = append(children, child)

			itemAccessors = append(itemAccessors, item)
			indexAccessors = append(indexAccessors, index)

			signals.OnCleanup(func() {
				cleanupCalls++
				children = removeChild(children, child)
			})

			return child
		})

		err := loom.Render("parent", forNode)
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
			assert.Equal(t, item, itemAccessors[i](), "item accessor should return correct item")
			assert.Equal(t, i, indexAccessors[i](), "index accessor should return correct index")
		}
	})

	t.Run("removes items from beginning", func(t *testing.T) {
		items, setItems := signals.Signal([]string{"A", "B", "C", "D", "E"})

		renderCalls := 0
		cleanupCalls := 0
		var children []*test.MockNode
		var itemAccessors []signals.Accessor[string]
		var indexAccessors []signals.Accessor[int]

		forNode := For(items, func(item signals.Accessor[string], index signals.Accessor[int]) loom.Node {
			renderCalls++

			child := test.NewMockNode(fmt.Sprintf("item-%s-%d", item(), index()))
			children = append(children, child)
			itemAccessors = append(itemAccessors, item)
			indexAccessors = append(indexAccessors, index)

			signals.OnCleanup(func() {
				cleanupCalls++
				children = removeChild(children, child)
			})

			return child
		})

		err := loom.Render("parent", forNode)
		assert.NoError(t, err)

		assert.Len(t, children, 5, "should have created five children")
		setItems([]string{"C", "D", "E"})

		assert.Len(t, children, 3, "should have three children after removal")
		assert.Equal(t, 2, cleanupCalls, "cleanup should be called for each removed item")
		assert.Equal(t, 5, renderCalls, "render function should be called only for initial items")

		assert.Equal(t, "item-A-0", children[0].ID(), "first child should be A-0")  // first child should still be A after removing B and C. only the accessors are updated
		assert.Equal(t, "item-B-1", children[1].ID(), "second child should be B-1") // second child should still be B after removing B and C. only the accessors are updated
		assert.Equal(t, "item-C-2", children[2].ID(), "third child should be C-2")  // third child should still be C after removing B and C. only the accessors are updated

		for i, item := range items() {
			assert.Equal(t, 1, children[i].MountCalls(), "child should be mounted once")
			assert.Equal(t, 0, children[i].UnmountCalls(), "child should not be unmounted")
			assert.Equal(t, item, itemAccessors[i](), "item accessor should return correct item")
			assert.Equal(t, i, indexAccessors[i](), "index accessor should return correct index")
		}
	})

	t.Run("removes items from middle", func(t *testing.T) {
		items, setItems := signals.Signal([]string{"A", "B", "C", "D", "E"})

		renderCalls := 0
		cleanupCalls := 0
		var children []*test.MockNode
		var itemAccessors []signals.Accessor[string]
		var indexAccessors []signals.Accessor[int]

		forNode := For(items, func(item signals.Accessor[string], index signals.Accessor[int]) loom.Node {
			renderCalls++

			child := test.NewMockNode(fmt.Sprintf("item-%s-%d", item(), index()))
			children = append(children, child)
			itemAccessors = append(itemAccessors, item)
			indexAccessors = append(indexAccessors, index)

			signals.OnCleanup(func() {
				cleanupCalls++
				children = removeChild(children, child)
			})

			return child
		})

		err := loom.Render("parent", forNode)
		assert.NoError(t, err)

		assert.Len(t, children, 5, "should have created five children")

		setItems([]string{"A", "C", "E"})

		assert.Len(t, children, 3, "should have three children after removal")
		assert.Equal(t, 2, cleanupCalls, "cleanup should be called for each removed item")
		assert.Equal(t, 5, renderCalls, "render function should be called only for initial items")

		assert.Equal(t, "item-A-0", children[0].ID(), "first child should be A-0")
		assert.Equal(t, "item-B-1", children[1].ID(), "second child should be B-1")
		assert.Equal(t, "item-C-2", children[2].ID(), "third child should be C-2")

		for i, item := range items() {
			assert.Equal(t, 1, children[i].MountCalls(), "child should be mounted once")
			assert.Equal(t, 0, children[i].UnmountCalls(), "child should not be unmounted")
			assert.Equal(t, item, itemAccessors[i](), "item accessor should return correct item")
			assert.Equal(t, i, indexAccessors[i](), "index accessor should return correct index")
		}
	})

	t.Run("clears all items", func(t *testing.T) {
		items, setItems := signals.Signal([]string{"A", "B", "C"})

		cleanupCalls := 0
		var children []*test.MockNode

		forNode := For(items, func(item signals.Accessor[string], index signals.Accessor[int]) loom.Node {
			child := test.NewMockNode(fmt.Sprintf("item-%s-%d", item(), index()))
			children = append(children, child)

			signals.OnCleanup(func() {
				cleanupCalls++
			})

			return child
		})

		err := loom.Render("parent", forNode)
		assert.NoError(t, err)

		assert.Len(t, children, 3, "should have created three children")

		setItems([]string{})

		assert.Equal(t, 3, cleanupCalls, "cleanup should be called for each removed item")

		for _, child := range children {
			assert.Equal(t, 1, child.MountCalls(), "child should be mounted once")
			assert.Equal(t, 1, child.UnmountCalls(), "child should be unmounted once")
		}
	})

	t.Run("keyer compares by key not value", func(t *testing.T) {
		type uncomparable struct {
			name string
			fn   func()
		}

		items, setItems := signals.Signal([]uncomparable{
			{name: "A", fn: func() {}},
			{name: "B", fn: func() {}},
			{name: "C", fn: func() {}},
		})

		renderCalls := 0
		cleanupCalls := 0
		var children []*test.MockNode
		var itemAccessors []signals.Accessor[uncomparable]
		var indexAccessors []signals.Accessor[int]

		forNode := For(
			items,
			func(item uncomparable) any { return item.name },
			func(item signals.Accessor[uncomparable], index signals.Accessor[int]) loom.Node {
				renderCalls++

				child := test.NewMockNode(fmt.Sprintf("item-%s-%d", item().name, index()))
				children = append(children, child)

				itemAccessors = append(itemAccessors, item)
				indexAccessors = append(indexAccessors, index)

				signals.OnCleanup(func() {
					cleanupCalls++
					children = removeChild(children, child)
				})

				return child
			})

		err := loom.Render("parent", forNode)
		assert.NoError(t, err)

		assert.Len(t, children, 3, "should have created three children")

		setItems([]uncomparable{
			{name: "B", fn: func() {}},
			{name: "C", fn: func() {}},
			{name: "A", fn: func() {}},
		})

		assert.Len(t, children, 3, "should still have three children")
		assert.Equal(t, 3, renderCalls, "render function should not be called for existing items")
		assert.Equal(t, 0, cleanupCalls, "cleanup should not be called for existing items")

		assert.Equal(t, "item-A-0", children[0].ID(), "first child should be A-0")
		assert.Equal(t, "item-B-1", children[1].ID(), "second child should be B-1")
		assert.Equal(t, "item-C-2", children[2].ID(), "third child should be C-2")

		for i, item := range items() {
			assert.Equal(t, 1, children[i].MountCalls(), "child should be mounted once")
			assert.Equal(t, 0, children[i].UnmountCalls(), "child should not be unmounted")
			assert.Equal(t, item.name, itemAccessors[i]().name, "item accessor should return correct item")
			assert.Equal(t, i, indexAccessors[i](), "index accessor should return correct index")
		}

		setItems([]uncomparable{
			{name: "A", fn: func() {}},
			{name: "B", fn: func() {}},
		})

		assert.Len(t, children, 2, "should have two children after removal")
		assert.Equal(t, 3, renderCalls, "render function should not be called for existing items")
		assert.Equal(t, 1, cleanupCalls, "cleanup should be called for removed item")

		assert.Equal(t, "item-A-0", children[0].ID(), "first child should be A-0")
		assert.Equal(t, "item-B-1", children[1].ID(), "second child should be B-1")

		for i, item := range items() {
			assert.Equal(t, 1, children[i].MountCalls(), "child should be mounted once")
			assert.Equal(t, 0, children[i].UnmountCalls(), "child should not be unmounted")
			assert.Equal(t, item.name, itemAccessors[i]().name, "item accessor should return correct item")
			assert.Equal(t, i, indexAccessors[i](), "index accessor should return correct index")
		}
	})

	t.Run("panics on mapper as keyer with extra mapper", func(t *testing.T) {
		items, _ := signals.Signal([]string{"A"})

		assert.Panics(t, func() {
			For(
				items,
				func(item signals.Accessor[string], index signals.Accessor[int]) loom.Node {
					return test.NewMockNode("item1")
				},
				func(item signals.Accessor[string], index signals.Accessor[int]) loom.Node {
					return test.NewMockNode("item2")
				},
			)
		})
	})

	t.Run("panics on keyer without mapper", func(t *testing.T) {
		items, _ := signals.Signal([]string{"A"})

		assert.Panics(t, func() {
			For(items, func(s string) any { return s })
		})
	})
}
