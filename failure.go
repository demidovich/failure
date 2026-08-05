package failure

import (
	"fmt"
	"io"
)

type Error interface {
	Error() string
	Stack() *Stack
}

type failure struct {
	message string
	stack   *Stack
}

// New makes an Error with formatted message from the given value.
func New(format string, args ...any) Error {
	if len(args) > 0 {
		format = fmt.Sprintf(format, args...)
	}

	return &failure{
		message: format,
		stack:   newStack(),
	}
}

func (f *failure) Error() string {
	return f.message
}

func (f *failure) Stack() *Stack {
	return f.stack
}

func (f *failure) Format(s fmt.State, verb rune) {
	format(s, verb, f.message, f.stack.String())
}

type wrappedFailure struct {
	message string
	stack   *Stack
	cause   error
}

// Wrap makes an wrapped Error with formatted message from the given value.
func Wrap(err error, format string, args ...any) Error {
	if err == nil {
		return nil
	}

	if len(args) > 0 {
		format = fmt.Sprintf(format, args...)
	}

	w := &wrappedFailure{
		message: format,
		cause:   err,
	}

	if s, ok := ExtractStack(err); ok {
		w.stack = s
	} else {
		w.stack = newStack()
	}

	return w
}

// WrapDeferred makes an deferred Error with formatted message from the given value.
func WrapDeferred(errP *error, format string, args ...any) {
	if *errP == nil {
		return
	}

	if len(args) > 0 {
		format = fmt.Sprintf(format, args...)
	}

	var stack *Stack
	if s, ok := ExtractStack(*errP); ok {
		stack = s
	} else {
		stack = newStack()
	}

	*errP = &wrappedFailure{
		message: format,
		cause:   *errP,
		stack:   stack,
	}
}

func (w *wrappedFailure) Error() string {
	if w.cause.Error() == "" {
		return w.message
	} else {
		return w.message + ": " + w.cause.Error()
	}
}

func (w *wrappedFailure) Stack() *Stack {
	return w.stack
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
