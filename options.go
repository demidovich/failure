package fail

import (
	"os"
	"path/filepath"
	"strings"
)

var (
	appRoot        string
	stackSize      int       = 32
	stackMode      StackMode = StackModeApplication
	stackLogPrefix string    = "-> "
)

type StackMode string

const (
	StackModeNone        StackMode = "none"
	StackModeCaller      StackMode = "caller"
	StackModeApplication StackMode = "application"
	StackModeFull        StackMode = "full"
)

func init() {
	SetAppRoot(
		guessedAppRoot(),
	)
}

func guessedAppRoot() string {
	pwd, _ := os.Getwd()
	for {
		if _, err := os.Stat(pwd + "/go.mod"); err == nil {
			return pwd
		}
		pwd = filepath.Dir(pwd)
		if pwd == "/" {
			break
		}
	}

	return ""
}

// Set appRoot variable for shorten file names in stack trace logs
func SetAppRoot(value string) {
	appRoot = strings.TrimSuffix(value, "/") + "/"
}

// Set stackMode variable for change stack trace verbosity
func SetStackMode(value StackMode) {
	stackMode = value
}

// Set stackLogPrefix variable for readability logs
func SetStackLogPrefix(value string) {
	stackLogPrefix = value
}
