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

	formatter := func(f runtime.Frame) string {
		return fmt.Sprintf(" -> %s, %s:%d", f.Function, failure.RelativePath(f.File), f.Line)
	}
	failure.SetStackframeFormatter(formatter)

	err := read()
	fmt.Printf("%+v\n", err)
}

func read() error {
	return missingRead()
}

func missingRead() error {
	_, err := os.ReadFile("/tmp/missing_file")
	return failure.Wrapf(err, "read file error")
}
