package components

import (
	"testing"

	"github.com/AnatoleLucet/loom"
	"github.com/AnatoleLucet/loom/test"
	"github.com/stretchr/testify/assert"
)

func TestOwn(t *testing.T) {
	t.Run("renders children", func(t *testing.T) {
		owner := NewOwner()
		child := test.NewMockNode("child")

		own := Own(owner, Fragment(Fragment(), child, Fragment()))
		_, err := loom.Render("parent", own)
		assert.NoError(t, err)

		assert.Equal(t, 1, child.MountCalls(), "child should be mounted once")
	})

	t.Run("should owns children", func(t *testing.T) {
		owner := NewOwner()
		child := test.NewMockNode("child")

		cleanupCalls := 0
		child.OnMount(func() {
			OnCleanup(func() { cleanupCalls++ })
		})

		own := Own(owner, child)
		_, err := loom.Render("parent", own)
		assert.NoError(t, err)

		assert.Equal(t, 1, child.MountCalls(), "child should be mounted once")

		owner.Dispose()
		assert.Equal(t, 1, cleanupCalls, "child's cleanup should be called once")
	})
}
