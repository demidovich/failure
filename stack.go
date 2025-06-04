package failure

import (
	"fmt"
	"runtime"
	"strings"
)

var stackFrameFormatter = func(f runtime.Frame) string {
	return fmt.Sprintf("%s%s:%d (%s)", stackPrefix, relativePath(f.File), f.Line, f.Function)
}

func SetStackFrameFormatter(f func(f runtime.Frame) string) {
	stackFrameFormatter = f
}

func newStack(skipFrames int) *runtime.Frames {
	if stackMode == StackModeNone {
		return nil
	}

	const depth = 32
	var pcs = make([]uintptr, depth)

	size := runtime.Callers(skipFrames, pcs[:])
	pcs = pcs[:size]

	return runtime.CallersFrames(pcs)
}

func stackToSlice(frames *runtime.Frames) []string {
	result := make([]string, 0, 32)
	for {
		frame, more := frames.Next()
		if stackMode == StackModeRoot && isExternalFile(frame.File) {
			break
		}
		result = append(
			result,
			stackFrameFormatter(frame),
		)
		if stackMode == StackModeCaller {
			break
		}
		if !more {
			break
		}
	}
	return result
}

func stackToString(frames *runtime.Frames) string {
	b := strings.Builder{}
	for {
		frame, more := frames.Next()
		if stackMode == StackModeRoot && isExternalFile(frame.File) {
			break
		}
		b.WriteString("\n")
		b.WriteString(stackFrameFormatter(frame))
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
	return strings.TrimPrefix(file, stackRootDir)
}

func isExternalFile(file string) bool {
	return !strings.HasPrefix(file, stackRootDir)
}
