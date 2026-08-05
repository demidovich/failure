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
		return fmt.Sprintf(" -> %s, %s:%d", f.Function, failure.RelativePath(f.File), f.Line)
	})

	err := read()
	fmt.Printf("%+v\n", err)
}

func read() error {
	return missingRead()
}

func missingRead() error {
	_, err := os.ReadFile("/tmp/missing_file")
	return failure.Wrap(err, "read file error")
}
