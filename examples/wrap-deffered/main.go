package main

import (
	"errors"
	"fmt"

	"github.com/demidovich/failure"
)

func main() {
	failure.SetStackMode(failure.StackModeRoot)
	failure.SetStackRootDir("./")

	err := a()
	fmt.Printf("%+v\n", err)
}

func a() (err error) {
	defer failure.WrapDeferred(&err, "a error")

	err = b()
	if err != nil {
		return
	}

	err = c()
	if err != nil {
		return
	}

	err = d()
	if err != nil {
		return
	}

	return nil
}

func b() error {
	return nil
}

func c() error {
	return nil
}

func d() error {
	return errors.New("c error")
}
