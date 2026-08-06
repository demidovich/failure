package main

import (
	"errors"
	"fmt"

	"github.com/demidovich/failure"
)

func main() {
	failure.SetStackModeRoot(".")

	err := a()
	fmt.Printf("%+v\n", err)
}

func a() (err error) {
	defer failure.WrapDeferred(&err, "a error")

	if err = b(); err != nil {
		return
	}

	if err = c(); err != nil {
		return
	}

	if err = d(); err != nil {
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
