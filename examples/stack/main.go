package main

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/demidovich/failure"
)

func main() {
	failure.SetStackModeRoot(".") // Do not do this in production. Use pwd only.
	failure.SetStackFrameFormatter(func(w io.Writer, f runtime.Frame) {
		io.WriteString(w, failure.RelativePath(f.File))
		io.WriteString(w, ":")
		io.WriteString(w, strconv.Itoa(f.Line))

		if idx := strings.LastIndex(f.Function, "."); idx != -1 {
			io.WriteString(w, " - ")
			io.WriteString(w, f.Function[idx+1:])
			io.WriteString(w, "()")
		}
	})

	err := read()

	// Example 1

	stack, ok := failure.ExtractStack(err)
	if ok {
		for _, line := range stack.FramesFormatted() {
			fmt.Println(line)
		}
	}

	// Example 2

	stack, ok = failure.ExtractStack(err)
	if ok {
		fmt.Println("")
		for _, frame := range stack.Frames() {
			fmt.Printf("%s:%d - %s\n", failure.RelativePath(frame.File), frame.Line, frame.Function)
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
