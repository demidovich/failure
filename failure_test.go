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

func TestFailure(t *testing.T) {

	t.Run("ok New", func(t *testing.T) {
		err := New("foo")
		assert.Equal(t, "foo", err.Error())
	})

	t.Run("ok New with formatted message", func(t *testing.T) {
		err := New("foo %s", "bar")
		assert.Equal(t, "foo bar", err.Error())
	})

	t.Run("ok Wrap", func(t *testing.T) {
		errA := errors.New("error A")
		errB := fmt.Errorf("%w", errA)
		errC := Wrap(errB, "error C")

		assert.True(t, errors.Is(errC, errA))
	})

	t.Run("ok Wrap nil", func(t *testing.T) {
		err := Wrap(nil, "")

		assert.Nil(t, err)
	})

	t.Run("ok Wrap with formatted message", func(t *testing.T) {
		errA := errors.New("error A")
		errB := fmt.Errorf("%w", errA)
		errC := Wrap(errB, "error C")

		assert.True(t, errors.Is(errC, errA))
	})

	t.Run("ok Format with verb v", func(t *testing.T) {
		stackMode = "none"
		errA := errors.New("error A")
		errB := Wrap(errA, "error B")
		msg := fmt.Sprintf("%v", errB)

		assert.Equal(t, "error B: error A", msg)
	})

	t.Run("ok Wrap Cause", func(t *testing.T) {
		errA := errors.New("error A")
		errB := Wrap(errA, "error B")

		var cause error
		if e, ok := errB.(testingWrappedFailure); ok {
			cause = e.Unwrap()
		}

		assert.Equal(t, "error A", cause.Error())
	})

	t.Run("ok WrapCause without message", func(t *testing.T) {
		errA := errors.New("")
		errB := Wrap(errA, "error B")

		assert.Equal(t, "error B", errB.Error())
	})

	t.Run("ok New stack created", func(t *testing.T) {
		stackMode = StackModeFull
		err := New("foo")

		assert.True(t, len(err.Stack()) > 0)
	})

	t.Run("ok Wrap stack created", func(t *testing.T) {
		stackMode = StackModeFull
		err := Wrap(errors.New("foo"), "bar")

		assert.True(t, len(err.Stack()) > 0)
	})

	var wrapDeferredStderr = func() (err error) {
		defer WrapDeferred(&err, "foo")
		err = errors.New("bar")
		return
	}

	t.Run("ok WrapDeferred wrap stderr", func(t *testing.T) {
		err := wrapDeferredStderr()
		assert.Equal(t, "foo: bar", err.Error())
	})

	var wrapDeferredFailure = func() (err error) {
		defer WrapDeferred(&err, "foo")
		err = New("bar")
		return
	}

	t.Run("ok WrapDeferred wrap failure err", func(t *testing.T) {
		err := wrapDeferredFailure()
		assert.Equal(t, "foo: bar", err.Error())
	})

	var wrapDeferredFailureWrapped = func() (err error) {
		defer WrapDeferred(&err, "baz")
		err = New("bar")
		err = Wrap(err, "foo")
		return
	}

	t.Run("ok WrapDeferred wrap failure wrapped err", func(t *testing.T) {
		err := wrapDeferredFailureWrapped()
		assert.Equal(t, "baz: foo: bar", err.Error())
	})

	var wrapDeferredNil = func() (err error) {
		defer WrapDeferred(&err, "foo")
		return
	}

	t.Run("ok WrapDeferred wrap nil", func(t *testing.T) {
		err := wrapDeferredNil()
		assert.Nil(t, err)
	})
}
