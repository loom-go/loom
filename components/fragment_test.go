package components

import (
	"testing"

	"github.com/AnatoleLucet/loom"
	"github.com/AnatoleLucet/loom/test"
	"github.com/stretchr/testify/assert"
)

func TestFragment(t *testing.T) {
	t.Run("renders children", func(t *testing.T) {
		child1 := test.NewMockNode("child1")
		child2 := test.NewMockNode("child2")

		fragment := Fragment(child1, child2)
		_, err := loom.Render("parent", fragment)
		assert.NoError(t, err)

		assert.Equal(t, 1, child1.MountCalls(), "child1 should be mounted once")
		assert.Equal(t, 1, child2.MountCalls(), "child2 should be mounted once")
	})
}
