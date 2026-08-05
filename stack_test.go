package failure

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:noinline
func testingExtractStackMock() error {
	return New("foo")
}

func TestStack_ExtractStack(t *testing.T) {
	t.Run("extract_from_failure_error", func(t *testing.T) {
		stackMode = StackModeFull
		err := testingExtractStackMock()

		stack, ok := ExtractStack(err)

		assert.True(t, ok)
		assert.NotEmpty(t, stack.Slice())
		assert.Contains(t, stack.Slice()[0], "testingExtractStackMock")
	})

	t.Run("extract_from_failure_wrapped_error", func(t *testing.T) {
		stackMode = StackModeFull
		err := testingExtractStackMock()
		err = Wrap(err, "")

		stack, ok := ExtractStack(err)

		assert.True(t, ok)
		assert.NotEmpty(t, stack.Slice())
		assert.Contains(t, stack.Slice()[0], "testingExtractStackMock")
	})

	t.Run("extract_from_deep_wrapped_error", func(t *testing.T) {
		stackMode = StackModeFull
		err := testingExtractStackMock()
		err = fmt.Errorf("%w", err)
		err = Wrap(err, "")
		err = fmt.Errorf("%w", err)

		stack, ok := ExtractStack(err)

		assert.True(t, ok)
		assert.NotEmpty(t, stack)
		assert.Contains(t, stack.Slice()[0], "testingExtractStackMock")
	})
}

func TestStack(t *testing.T) {
	t.Run("frames", func(t *testing.T) {
		stackMode = StackModeFull
		stack := newStack()

		assert.NotNil(t, stack.Frames())
	})

	t.Run("frames_on_nil_stack", func(t *testing.T) {
		var stack *Stack

		assert.Nil(t, stack.Frames())
	})

	t.Run("slice", func(t *testing.T) {
		stackMode = StackModeFull
		stack := newStack()

		assert.NotEmpty(t, stack.Slice())
	})

	t.Run("slice_on_nil_stack", func(t *testing.T) {
		var stack *Stack

		assert.Empty(t, stack.Slice())
	})

	t.Run("string", func(t *testing.T) {
		stackMode = StackModeFull
		stack := newStack()

		assert.NotEmpty(t, stack.String())
	})

	t.Run("string_on_nil_stack", func(t *testing.T) {
		var stack *Stack

		assert.Empty(t, stack.String())
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
