package failure

import (
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
