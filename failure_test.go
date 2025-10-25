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

	t.Run("New", func(t *testing.T) {
		err := New("foo")
		assert.Equal(t, "foo", err.Error())
	})

	t.Run("Newf", func(t *testing.T) {
		err := Newf("foo %s", "bar")
		assert.Equal(t, "foo bar", err.Error())
	})

	t.Run("Wrap", func(t *testing.T) {
		errA := errors.New("error A")
		errB := fmt.Errorf("%w", errA)
		errC := Wrapf(errB, "error C")

		assert.True(t, errors.Is(errC, errA))
	})

	t.Run("Wrap nil", func(t *testing.T) {
		err := Wrap(nil)

		assert.Nil(t, err)
	})

	t.Run("Wrapf", func(t *testing.T) {
		errA := errors.New("error A")
		errB := Wrapf(errA, "formatted error %s", "B")

		assert.Equal(t, "formatted error B: error A", errB.Error())
	})

	t.Run("Wrapf nil", func(t *testing.T) {
		err := Wrapf(nil, "")

		assert.Nil(t, err)
	})

	t.Run("Wrapf format", func(t *testing.T) {
		stackMode = "none"
		errA := errors.New("error A")
		errB := Wrapf(errA, "error B")
		msg := fmt.Sprintf("%v", errB)

		assert.Equal(t, "error B: error A", msg)
	})

	t.Run("WrapCause", func(t *testing.T) {
		errA := errors.New("error A")
		errB := Wrapf(errA, "error B")

		var cause error
		if e, ok := errB.(testingWrappedFailure); ok {
			cause = e.Unwrap()
		}

		assert.Equal(t, "error A", cause.Error())
	})

	t.Run("WrapCause without message", func(t *testing.T) {
		errA := errors.New("")
		errB := Wrapf(errA, "error B")

		assert.Equal(t, "error B", errB.Error())
	})

	t.Run("WrapCause without message", func(t *testing.T) {
		errA := errors.New("")
		errB := Wrapf(errA, "error B")

		assert.Equal(t, "error B", errB.Error())
	})

	t.Run("New stack", func(t *testing.T) {
		stackMode = StackModeFull
		err := New("foo")

		assert.True(t, len(err.Stack()) > 0)
	})

	t.Run("Wrap stack", func(t *testing.T) {
		stackMode = StackModeFull
		err := Wrapf(errors.New("foo"), "bar")

		assert.True(t, len(err.Stack()) > 0)
	})
}
