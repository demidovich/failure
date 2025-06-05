package failure

import (
	"path/filepath"
)

var (
	stackMode    StackMode = StackModeFull
	stackRootDir string
	stackPrefix  string
)

type StackMode string

const (
	StackModeNone   StackMode = "none"
	StackModeCaller StackMode = "caller"
	StackModeRoot   StackMode = "root"
	StackModeFull   StackMode = "full"
)

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
