package failure

import (
	"fmt"
	"runtime"
	"strings"
)

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

func stackToString(frames *runtime.Frames) string {
	b := strings.Builder{}
	for {
		frame, more := frames.Next()
		if stackMode == StackModeRoot && isExternalFile(frame.File) {
			break
		}
		b.WriteString(
			fmt.Sprintf("\n%s%s:%d (%s)", stackPrefix, relativePath(frame.File), frame.Line, frame.Function),
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
	return strings.TrimPrefix(file, stackRootDir)
}

func isExternalFile(file string) bool {
	return !strings.HasPrefix(file, stackRootDir)
}
