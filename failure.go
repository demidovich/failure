package failure

import (
	"errors"
	"fmt"
	"io"
	"runtime"
)

func As(err error, target any) bool {
	return errors.As(err, target)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

type failure struct {
	message string
	stack   *runtime.Frames
}

func New(message string) error {
	return &failure{
		message: message,
		stack:   newStack(3),
	}
}

func Newf(format string, args ...any) error {
	return &failure{
		message: fmt.Sprintf(format, args...),
		stack:   newStack(3),
	}
}

func (b *failure) Error() string {
	return b.message
}

// func (b failure) Caller() runtime.Frame {
// 	frame, _ := b.stack.Next()
// 	return frame
// }

func (b *failure) Format(s fmt.State, verb rune) {
	format(s, verb, b.message, b.stack)
}

type wrappedFailure struct {
	message string
	stack   *runtime.Frames
	cause   error
}

func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}

	return &wrappedFailure{
		message: message,
		stack:   newStack(3),
		cause:   err,
	}
}

func Wrapf(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}

	return &wrappedFailure{
		message: fmt.Sprintf(format, args...),
		stack:   newStack(3),
		cause:   err,
	}
}

func (w *wrappedFailure) Error() string {
	return w.message + ": " + w.cause.Error()
}

func (w *wrappedFailure) Format(s fmt.State, verb rune) {
	format(s, verb, w.Error(), w.stack)
}

func (w *wrappedFailure) Cause() error {
	return w.cause
}

func (w *wrappedFailure) Unwrap() error {
	return w.cause
}

// Universal formatter for wrapped and unwrapped errors
func format(s fmt.State, verb rune, message string, stack *runtime.Frames) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = io.WriteString(s, message)
			switch stackMode {
			case StackModeNone:
			case StackModeCaller:
				_, _ = io.WriteString(s, "\n\nCaller:")
				_, _ = io.WriteString(s, stackToString(stack))
			default:
				_, _ = io.WriteString(s, "\n\nStack Trace:")
				_, _ = io.WriteString(s, stackToString(stack))
			}
			return
		}
		fallthrough
	case 's', 'q':
		_, _ = io.WriteString(s, message)
	}
}
