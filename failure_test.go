package failure

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	err := New("foo")
	assert.Equal(t, "foo", err.Error())
}

func TestNewf(t *testing.T) {
	err := Newf("foo %s", "bar")
	assert.Equal(t, "foo bar", err.Error())
}

func TestWrap(t *testing.T) {
	errA := errors.New("error A")
	errB := fmt.Errorf("%w", errA)
	errC := Wrap(errB, "error C")

	assert.True(t, errors.Is(errC, errA))
}

func TestWrapNil(t *testing.T) {
	err := Wrap(nil, "no error")

	assert.Nil(t, err)
}

func TestWrapf(t *testing.T) {
	errA := errors.New("error A")
	errB := Wrapf(errA, "formatted error %s", "B")

	assert.Equal(t, "formatted error B: error A", errB.Error())
}

func TestWrapFormat(t *testing.T) {
	stackMode = "none"
	errA := errors.New("error A")
	errB := Wrap(errA, "error B")
	msg := fmt.Sprintf("%v", errB)

	assert.Equal(t, "error B: error A", msg)
}

type testingFailureCause interface {
	Cause() error
}

func TestWrapCause(t *testing.T) {
	errA := errors.New("error A")
	errB := Wrap(errA, "error B")

	var cause error
	if e, ok := errB.(testingFailureCause); ok {
		cause = e.Cause()
	}

	assert.Equal(t, "error A", cause.Error())
}

func TestWrapCauseWithoutMessage(t *testing.T) {
	errA := errors.New("")
	errB := Wrap(errA, "error B")

	assert.Equal(t, "error B", errB.Error())
}

func TestWrapfNil(t *testing.T) {
	err := Wrapf(nil, "no error")

	assert.Nil(t, err)
}

func TestIs(t *testing.T) {
	errA := errors.New("error A")
	errB := fmt.Errorf("%w", errA)
	errC := Wrap(errB, "error C")

	assert.True(t, Is(errC, errA))
}

func TestAs(t *testing.T) {
	errA := errors.New("error A")
	errB := Wrap(errA, "error B")

	assert.True(t, As(errB, &errA))
}

func TestStack(t *testing.T) {
	stackMode = StackModeFull
	err := New("foo")

	assert.True(t, len(err.Stack()) > 0)
}

func TestWrapStack(t *testing.T) {
	stackMode = StackModeFull
	err := Wrap(errors.New("foo"), "bar")

	assert.True(t, len(err.Stack()) > 0)
}
