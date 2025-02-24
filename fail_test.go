package fail

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func foo() basic {
	return New("foo")
}

func bar() basic {
	return foo()
}

func TestError(t *testing.T) {
	err := New("foo")
	assert.Equal(t, "foo", err.Error())
}

func TestStackNone(t *testing.T) {
	stackMode = StackModeNone

	err := bar()

	assert.Nil(t, err.Stack())
}

func TestStack(t *testing.T) {
	stackMode = StackModeApplication

	err := bar()
	stack := err.Stack()

	frame1, _ := stack.Next()
	assert.Equal(t, "github.com/demidovich/fail.foo", frame1.Function)
	assert.Equal(t, "fail_test.go", relativePath(frame1.File))
	assert.Equal(t, 12, frame1.Line)

	frame2, _ := stack.Next()
	assert.Equal(t, "github.com/demidovich/fail.bar", frame2.Function)
	assert.Equal(t, "fail_test.go", relativePath(frame2.File))
	assert.Equal(t, 16, frame2.Line)
}

func TestFormatModeNone(t *testing.T) {
	SetStackMode(StackModeNone)

	err := bar()
	msg := fmt.Sprintf("%+v", err)

	assert.Equal(t, "foo", msg)
}

func TestFormatModeCaller(t *testing.T) {
	SetStackMode(StackModeCaller)
	SetStackLogPrefix("")

	err := bar()
	msg := fmt.Sprintf("%+v", err)

	expected := strings.Join([]string{
		"foo\n",
		"Caller:",
		"fail_test.go:12 (github.com/demidovich/fail.foo)",
	}, "\n")

	assert.Equal(t, expected, msg)
}

func TestFormatModeApplication(t *testing.T) {
	SetStackMode(StackModeApplication)
	SetStackLogPrefix("")

	err := bar()
	msg := fmt.Sprintf("%+v", err)

	expected := strings.Join([]string{
		"foo\n",
		"Stack Trace:",
		"fail_test.go:12 (github.com/demidovich/fail.foo)",
		"fail_test.go:16 (github.com/demidovich/fail.bar)",
		"fail_test.go:78 (github.com/demidovich/fail.TestFormatModeApplication)",
	}, "\n")

	assert.Equal(t, expected, msg)
}
