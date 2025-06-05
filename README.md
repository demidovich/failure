# failure

[![GoReport][report-img]][report] [![Coverage Status][cov-img]][cov]

Package `failure` is an error handling library for Go with readable stack traces.

[report-img]: https://goreportcard.com/badge/github.com/demidovich/failure
[report]: https://goreportcard.com/report/github.com/demidovich/failure
[cov-img]: https://codecov.io/gh/demidovich/failure/branch/master/graph/badge.svg
[cov]: https://codecov.io/gh/demidovich/failure

## Usage

```go
package main

import (
	"fmt"
	"os"

	"github.com/demidovich/failure"
)

func main() {
	failure.SetStackMode(failure.StackModeRoot)
	failure.SetStackRootDir("./")
	failure.SetStackPrefix("  --- ")

	err := read()
	fmt.Printf("%+v\n", err)
}

func read() error {
	return missingRead()
}

func missingRead() error {
	_, err := os.ReadFile("/tmp/missing_file")
	return failure.Wrap(err, "read file error")
}
```

```
read file error: open /tmp/missing_file: no such file or directory

Stack Trace:
 --- main.go:25 (main.missingRead)
 --- main.go:20 (main.read)
 --- main.go:15 (main.main)
```

## StackMode

Сontrol of stack trace display mode.

### Full

```go
failure.SetStackMode(failure.StackModeFull)
```

```
read file error: open /tmp/missing_file: no such file or directory

Stack Trace:
 --- /mnt/hdata/code/go/failure/examples/basic/main.go:25 (main.missingRead)
 --- /mnt/hdata/code/go/failure/examples/basic/main.go:20 (main.read)
 --- /mnt/hdata/code/go/failure/examples/basic/main.go:15 (main.main)
 --- /usr/lib/go-1.24/src/runtime/proc.go:283 (runtime.main)
 --- /usr/lib/go-1.24/src/runtime/asm_amd64.s:1700 (runtime.goexit)
```

### Root

```go
failure.SetStackMode(failure.StackModeRoot)
failure.SetStackRootDir("./")
```

```
read file error: open /tmp/missing_file: no such file or directory

Stack Trace:
 --- main.go:25 (main.missingRead)
 --- main.go:20 (main.read)
 --- main.go:15 (main.main)
```

### Caller

```go
failure.SetStackMode(failure.StackModeCaller)
failure.SetStackRootDir("./")
```

```
read file error: open /tmp/missing_file: no such file or directory

Caller:
 --- main.go:25 (main.missingRead)
```

### None

```go
failure.SetStackMode(failure.StackModeNone)
```

```
read file error: open /tmp/missing_file: no such file or directory
```

## Stack

```go
if e, ok := err.(failure.Error); ok {
	for _, line := range e.Stack() {
		fmt.Println(line)
	}
}
```

```
main.go:31 (main.missingRead)
main.go:26 (main.read)
main.go:15 (main.main)
```

## Stackframe formatting

```go
formatter := func(f runtime.Frame) string {
    return fmt.Sprintf(" -> %s, %s:%d", f.Function, failure.RelativePath(f.File), f.Line)
}

failure.SetStackframeFormatter(formatter)
```

```
 -> main.missingRead, main.go:30
 -> main.read, main.go:25
 -> main.main, main.go:20
```
