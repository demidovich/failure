package failure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	t.Run("SliceCache", func(t *testing.T) {
		stackMode = StackModeFull
		stack := newStack()

		assert.False(t, stack.hasSlice)
		assert.NotEmpty(t, stack.Slice())

		assert.True(t, stack.hasSlice)
		assert.NotEmpty(t, stack.Slice())
	})
}
