package failure

import (
	"fmt"
	"strings"
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
		assert.NotNil(t, stack)
		assert.Contains(t, stack.Frames()[0].Function, "testingExtractStackMock")
	})

	t.Run("extract_from_failure_wrapped_error", func(t *testing.T) {
		stackMode = StackModeFull
		err := testingExtractStackMock()
		err = Wrap(err, "")

		stack, ok := ExtractStack(err)

		assert.True(t, ok)
		assert.NotNil(t, stack)
		assert.Contains(t, stack.Frames()[0].Function, "testingExtractStackMock")
	})

	t.Run("extract_from_deep_wrapped_error", func(t *testing.T) {
		stackMode = StackModeFull
		err := testingExtractStackMock()
		err = fmt.Errorf("%w", err)
		err = Wrap(err, "")
		err = fmt.Errorf("%w", err)

		stack, ok := ExtractStack(err)

		assert.True(t, ok)
		assert.NotNil(t, stack)
		assert.Contains(t, stack.Frames()[0].Function, "testingExtractStackMock")
	})
}

func TestStack(t *testing.T) {
	t.Run("frames", func(t *testing.T) {
		stackMode = StackModeFull
		stack := newStack()

		assert.NotEmpty(t, stack.Frames())
	})

	t.Run("frames_on_nil_stack", func(t *testing.T) {
		var stack *Stack

		assert.Empty(t, stack.Frames())
	})

	t.Run("frames_formatted", func(t *testing.T) {
		stackMode = StackModeFull
		stack := newStack()

		assert.NotEmpty(t, stack.FramesFormatted())
	})

	t.Run("frames_formatted_on_nil_stack", func(t *testing.T) {
		var stack *Stack

		assert.Empty(t, stack.FramesFormatted())
	})

	t.Run("serialize", func(t *testing.T) {
		stackMode = StackModeFull
		stack := newStack()
		serialized := stack.Serialize(", ")
		splitted := strings.Split(serialized, ", ")

		fmt.Println(serialized)

		assert.Len(t, stack.Frames(), len(splitted))
	})

	t.Run("serialize_on_nil_stack", func(t *testing.T) {
		var stack *Stack

		assert.Empty(t, stack.Serialize(", "))
	})
}
