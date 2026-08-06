package failure

import (
	"io"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

type StackMode int

const (
	StackModeNone StackMode = iota + 1
	StackModeCaller
	StackModeRoot
	StackModeFull
)

const (
	stackSkipFrames = 3
)

var (
	stackMode       = StackModeFull
	stackRootDir    string
	stackMaxDepth   = 32
	stackLinePrefix = " -> "
	optionsMu       = sync.Mutex{}
)

var stackFrameFormatter = func(w io.Writer, f runtime.Frame) {
	io.WriteString(w, RelativePath(f.File))
	io.WriteString(w, ":")
	io.WriteString(w, strconv.Itoa(f.Line))

	if idx := strings.LastIndex(f.Function, "."); idx != -1 {
		io.WriteString(w, " - ")
		io.WriteString(w, f.Function[idx+1:])
		io.WriteString(w, "()")
	}
}

// SetStackModeNone disables the stack trace.
func SetStackModeNone() {
	optionsMu.Lock()
	defer optionsMu.Unlock()

	stackMode = StackModeNone
}

// SetStackModeCaller enables a stack trace with the initial caller of the failure.
func SetStackModeCaller() {
	optionsMu.Lock()
	defer optionsMu.Unlock()

	stackMode = StackModeCaller
}

// SetStackModeRoot enables a stack trace limited to project files only.
//   - buildPath: Absolute path to the application root.
//     Don't pass relative paths.
//     Use the value from ldflags.
func SetStackModeRoot(buildPath string) {
	optionsMu.Lock()
	defer optionsMu.Unlock()

	setStackRootDir(buildPath)
	stackMode = StackModeRoot
}

// SetStackModeFull enables a full stack trace.
func SetStackModeFull() {
	optionsMu.Lock()
	defer optionsMu.Unlock()

	stackMode = StackModeFull
}

// SetStackRootDir set stackRootDir variable for shorten file paths in stack trace.
//   - buildPath: Absolute path to the application root.
//     Don't pass relative paths.
//     Use the value from ldflags.
func SetStackRootDir(buildPath string) {
	optionsMu.Lock()
	defer optionsMu.Unlock()

	setStackRootDir(buildPath)
}

func setStackRootDir(buildPath string) {
	stackRootDir, _ = filepath.Abs(buildPath)
	if !strings.HasSuffix(stackRootDir, "/") {
		stackRootDir += "/"
	}
}

// SetStackFrameFormatter set custom stack frame formatter
func SetStackFrameFormatter(f func(w io.Writer, f runtime.Frame)) {
	optionsMu.Lock()
	defer optionsMu.Unlock()

	stackFrameFormatter = f
}

// SetStackMaxDepth set stack trace depth
func SetStackMaxDepth(value int) {
	optionsMu.Lock()
	defer optionsMu.Unlock()

	stackMaxDepth = value
}
