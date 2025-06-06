package failure

import (
	"errors"
	"fmt"
	"io"
)

func As(err error, target any) bool {
	return errors.As(err, target)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

type Error interface {
	Error() string
	Stack() []string
}

type failure struct {
	err   error
	stack *stack
}

// New makes an Error from error or string.
func New(e any) Error {
	var err error
	var stack *stack

	switch e := e.(type) {
	case *failure:
		err = e.err
		stack = e.stack
	case error:
		err = e
		stack = newStack()
	default:
		err = fmt.Errorf("%v", e)
		stack = newStack()
	}

	return &failure{
		err:   err,
		stack: stack,
	}
}

// Newf makes an Error from formatted string.
func Newf(format string, args ...any) Error {
	return &failure{
		err:   fmt.Errorf(format, args...),
		stack: newStack(),
	}
}

// Wrap makes an Error from given error.
func Wrap(e any) Error {
	if e == nil {
		return nil
	}

	var err error
	var stack *stack

	switch e := e.(type) {
	case *failure:
		return e
	case error:
		err = e
		stack = newStack()
	default:
		err = fmt.Errorf("%v", e)
		stack = newStack()
	}

	return &failure{
		err:   err,
		stack: stack,
	}
}

// Wrap makes an Error from given error.
func Wrapf(e error, format string, args any) Error {
	if e == nil {
		return nil
	}

	var err error
	var stack *stack

	switch e := e.(type) {
	case *failure:
		return e
	case error:
		err = e
		stack = newStack()
	default:
		err = fmt.Errorf("%v", e)
		stack = newStack()
	}

	return &failure{
		err:   err,
		stack: stack,
	}
}

func (f *failure) Error() string {
	return f.err.Error()
}

func (f *failure) Unwrap() error {
	return f.err
}

func (f *failure) Stack() []string {
	return f.stack.Slice()
}

func (f *failure) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = io.WriteString(s, f.err.Error())
			switch stackMode {
			case StackModeNone:
			case StackModeCaller:
				_, _ = io.WriteString(s, "\n\nCaller:")
				_, _ = io.WriteString(s, f.stack.String())
			default:
				_, _ = io.WriteString(s, "\n\nStack Trace:")
				_, _ = io.WriteString(s, f.stack.String())
			}
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, f.err.Error())
	case 'q':
		fmt.Fprintf(s, "%q", f.err.Error())
	}
}
