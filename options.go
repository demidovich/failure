package failure

import (
	"fmt"
	"path/filepath"
	"runtime"
)

type StackMode string

const (
	StackModeNone   StackMode = "none"
	StackModeCaller StackMode = "caller"
	StackModeRoot   StackMode = "root"
	StackModeFull   StackMode = "full"
	skipStackFrames           = 3
)

var (
	stackMode           StackMode = StackModeFull
	stackRootDir        string
	stackPrefix         string
	stackDepth          int = 32
	stackframeFormatter     = func(f runtime.Frame) string {
		return fmt.Sprintf("%s%s:%d (%s)", stackPrefix, RelativePath(f.File), f.Line, f.Function)
	}
)

func SetStackframeFormatter(f func(f runtime.Frame) string) {
	stackframeFormatter = f
}

// Set stackMode variable for change stack trace verbosity
func SetStackMode(value StackMode) {
	stackMode = value
}

// Set stackRootDir variable for shorten file names in stack trace
func SetStackRootDir(value string) {
	stackRootDir, _ = filepath.Abs(value)
	if stackRootDir != "/" {
		stackRootDir = stackRootDir + "/"
	}
}

// Set stackPrefix variable for readability logs
func SetStackPrefix(value string) {
	stackPrefix = value
}

func SetStackDepth(value int) {
	if value < 1 {
		panic("stack depth cannot be less than 1")
	}
	stackDepth = value
}
