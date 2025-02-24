package fail

import (
	"fmt"
	"io"
	"runtime"
)

type basic struct {
	message string
	stack   *runtime.Frames
}

func New(message string) error {
	return &basic{
		message: message,
		stack:   newStack(3),
	}
}

func Newf(format string, args ...interface{}) error {
	return &basic{
		message: fmt.Sprintf(format, args...),
		stack:   newStack(3),
	}
}

func (b *basic) Error() string {
	return b.message
}

// func (b basic) Caller() runtime.Frame {
// 	frame, _ := b.stack.Next()
// 	return frame
// }

func (b *basic) Format(s fmt.State, verb rune) {
	format(s, verb, b.message, b.stack)
}

type wrapped struct {
	message string
	stack   *runtime.Frames
	cause   error
}

func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}

	return &wrapped{
		message: message,
		stack:   newStack(3),
		cause:   err,
	}
}

func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	return &wrapped{
		message: fmt.Sprintf(format, args...),
		stack:   newStack(3),
		cause:   err,
	}
}

func (w *wrapped) Error() string {
	return w.message + ": " + w.cause.Error()
}

func (w *wrapped) Format(s fmt.State, verb rune) {
	format(s, verb, w.Error(), w.stack)
}

func (w *wrapped) Cause() error {
	return w.cause
}

func (w *wrapped) Unwrap() error {
	return w.cause
}

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
