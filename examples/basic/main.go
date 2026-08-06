package main

import (
	"fmt"
	"os"

	"github.com/demidovich/failure"
)

func main() {
	failure.SetStackModeRoot(".") // Do not do this in production. Use pwd only.

	err := read()

	// Example 1

	fmt.Printf("%+v\n\n", err)

	// Example 2

	stack, ok := failure.ExtractStack(err)
	if ok {
		fmt.Println(stack.Serialize(", "))
	}

	// Example 3

	stack, ok = failure.ExtractStack(err)
	if ok {
		for _, line := range stack.FramesFormatted() {
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
