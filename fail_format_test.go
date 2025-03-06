package fail

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func foo() error {
	return New("foo")
}

func bar() error {
	return foo()
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
		"fail_format_test.go:12 (github.com/demidovich/fail.foo)",
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
		"fail_format_test.go:12 (github.com/demidovich/fail.foo)",
		"fail_format_test.go:16 (github.com/demidovich/fail.bar)",
		"fail_format_test.go:48 (github.com/demidovich/fail.TestFormatStackModeApplication)",
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
		"fail_format_test.go:12 (github.com/demidovich/fail.foo)",
		"fail_format_test.go:16 (github.com/demidovich/fail.bar)",
		"fail_format_test.go:66 (github.com/demidovich/fail.TestFormatStackModeFull)",
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
		{"foo", "%q"},
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
