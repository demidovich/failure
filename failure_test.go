package failure

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFailure_New(t *testing.T) {
	t.Run("new_default", func(t *testing.T) {
		err := New("foo")
		assert.Equal(t, "foo", err.Error())
	})

	t.Run("new_formatted_message", func(t *testing.T) {
		err := New("foo: %s", "bar")
		assert.Equal(t, "foo: bar", err.Error())
	})

	t.Run("new_stack", func(t *testing.T) {
		stackMode = StackModeFull
		err := New("foo")
		stack, ok := ExtractStack(err)

		assert.True(t, ok)
		assert.NotEmpty(t, stack.Frames())
	})
}

func TestFailure_Wrap(t *testing.T) {
	t.Run("wrap_default", func(t *testing.T) {
		errA := errors.New("error A")
		errB := fmt.Errorf("%w", errA)
		errC := Wrap(errB, "error C")

		assert.ErrorIs(t, errC, errA)
	})

	t.Run("wrap_formatted_message", func(t *testing.T) {
		errA := errors.New("baz")
		errB := Wrap(errA, "foo: %s", "bar")

		assert.Equal(t, "foo: bar: baz", errB.Error())
	})

	t.Run("wrap_stack", func(t *testing.T) {
		stackMode = StackModeFull
		err := Wrap(errors.New("foo"), "bar")
		stack, ok := ExtractStack(err)

		assert.True(t, ok)
		assert.NotEmpty(t, stack.Frames())
	})

	t.Run("wrap_nil", func(t *testing.T) {
		err := Wrap(nil, "")

		assert.NoError(t, err)
	})

	t.Run("wrap_with_formatted_message", func(t *testing.T) {
		errA := errors.New("error A")
		errB := fmt.Errorf("%w", errA)
		errC := Wrap(errB, "error C")

		assert.ErrorIs(t, errC, errA)
	})

	t.Run("wrap_format_with_verb_v", func(t *testing.T) {
		stackMode = StackModeNone
		errA := errors.New("error A")
		errB := Wrap(errA, "error B")
		msg := fmt.Sprintf("%v", errB)

		assert.Equal(t, "Error: error B: error A", msg)
	})

	t.Run("wrap_unwrap_cause", func(t *testing.T) {
		errA := errors.New("error A")
		errB := Wrap(errA, "error B")
		wrap, ok := errB.(*wrappedFailure) //nolint:errorlint

		assert.True(t, ok)
		assert.Equal(t, "error A", wrap.Unwrap().Error())
	})

	t.Run("wrap_cause_with_empty_message", func(t *testing.T) {
		errA := errors.New("")
		errB := Wrap(errA, "error B")

		assert.Equal(t, "error B", errB.Error())
	})
}

func TestFailure_WrapDeferred(t *testing.T) {
	t.Run("wrap_deferred_on_stderror", func(t *testing.T) {
		var wrapDeferredStderr = func() (err error) {
			defer WrapDeferred(&err, "foo")
			err = errors.New("bar")
			return
		}

		err := wrapDeferredStderr()

		assert.Equal(t, "foo: bar", err.Error())
	})

	t.Run("wrap_deferred_with_formatted_message", func(t *testing.T) {
		var wrapDeferredFailure = func() (err error) {
			defer WrapDeferred(&err, "foo %s", "123")
			err = New("bar")
			return
		}

		err := wrapDeferredFailure()

		assert.Equal(t, "foo 123: bar", err.Error())
	})

	t.Run("wrap_deferred_on_wrapped_failure", func(t *testing.T) {
		var wrapDeferredFailureWrapped = func() (err error) {
			defer WrapDeferred(&err, "baz")
			err = New("bar")
			err = Wrap(err, "foo")
			return
		}

		err := wrapDeferredFailureWrapped()

		assert.Equal(t, "baz: foo: bar", err.Error())
	})

	t.Run("wrap_deferred_on_nil", func(t *testing.T) {
		var wrapDeferredNil = func() (err error) {
			defer WrapDeferred(&err, "foo")
			return
		}

		err := wrapDeferredNil()

		assert.NoError(t, err)
	})
}
