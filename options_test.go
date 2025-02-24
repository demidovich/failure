package fail

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGuessedAppRoot(t *testing.T) {
	pwd, _ := os.Getwd()
	guessed := guessedAppRoot()

	assert.Equal(t, pwd, guessed)
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
