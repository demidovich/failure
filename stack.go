package failure

import (
	"errors"
	"runtime"
	"strings"
)

// Stack contains a collection of stack frames and methods to work with them.
type Stack struct {
	frames []runtime.Frame
}

func newStack() *Stack {
	if stackMode == StackModeNone {
		return nil
	}

	pcs := make([]uintptr, stackMaxDepth)
	size := runtime.Callers(stackSkipFrames, pcs)
	if size < 1 {
		return nil
	}

	if stackMode == StackModeCaller {
		size = 1
	}

	frames := runtime.CallersFrames(pcs[:size])
	s := Stack{
		frames: make([]runtime.Frame, 0, size),
	}

	for {
		frame, more := frames.Next()
		if stackMode == StackModeRoot && IsExternalFile(frame.File) {
			break
		}
		s.frames = append(s.frames, frame)
		if stackMode == StackModeCaller {
			break
		}
		if !more {
			break
		}
	}

	return &s
}

// Frames return slice of runtime frame.
func (s *Stack) Frames() []runtime.Frame {
	if s == nil {
		return nil
	}

	return s.frames
}

// FramesFormatted return slice of formatted runtime frames.
func (s *Stack) FramesFormatted() []string {
	if s == nil {
		return nil
	}

	f := make([]string, 0, len(s.frames))
	b := strings.Builder{}
	b.Grow(256) // The number of bytes that stackFrameFormatter will write is unknown.
	for i := range s.frames {
		stackFrameFormatter(&b, s.frames[i])
		f = append(f, b.String())
		b.Reset()
	}

	return f
}

// Serialize returns formatted runtime frames as a string.
func (s *Stack) Serialize(sep string) string {
	if s == nil {
		return ""
	}

	b := strings.Builder{}
	b.Grow(len(s.frames) * 100) // The number of bytes that stackFrameFormatter will write is unknown.
	for i := range s.frames {
		if i > 0 {
			b.WriteString(sep)
		}
		stackFrameFormatter(&b, s.frames[i])
	}

	return b.String()
}

// RelativePath returns a shortened path if the application root was specified.
func RelativePath(file string) string {
	return strings.TrimPrefix(file, stackRootDir)
}

// ExtractStack failure stack trace from wrapped chain.
func ExtractStack(err error) (*Stack, bool) {
	for {
		if err == nil {
			return nil, false
		}
		if e, ok := err.(Error); ok { //nolint:errorlint
			s := e.Stack()
			if s != nil {
				return s, true
			}
		}
		err = errors.Unwrap(err)
	}
}

// IsExternalFile returns true if the file path does not start with the failure root prefix.
func IsExternalFile(file string) bool {
	return !strings.HasPrefix(file, stackRootDir)
}
