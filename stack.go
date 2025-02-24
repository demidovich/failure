package fail

import (
	"fmt"
	"runtime"
	"strings"
)

func newStack(skipFrames int) *runtime.Frames {
	var stack = make([]uintptr, stackSize)

	depth := runtime.Callers(skipFrames, stack[:])
	stack = stack[:depth]

	return runtime.CallersFrames(stack)
}

func stackToString(frames *runtime.Frames) string {
	b := strings.Builder{}
	for {
		frame, more := frames.Next()
		if stackMode == StackModeApplication && !isApplicationFile(frame.File) {
			break
		}
		b.WriteString(
			fmt.Sprintf("\n%s%s:%d (%s)", stackLogPrefix, relativePath(frame.File), frame.Line, frame.Function),
		)
		if stackMode == StackModeCaller {
			break
		}
		if !more {
			break
		}
	}

	return b.String()
}

func relativePath(file string) string {
	return strings.TrimPrefix(file, appRoot)
}

func isApplicationFile(file string) bool {
	return strings.HasPrefix(file, appRoot)
}
