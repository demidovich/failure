package failure

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testingWrappedFailure interface {
	Unwrap() error
}

func TestFailure_New(t *testing.T) {
	t.Run("new_default", func(t *testing.T) {
		err := New("foo")
		assert.Equal(t, "foo", err.Error())
	})

	t.Run("new_with_formatted_message", func(t *testing.T) {
		err := New("foo %s", "bar")
		assert.Equal(t, "foo bar", err.Error())
	})

	t.Run("new_stack", func(t *testing.T) {
		stackMode = StackModeFull
		err := New("foo")

		assert.True(t, len(err.Stack()) > 0)
	})
}

func TestFailure_Wrap(t *testing.T) {
	t.Run("wrap_default", func(t *testing.T) {
		errA := errors.New("error A")
		errB := fmt.Errorf("%w", errA)
		errC := Wrap(errB, "error C")

		assert.True(t, errors.Is(errC, errA))
	})

	t.Run("wrap_stack", func(t *testing.T) {
		stackMode = StackModeFull
		err := Wrap(errors.New("foo"), "bar")

		assert.True(t, len(err.Stack()) > 0)
	})

	t.Run("wrap_nil", func(t *testing.T) {
		err := Wrap(nil, "")

		assert.Nil(t, err)
	})

	t.Run("wrap_with_formatted_message", func(t *testing.T) {
		errA := errors.New("error A")
		errB := fmt.Errorf("%w", errA)
		errC := Wrap(errB, "error C")

		assert.True(t, errors.Is(errC, errA))
	})

	t.Run("wrap_format_with_verb_v", func(t *testing.T) {
		stackMode = "none"
		errA := errors.New("error A")
		errB := Wrap(errA, "error B")
		msg := fmt.Sprintf("%v", errB)

		assert.Equal(t, "error B: error A", msg)
	})

	t.Run("wrap_unwrap_cause", func(t *testing.T) {
		errA := errors.New("error A")
		errB := Wrap(errA, "error B")

		var cause error
		if e, ok := errB.(testingWrappedFailure); ok {
			cause = e.Unwrap()
		}

		assert.Equal(t, "error A", cause.Error())
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

	t.Run("wrap_deferred_on_failure_error", func(t *testing.T) {
		var wrapDeferredFailure = func() (err error) {
			defer WrapDeferred(&err, "foo")
			err = New("bar")
			return
		}

		err := wrapDeferredFailure()

		assert.Equal(t, "foo: bar", err.Error())
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

		assert.Nil(t, err)
	})
}
