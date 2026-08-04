package failure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	t.Run("frames", func(t *testing.T) {
		stackMode = StackModeFull
		stack := newStack()

		assert.NotNil(t, stack.Frames())
	})

	t.Run("slice", func(t *testing.T) {
		stackMode = StackModeFull
		stack := newStack()

		assert.True(t, len(stack.Slice()) > 0)
	})

	t.Run("string", func(t *testing.T) {
		stackMode = StackModeFull
		stack := newStack()

		assert.NotEqual(t, "", stack.String())
	})

	t.Run("slice_cache", func(t *testing.T) {
		stackMode = StackModeFull
		stack := newStack()

		assert.False(t, stack.hasSlice)
		assert.NotEmpty(t, stack.Slice())

		assert.True(t, stack.hasSlice)
		assert.NotEmpty(t, stack.Slice())
	})
}
