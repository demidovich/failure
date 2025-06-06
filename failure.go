package failure

import (
	"fmt"
	"io"
)

type Error interface {
	Error() string
	Stack() []string
}

type failure struct {
	message string
	stack   stack
}

func New(message string) Error {
	return &failure{
		message: message,
		stack:   newStack(),
	}
}

func Newf(format string, args ...any) Error {
	return &failure{
		message: fmt.Sprintf(format, args...),
		stack:   newStack(),
	}
}

func (f *failure) Error() string {
	return f.message
}

func (f *failure) Stack() []string {
	return f.stack.Slice()
}

func (f *failure) Format(s fmt.State, verb rune) {
	format(s, verb, f.message, f.stack.String())
}

type wrappedFailure struct {
	message string
	stack   stack
	cause   error
}

func Wrap(err error, message string) Error {
	if err == nil {
		return nil
	}

	return &wrappedFailure{
		message: message,
		stack:   newStack(),
		cause:   err,
	}
}

func Wrapf(err error, format string, args ...any) Error {
	if err == nil {
		return nil
	}

	return &wrappedFailure{
		message: fmt.Sprintf(format, args...),
		stack:   newStack(),
		cause:   err,
	}
}

func (w *wrappedFailure) Error() string {
	if w.cause.Error() == "" {
		return w.message
	} else {
		return w.message + ": " + w.cause.Error()
	}
}

func (w *wrappedFailure) Stack() []string {
	return w.stack.Slice()
}

func (w *wrappedFailure) Format(s fmt.State, verb rune) {
	format(s, verb, w.Error(), w.stack.String())
}

func (w *wrappedFailure) Unwrap() error {
	return w.cause
}

// Universal formatter for wrapped and unwrapped errors
func format(s fmt.State, verb rune, message string, stack string) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = io.WriteString(s, message)
			switch stackMode {
			case StackModeNone:
			case StackModeCaller:
				_, _ = io.WriteString(s, "\n\nCaller:")
				_, _ = io.WriteString(s, stack)
			default:
				_, _ = io.WriteString(s, "\n\nStack Trace:")
				_, _ = io.WriteString(s, stack)
			}
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, message)
	case 'q':
		fmt.Fprintf(s, "%q", message)
	}
}
