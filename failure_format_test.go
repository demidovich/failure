package failure

import (
	"fmt"
	"io"
	"runtime"
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

func TestFormat(t *testing.T) {
	frameFormatter := func(w io.Writer, f runtime.Frame) {
		io.WriteString(w, RelativePath(f.File))
		io.WriteString(w, " ")
		io.WriteString(w, f.Function)
	}

	t.Run("stack_mode_none", func(t *testing.T) {
		SetStackModeNone()

		var (
			err    = bar()
			actual = fmt.Sprintf("%+v", err)
		)

		assert.Equal(t, "Error: foo", actual)
	})

	t.Run("stack_mode_caller", func(t *testing.T) {
		SetStackModeCaller()
		SetStackRootDir(".")
		SetStackFrameFormatter(frameFormatter)

		var (
			err      = bar()
			actual   = fmt.Sprintf("%+v", err)
			expected = strings.Join([]string{
				"Error: foo",
				"Caller: failure_format_test.go github.com/demidovich/failure.foo",
			}, "\n")
		)

		assert.Equal(t, expected, actual)
	})

	t.Run("stack_mode_root", func(t *testing.T) {
		SetStackModeRoot(".")
		SetStackFrameFormatter(frameFormatter)

		var (
			err      = bar()
			actual   = fmt.Sprintf("%+v", err)
			prefix   = stackLinePrefix
			expected = strings.Join([]string{
				"Error: foo\n",
				"Stack Trace:",
				prefix + "failure_format_test.go github.com/demidovich/failure.foo",
				prefix + "failure_format_test.go github.com/demidovich/failure.bar",
			}, "\n")
		)

		assert.Contains(t, actual, expected)
	})

	t.Run("stack_mode_full", func(t *testing.T) {
		SetStackModeFull()
		SetStackRootDir(".")
		SetStackFrameFormatter(frameFormatter)

		var (
			err      = bar()
			actual   = fmt.Sprintf("%+v", err)
			prefix   = stackLinePrefix
			expected = strings.Join([]string{
				"Error: foo\n",
				"Stack Trace:",
				prefix + "failure_format_test.go github.com/demidovich/failure.foo",
				prefix + "failure_format_test.go github.com/demidovich/failure.bar",
			}, "\n")
		)

		assert.Contains(t, actual, expected)
	})

	t.Run("format_verbs", func(t *testing.T) {
		SetStackModeFull()

		tests := []struct {
			message  string
			expected string
			verb     string
		}{
			{"foo", "Error: foo", "%v"},
			{"foo", "foo", "%s"},
			{"foo", "\"foo\"", "%q"},
		}

		for _, tt := range tests {
			err := New(tt.message)
			actual := fmt.Sprintf(tt.verb, err)
			if tt.expected != actual {
				t.Errorf("Verb: %s, expected: %s, actual: %s", tt.verb, tt.expected, actual)
			}
		}
	})

	t.Run("format_s_verb_without_stack", func(t *testing.T) {
		SetStackModeNone()

		var (
			err    = bar()
			actual = fmt.Sprintf("%s", err)
		)

		assert.Equal(t, "foo", actual)
	})
}
