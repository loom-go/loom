package components

import (
	"context"
	"sync"
	"testing"

	"github.com/loom-go/loom"
	"github.com/stretchr/testify/assert"
)

func TestSelf(t *testing.T) {
	t.Run("gives the current Component", func(t *testing.T) {
		self := Self()

		assert.NotNil(t, self, "Self should not be nil")
		assert.False(t, self.IsDisposed(), "Self should not be disposed")
	})

	t.Run("marks as disposed", func(t *testing.T) {
		visible, setVisible := Signal(true)

		var self loom.Component
		show := Show(visible, func() loom.Node {
			self = Self()
			return nil
		})
		_, err := loom.Render("parent", show)
		assert.NoError(t, err)

		assert.NotNil(t, self, "Self should not be nil")
		assert.False(t, self.IsDisposed(), "Self should not be disposed")

		setVisible(false)
		assert.True(t, self.IsDisposed(), "Self should be disposed after visibility is set to false")
	})

	t.Run("closes disposed channel", func(t *testing.T) {
		var wg sync.WaitGroup

		visible, setVisible := Signal(true)

		cleaned := false
		show := Show(visible, func() loom.Node {
			wg.Add(1)
			go func(self loom.Component) {
				defer wg.Done()

				<-self.Disposed()
				cleaned = true
			}(Self())

			return nil
		})
		_, err := loom.Render("parent", show)
		assert.NoError(t, err)

		assert.False(t, cleaned, "Cleaned should be false before visibility is set to false")

		setVisible(false)
		wg.Wait()
		assert.True(t, cleaned, "Cleaned should be true after visibility is set to false")
	})

	t.Run("context gets cancelled on dispose", func(t *testing.T) {
		var wg sync.WaitGroup

		visible, setVisible := Signal(true)

		cleaned := false
		show := Show(visible, func() loom.Node {
			wg.Add(1)
			go func(ctx context.Context) {
				defer wg.Done()

				<-ctx.Done()
				cleaned = true
			}(Self().Context())

			return nil
		})
		_, err := loom.Render("parent", show)
		assert.NoError(t, err)

		assert.False(t, cleaned, "Cleaned should be false before visibility is set to false")

		setVisible(false)
		wg.Wait()
		assert.True(t, cleaned, "Cleaned should be true after visibility is set to false")
	})
}
