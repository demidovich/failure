package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/demidovich/failure"
)

func main() {
	failure.SetStackMode(failure.StackModeRoot)
	failure.SetStackRootDir("./")
	failure.SetStackframeFormatter(func(f runtime.Frame) string {
		return fmt.Sprintf("%s (%d)", failure.RelativePath(f.File), f.Line)
	})

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
