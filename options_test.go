package failure

import (
	"io"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	t.Run("ok_set_stack_mode_none", func(t *testing.T) {
		stackMode = StackModeCaller
		SetStackModeNone()

		assert.Equal(t, StackModeNone, stackMode)
	})

	t.Run("ok_set_stack_mode_caller", func(t *testing.T) {
		stackMode = StackModeNone
		SetStackModeCaller()

		assert.Equal(t, StackModeCaller, stackMode)
	})

	t.Run("ok_set_stack_mode_root", func(t *testing.T) {
		stackMode = StackModeNone
		SetStackModeRoot("/tmp")

		assert.Equal(t, StackModeRoot, stackMode)
		assert.Equal(t, "/tmp/", stackRootDir)
	})

	t.Run("ok_set_stack_mode_full", func(t *testing.T) {
		stackMode = StackModeNone
		SetStackModeFull()

		assert.Equal(t, StackModeFull, stackMode)
	})

	t.Run("ok_set_stack_root_dir", func(t *testing.T) {
		value := "/tmp"
		SetStackRootDir(value)

		assert.Equal(t, value+"/", stackRootDir)
	})

	t.Run("ok_set_stack_formatter", func(t *testing.T) {
		t.Run("ok_set_stack_formatter", func(t *testing.T) {
			SetStackFrameFormatter(func(w io.Writer, f runtime.Frame) {
				io.WriteString(w, f.File)
				io.WriteString(w, strconv.Itoa(f.Line))
				io.WriteString(w, f.Function)
			})

			b := &strings.Builder{}
			stackFrameFormatter(b, runtime.Frame{
				File:     "k",
				Line:     1,
				Function: "m",
			})

			assert.Equal(t, "k1m", b.String())
		})
	})

	t.Run("ok_set_stack_max_depth", func(t *testing.T) {
		stackMaxDepth = 1
		SetStackMaxDepth(2)

		assert.Equal(t, 2, stackMaxDepth)
	})
}
