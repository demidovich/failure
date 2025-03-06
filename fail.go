package fail

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

type fail struct {
	message string
	stack   *runtime.Frames
}

func New(message string) error {
	return &fail{
		message: message,
		stack:   newStack(3),
	}
}

func Newf(format string, args ...any) error {
	return &fail{
		message: fmt.Sprintf(format, args...),
		stack:   newStack(3),
	}
}

func (b *fail) Error() string {
	return b.message
}

// func (b fail) Caller() runtime.Frame {
// 	frame, _ := b.stack.Next()
// 	return frame
// }

func (b *fail) Format(s fmt.State, verb rune) {
	format(s, verb, b.message, b.stack)
}

type wrappedFail struct {
	message string
	stack   *runtime.Frames
	cause   error
}

func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}

	return &wrappedFail{
		message: message,
		stack:   newStack(3),
		cause:   err,
	}
}

func Wrapf(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}

	return &wrappedFail{
		message: fmt.Sprintf(format, args...),
		stack:   newStack(3),
		cause:   err,
	}
}

func (w *wrappedFail) Error() string {
	return w.message + ": " + w.cause.Error()
}

func (w *wrappedFail) Format(s fmt.State, verb rune) {
	format(s, verb, w.Error(), w.stack)
}

func (w *wrappedFail) Cause() error {
	return w.cause
}

func (w *wrappedFail) Unwrap() error {
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
