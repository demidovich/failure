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

	stack, ok := failure.ExtractStack(err)
	if ok {
		for _, line := range stack.Slice() {
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
