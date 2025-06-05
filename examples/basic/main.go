package main

import (
	"fmt"
	"os"

	"github.com/demidovich/failure"
)

func main() {
	failure.SetStackMode(failure.StackModeRoot)
	failure.SetStackRootDir("./")
	failure.SetStackPrefix(" --- ")

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
