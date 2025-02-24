package fail

import (
	"fmt"
	"io"
	"runtime"
)

func New(message string) basic {
	if stackMode == StackModeNone {
		return basic{
			message: message,
			stack:   nil,
		}
	}

	return basic{
		message: message,
		stack:   newStack(3),
	}
}

type basic struct {
	message string
	stack   *runtime.Frames
}

func (e basic) Error() string {
	return e.message
}

func (e basic) Stack() *runtime.Frames {
	return e.stack
}

// func (e basic) Caller() runtime.Frame {
// 	frame, _ := e.stack.Next()
// 	return frame
// }

func (e basic) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = io.WriteString(s, e.message)
			switch stackMode {
			case StackModeNone:
			case StackModeCaller:
				_, _ = io.WriteString(s, "\n\nCaller:")
				_, _ = io.WriteString(s, stackToString(e.stack))
			default:
				_, _ = io.WriteString(s, "\n\nStack Trace:")
				_, _ = io.WriteString(s, stackToString(e.stack))
			}
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, e.message)
	case 'q':
		fmt.Fprintf(s, "%q", e.message)
	}
}
