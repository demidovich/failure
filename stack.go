package failure

import (
	"runtime"
	"strings"
)

type stack struct {
	frames   *runtime.Frames
	hasSlice bool
	slice    []string
}

func newStack() stack {
	if stackMode == StackModeNone {
		return stack{}
	}

	var pcs = make([]uintptr, stackDepth)

	size := runtime.Callers(skipStackFrames, pcs[:])
	pcs = pcs[:size]

	return stack{
		frames: runtime.CallersFrames(pcs),
	}
}

func (s *stack) Frames() *runtime.Frames {
	return s.frames
}

func (s *stack) Slice() []string {
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

func (s *stack) String() string {
	b := strings.Builder{}
	for _, line := range s.Slice() {
		b.WriteString("\n")
		b.WriteString(line)
	}

	return b.String()
}

func (s *stack) isExternalFile(file string) bool {
	return !strings.HasPrefix(file, stackRootDir)
}

// RelativePath returns a shortened path if the application root was specified
func RelativePath(file string) string {
	return strings.TrimPrefix(file, stackRootDir)
}
