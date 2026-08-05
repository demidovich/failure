package failure

import (
	"errors"
	"runtime"
	"strings"
)

type Stack struct {
	frames   *runtime.Frames
	hasSlice bool
	slice    []string
}

func newStack() *Stack {
	if stackMode == StackModeNone {
		return nil
	}

	var pcs = make([]uintptr, stackDepth)
	size := runtime.Callers(skipStackFrames, pcs[:])
	pcs = pcs[:size]

	return &Stack{
		frames: runtime.CallersFrames(pcs),
	}
}

func (s *Stack) Frames() *runtime.Frames {
	if s == nil {
		return nil
	}

	return s.frames
}

func (s *Stack) Slice() []string {
	if s == nil || s.frames == nil {
		return nil
	}

	if s.hasSlice {
		return s.slice
	}

	s.slice = make([]string, 0, stackDepth)
	if s.frames == nil {
		return s.slice
	}

	for {
		frame, more := s.frames.Next()
		if stackMode == StackModeRoot && s.isExternalFile(frame.File) {
			break
		}
		s.slice = append(
			s.slice,
			stackframeFormatter(frame),
		)
		if stackMode == StackModeCaller {
			break
		}
		if !more {
			break
		}
	}

	s.hasSlice = true
	return s.slice
}

func (s *Stack) String() string {
	if s == nil || s.frames == nil {
		return ""
	}

	b := strings.Builder{}
	for _, line := range s.Slice() {
		b.WriteString("\n")
		b.WriteString(line)
	}

	return b.String()
}

func (s *Stack) isExternalFile(file string) bool {
	return !strings.HasPrefix(file, stackRootDir)
}

// RelativePath returns a shortened path if the application root was specified
func RelativePath(file string) string {
	return strings.TrimPrefix(file, stackRootDir)
}

// ExtractStack extract failure stack
func ExtractStack(err error) (*Stack, bool) {
	for {
		if err == nil {
			return nil, false
		}
		if e, ok := err.(Error); ok {
			s := e.Stack()
			if s != nil {
				return s, true
			}
		}
		err = errors.Unwrap(err)
	}
}
