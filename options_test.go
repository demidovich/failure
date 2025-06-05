package failure

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetStackRootDir(t *testing.T) {
	value := "/tmp"
	SetStackRootDir(value)

	assert.Equal(t, value+"/", stackRootDir)
}

func TestSetStackMode(t *testing.T) {
	value := StackModeNone
	SetStackMode(value)

	assert.Equal(t, value, stackMode)
}

func TestSetStackPrefix(t *testing.T) {
	value := "----->"
	SetStackPrefix(value)

	assert.Equal(t, value, stackPrefix)
}

func TestSetStackframeFormatter(t *testing.T) {
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
}
