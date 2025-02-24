package fail

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirectoryWithExistingFile(t *testing.T) {
	pwd, _ := os.Getwd()
	found := directoryWithFile("go.mod")

	assert.Equal(t, pwd, found)
}

func TestDirectoryWithMissingFile(t *testing.T) {
	found := directoryWithFile("missing.file")

	assert.Equal(t, "", found)
}

func TestSetAppRoot(t *testing.T) {
	value := "/tmp"
	SetAppRoot(value)

	assert.Equal(t, value+"/", appRoot)
}

func TestSetStackMode(t *testing.T) {
	value := StackModeNone
	SetStackMode(value)

	assert.Equal(t, value, stackMode)
}

func TestSetStackLogPrefix(t *testing.T) {
	value := "----->"
	SetStackLogPrefix(value)

	assert.Equal(t, value, stackLogPrefix)
}
