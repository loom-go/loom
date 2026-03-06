package components

import (
	"testing"

	"github.com/AnatoleLucet/loom"
	"github.com/stretchr/testify/assert"
)

func TestRef(t *testing.T) {
	type Element struct {
		id string
	}

	t.Run("assigns the parent to ptr", func(t *testing.T) {
		var ref Element

		parent := Element{"parent"}
		_, err := loom.Render(parent, Apply(Ref{Ptr: &ref}))
		assert.NoError(t, err)

		assert.NotNil(t, ref, "ref should not be nil")
		assert.Equal(t, parent, ref, "ref should be equal to the parent")
	})

	t.Run("errors if ptr type does not match the parent type", func(t *testing.T) {
		var ref string

		parent := Element{"parent"}
		_, err := loom.Render(parent, Apply(Ref{Ptr: &ref}))
		assert.Error(t, err)
	})

	t.Run("calls fn with the parent", func(t *testing.T) {
		var calledWith Element
		onRef := func(p Element) {
			calledWith = p
		}

		parent := Element{"parent"}
		_, err := loom.Render(parent, Apply(Ref{Fn: onRef}))
		assert.NoError(t, err)

		assert.NotNil(t, calledWith, "calledWith should not be nil")
		assert.Equal(t, parent, calledWith, "calledWith should be equal to the parent")
	})

	t.Run("errors if fn type does not match the parent type", func(t *testing.T) {
		onRef := func(p string) {}

		parent := Element{"parent"}
		_, err := loom.Render(parent, Apply(Ref{Fn: onRef}))
		assert.Error(t, err)
	})
}
