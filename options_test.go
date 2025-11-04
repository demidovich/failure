package failure

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	t.Run("ok SetStackRootDir", func(t *testing.T) {
		value := "/tmp"
		SetStackRootDir(value)

		assert.Equal(t, value+"/", stackRootDir)
	})

	t.Run("ok SetStackMode", func(t *testing.T) {
		value := StackModeNone
		SetStackMode(value)

		assert.Equal(t, value, stackMode)
	})

	t.Run("ok SetStackPrefix", func(t *testing.T) {
		value := "----->"
		SetStackPrefix(value)

		assert.Equal(t, value, stackPrefix)
	})

	t.Run("ok SetStackDeth", func(t *testing.T) {
		value := 10
		SetStackDepth(value)

		assert.Equal(t, value, stackDepth)
	})

	t.Run("fails SetStackDeth on invalid value", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("don`t panic")
			}
		}()

		SetStackDepth(0)
		SetStackDepth(-1)
	})

	t.Run("ok SetStackframeFormatter", func(t *testing.T) {
		formatter := func(f runtime.Frame) string {
			return fmt.Sprintf("%s %d %s", f.File, f.Line, f.Function)
		}

		frame := runtime.Frame{
			File:     "a",
			Line:     10,
			Function: "b",
		}

		SetStackframeFormatter(formatter)

		assert.Equal(t, "a 10 b", stackframeFormatter(frame))
	})
}
