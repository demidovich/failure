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

func TestFormatStackModeNone(t *testing.T) {
	SetStackMode(StackModeNone)

	err := bar()
	msg := fmt.Sprintf("%+v", err)

	assert.Equal(t, "foo", msg)
}

func TestFormatStackModeCaller(t *testing.T) {
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

func TestFormatStackModeApplication(t *testing.T) {
	SetStackMode(StackModeApplication)
	SetStackLogPrefix("")

	err := bar()
	msg := fmt.Sprintf("%+v", err)

	expected := strings.Join([]string{
		"foo\n",
		"Stack Trace:",
		"fail_test.go:12 (github.com/demidovich/fail.foo)",
		"fail_test.go:16 (github.com/demidovich/fail.bar)",
		"fail_test.go:77 (github.com/demidovich/fail.TestFormatStackModeApplication)",
	}, "\n")

	assert.Equal(t, expected, msg)
}

func TestFormatStackModeFull(t *testing.T) {
	SetStackMode(StackModeFull)
	SetStackLogPrefix("")

	err := bar()
	msg := fmt.Sprintf("%+v", err)

	prefix := strings.Join([]string{
		"foo\n",
		"Stack Trace:",
		"fail_test.go:12 (github.com/demidovich/fail.foo)",
		"fail_test.go:16 (github.com/demidovich/fail.bar)",
		"fail_test.go:95 (github.com/demidovich/fail.TestFormatStackModeFull)",
	}, "\n")

	assert.Contains(t, msg, prefix)
}

func TestFormatVerbs(t *testing.T) {
	tests := []struct {
		expected string
		verb     string
	}{
		{"foo", "%v"},
		{"foo", "%s"},
	}

	SetStackMode(StackModeFull)

	for _, tt := range tests {
		err := New(tt.expected)
		actual := fmt.Sprintf(tt.verb, err)
		if tt.expected != actual {
			t.Errorf("Verb: %s, expected: %s, actual: %s", tt.verb, tt.expected, actual)
		}
	}
}

func TestFormatS(t *testing.T) {
	SetStackMode(StackModeNone)

	err := bar()
	msg := fmt.Sprintf("%s", err)

	assert.Equal(t, "foo", msg)
}
