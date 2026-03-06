package components

import (
	"testing"

	"github.com/AnatoleLucet/loom"
	"github.com/AnatoleLucet/loom/test"
	"github.com/stretchr/testify/assert"
)

func TestProvider(t *testing.T) {
	t.Run("renders children", func(t *testing.T) {
		_, ctx := NewContext("default")
		child := test.NewMockNode("child")

		provider := ctx.Provider("provided", func() loom.Node {
			return child
		})
		_, err := loom.Render("parent", provider)
		assert.NoError(t, err)

		assert.Equal(t, 1, child.MountCalls(), "child should be mounted once")
	})

	t.Run("provides context value to children", func(t *testing.T) {
		read, ctx := NewContext("default")

		reads := []string{}
		child := test.NewMockNode("child")
		child.OnMount(func() {
			reads = append(reads, read())
		})

		provider := ctx.Provider("provided", func() loom.Node {
			reads = append(reads, read())
			return child
		})
		_, err := loom.Render("parent", provider)
		assert.NoError(t, err)

		assert.Equal(t, []string{"provided", "provided"}, reads, "context value should be provided to children")
	})

	t.Run("nested providers", func(t *testing.T) {
		read, ctx := NewContext("default")

		reads := []string{}
		child := test.NewMockNode("child")
		child.OnMount(func() {
			reads = append(reads, read())
		})

		provider := ctx.Provider("outer", func() loom.Node {
			reads = append(reads, read())
			return ctx.Provider("inner", func() loom.Node {
				reads = append(reads, read())
				return child
			})
		})
		_, err := loom.Render("parent", provider)
		assert.NoError(t, err)

		assert.Equal(t, []string{"outer", "inner", "inner"}, reads, "nested providers should override context value")
	})

	t.Run("default context value", func(t *testing.T) {
		read, ctx := NewContext("default")

		reads := []string{}
		child := test.NewMockNode("child")
		child.OnMount(func() {
			reads = append(reads, read())
		})

		provider := ctx.Provider("provided", func() loom.Node {
			reads = append(reads, read())
			return child
		})

		container := test.NewMockNode("container", provider)
		container.OnMount(func() {
			reads = append(reads, read())
		})

		_, err := loom.Render("parent", container)
		assert.NoError(t, err)

		assert.Equal(t, []string{"default", "provided", "provided"}, reads, "default context value should be used outside provider")
	})

	t.Run("from accessor", func(t *testing.T) {
		read, ctx := NewContext("default")
		value, setValue := Signal("initial")

		reads := []string{}
		child := test.NewMockNode("child")
		child.OnMount(func() {
			reads = append(reads, read())
		})
		child.OnUpdate(func() {
			reads = append(reads, read())
		})

		provider := ctx.BindProvider(value, func() loom.Node {
			reads = append(reads, read())
			return Bind(func() loom.Node { return child }) // make child subscribe to read()
		})
		_, err := loom.Render("parent", provider)
		assert.NoError(t, err)

		assert.Equal(t, []string{"initial", "initial"}, reads, "context value should be provided to children from accessor")

		setValue("updated")
		assert.Equal(t, []string{"initial", "initial", "updated"}, reads, "context value should update in children when accessor changes")
	})

	t.Run("in Bind", func(t *testing.T) {
		read, ctx := NewContext("default")
		value, setValue := Signal("initial")

		reads := []string{}
		child := test.NewMockNode("child")
		child.OnMount(func() {
			reads = append(reads, read())
		})
		child.OnUpdate(func() {
			reads = append(reads, read())
		})

		// this test mainly exists to ensure the new Owner creation in ProviderBind&Provider doesn't
		// mess the context when recreated within a Bind
		provider := Bind(func() loom.Node {
			reads = append(reads, read())
			return ctx.BindProvider(value, func() loom.Node {
				reads = append(reads, read())
				return child
			})
		})

		_, err := loom.Render("parent", provider)
		assert.NoError(t, err)

		assert.Equal(t, []string{"default", "initial", "initial"}, reads, "context value should be provided to children from accessor within Bind")

		setValue("updated")
		assert.Equal(t, []string{"default", "initial", "initial", "default", "updated", "updated"}, reads, "context value should update in children when accessor changes within Bind")
	})
}
