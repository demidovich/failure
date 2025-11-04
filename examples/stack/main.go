package main

import (
	"fmt"
	"os"

	"github.com/demidovich/failure"
)

func main() {
	failure.SetStackMode(failure.StackModeRoot)
	failure.SetStackRootDir("./")

	err := read()

	if e, ok := err.(failure.Error); ok {
		for _, line := range e.Stack() {
			fmt.Println(line)
		}
	}
}

func read() error {
	return missingRead()
}

func missingRead() error {
	_, err := os.ReadFile("/tmp/missing_file")
	return failure.Wrap(err, "read file error")
}
