package main

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"

	"github.com/demidovich/failure"
)

func main() {
	failure.SetStackModeRoot(".") // Do not do this in production. Use pwd only.
	failure.SetStackFrameFormatter(func(w io.Writer, f runtime.Frame) {
		io.WriteString(w, f.Function)
		io.WriteString(w, ", ")
		io.WriteString(w, failure.RelativePath(f.File))
		io.WriteString(w, ":")
		io.WriteString(w, strconv.Itoa(f.Line))
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
