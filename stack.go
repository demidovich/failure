package failure

import (
	"fmt"
	"runtime"
	"strings"
)

const stackSkipFrames = 3

var stackframeFormatter = func(f runtime.Frame) string {
	return fmt.Sprintf("%s%s:%d (%s)", stackPrefix, RelativePath(f.File), f.Line, f.Function)
}

func SetStackframeFormatter(f func(f runtime.Frame) string) {
	stackframeFormatter = f
}

func RelativePath(file string) string {
	return strings.TrimPrefix(file, stackRootDir)
}

func newStack() *stack {
	if stackMode == StackModeNone {
		return &stack{}
	}

	const depth = 32
	var pcs = make([]uintptr, depth)

	size := runtime.Callers(stackSkipFrames, pcs[:])
	pcs = pcs[:size]

	return &stack{
		frames: runtime.CallersFrames(pcs),
	}
}

type stack struct {
	frames   *runtime.Frames
	hasSlice bool
	slice    []string
}

func (s *stack) Slice() []string {
	if s.hasSlice {
		return s.slice
	}

	s.slice = make([]string, 0, 32)
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
