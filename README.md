# failure

[![GoReport][report-img]][report] [![Coverage Status][cov-img]][cov]

[report-img]: https://goreportcard.com/badge/github.com/demidovich/failure
[report]: https://goreportcard.com/report/github.com/demidovich/failure
[cov-img]: https://codecov.io/gh/demidovich/failure/branch/master/graph/badge.svg
[cov]: https://codecov.io/gh/demidovich/failure

Package `failure` is an error handling library for Go with readable stack traces. 

- [Problem](#problem)
- [Usage](#usage)
- [Wrap deferred](#wrap-deferred)
- [Stack trace](#stack-trace)
- [Stack trace mode](#stack-trace-mode)
- [Stack frame formatting](#stack-frame-formatting)
- [Build](#build)
- [Benchmarks](#benchmarks)

Package features:

- Wrap an error multiple times without duplicating the stack trace
- Wrap errors using defer
- Customize stack trace formatting
- Method to get the stack trace from an error
- Method to serialize the stack trace into a string

## Problem

Nothing new here. This topic has been written about many times.

The Go philosophy says that stack traces are not needed because they are too slow and expensive. However, to fix a bug, you must know exactly where it happened — the file and the line number. Because of this, Go developers have to act like detectives, spending time on investigations just to find the source of the problem.

I came to Go from languages with Exceptions, and I know how much easier development is when you can find the exact place of an error instantly. This is not a problem in small services. But services do not stay small forever; they grow bigger over time. At first, you don't notice this problem. But as the project grows, it wastes more and more of your time.

If the error message is unique, everything is simple. But that rarely happens. Many developers even wrap errors and manually add the function name. It seems that people who refused to use stack traces have simply reinvented them.

What we see in logs

```
processing item #89234 terminated unexpectedly: sql: transaction has already been committed or rolled back: context canceled
```

What we actually need to see

```
broken_file.go:100
```

The ideal option

```
processing item #89234 terminated unexpectedly: sql: transaction has already been committed or rolled back: context canceled
broken_file.go:100
```

## Usage

```go
package main

import (
	"fmt"
	"os"

	"github.com/demidovich/failure"
)

func main() {
	failure.SetStackModeRoot(".") // Do not do this in production. Use pwd only.

	err := read()

	// Example 1

	fmt.Printf("%+v\n\n", err)

	// Example 2

	stack, ok := failure.ExtractStack(err)
	if ok {
		fmt.Println(stack.Serialize(", "))
	}

	// Example 3

	stack, ok = failure.ExtractStack(err)
	if ok {
		for _, line := range stack.FramesFormatted() {
			fmt.Println(line)
		}
	}
}

func read() error {
	return missingRead()
}

func missingRead() error {
	_, err := os.ReadFile("/tmp/missing_file")
	return failure.Wrap(err, "read file error")
}
```

Output example 1
```
Error: read file error: open /tmp/missing_file: no such file or directory

Stack Trace:
 -> main.go:42 - missingRead()
 -> main.go:37 - read()
 -> main.go:13 - main()
```

Output example 2
```
main.go:42 - missingRead(), main.go:37 - read(), main.go:13 - main()
```

Output example 3
```
main.go:42 - missingRead()
main.go:37 - read()
main.go:13 - main()
```

## Wrap deferred

```go
func a() (err error) {
	defer failure.WrapDeferred(&err, "a error")

	if err = b(); err != nil {
		return
	}

	if err = c(); err != nil { // failed
		return
	}

	if err = d(); err != nil {
		return
	}

	return nil
}
```

Output
```
Error: a error: c error

Stack Trace:
 -> main.go:29 - a()
 -> main.go:13 - main()
```

## Stack trace

How to get a stack trace from an error.

```go
if stack, ok := failure.ExtractStack(err); ok {

}
```

Slice of formatted stack frames.

```go
if stack, ok := failure.ExtractStack(err); ok {
    for _, line := range stack.FramesFormatted() {
        fmt.Println(line)
    }
}
```

Serialize stack frames. They will be formatted.

```go
if stack, ok := failure.ExtractStack(err); ok {
    fmt.Println(stack.Serialize(", "))
}
```

Slice of raw stack frames.

```go
if stack, ok := failure.ExtractStack(err); ok {
	for _, frame := range e.Frames() {
		fmt.Println(frame.Function)
	}
}
```

## Stack trace mode

Stack trace mode controls the stack trace details. Long paths and frames from internal Go libraries are just noise. You do not need them to find errors.

#### Full

```go
failure.SetStackModeFull()
```

```
Error: read file error: open /tmp/missing_file: no such file or directory

Stack Trace:
 -> /mnt/hdata/code/failure/examples/basic/main.go:35 - missingRead()
 -> /mnt/hdata/code/failure/examples/basic/main.go:30 - read()
 -> /mnt/hdata/code/failure/examples/basic/main.go:13 - main()
 -> /usr/lib/go-1.26/src/runtime/proc.go:290 - main()
 -> /usr/lib/go-1.26/src/runtime/asm_amd64.s:1771 - goexit()
```

#### Root

> [!NOTE]
> Do not use relative paths in `SetStackModeRoot()` and `SetStackRootDir()`. See more in the [build section](#build).

```go
failure.SetStackModeRoot(buildPath)
```

```
Error: read file error: open /tmp/missing_file: no such file or directory

Stack Trace:
 -> main.go:25 - missingRead()
 -> main.go:20 - read()
 -> main.go:15 - main()
```

#### Caller

Example 1

```go
failure.SetStackModeCaller()
```

```
Error: read file error: open /tmp/missing_file: no such file or directory
Caller: /mnt/hdata/code/failure/examples/basic/main.go:36 - missingRead()
```

Example 2

```go
failure.SetStackModeCaller()
failure.SetStackRootDir(buildPath)
```

```
Error: read file error: open /tmp/missing_file: no such file or directory
Caller: main.go:37 - missingRead()
```

#### None

```go
failure.SetStackModeNone()
```

```
Error: read file error: open /tmp/missing_file: no such file or directory
```

## Stack frame formatting

Custom stack frame formatting is configured during application initialization.

> [!NOTE]
> The line prefix ` -> ` is defined within the package and is added only when the error is converted into string.

Example 1

```go
failure.SetStackFrameFormatter(func(w io.Writer, f runtime.Frame) {
	io.WriteString(w, failure.RelativePath(f.File))
	io.WriteString(w, ":")
	io.WriteString(w, strconv.Itoa(f.Line))
})
```

```json
{
  "log.level": "ERROR",
  "message": "queue: SigninStarted: dial tcp :0: connect: connection refused",
  "error.stack_trace": "internal/queue/handlers/signin_started.go:46, internal/queue/queue.go:191"
}
```

Example 2

```go
failure.SetStackFrameFormatter(func(w io.Writer, f runtime.Frame) {
	io.WriteString(w, failure.RelativePath(f.File))
	io.WriteString(w, ":")
	io.WriteString(w, strconv.Itoa(f.Line))

	if idx := strings.LastIndex(f.Function, "/"); idx != -1 {
		io.WriteString(w, " - ")
		io.WriteString(w, f.Function[idx+1:])
		io.WriteString(w, "()")
	}
})
```

```json
{
  "log.level": "ERROR",
  "message": "queue: SigninStarted: dial tcp :0: connect: connection refused",
  "error.stack_trace": "internal/queue/handlers/signin_started.go:46 - handlers.signinStarted.Handle(), internal/queue/queue.go:191 - queue.(*Queue).processTask.func1()"
}
```

## Build

Do not use relative paths in `SetStackModeRoot()` and `SetStackRootDir()`. This will work with `go run` but will break with `go build`.

```go
var buildPath string

func main() {
	failure.SetStackModeRoot(buildPath)
    ...
```

```shell
go build -ldflags "-X main.buildPath=$(pwd)" ./cmd/app/app.go
```

## Benchmarks

I wanted to see how slow stack traces really are. Yes, there is a performance cost, but are a few extra allocations and nanoseconds that important? I think that in 99% of projects, this cost does not matter at all.

```
go clean -cache -testcache
go test -benchmem -count=1 -a ./benchmark/fmt_discard/ -bench=.
goos: linux
goarch: amd64
pkg: trace_errors/benchmark/fmt_discard
cpu: Intel(R) Core(TM) i5-10500T CPU @ 2.30GHz
Benchmark_std_errors_no_stack-12                     	 1000000	      1136 ns/op	     640 B/op	       9 allocs/op
Benchmark_std_errors_caller-12                       	  516914	      2481 ns/op	     968 B/op	      14 allocs/op
Benchmark_github_com_demidovich_failure_caller-12    	  381114	      4782 ns/op	    2457 B/op	      24 allocs/op
Benchmark_github_com_demidovich_failure_root-12      	  288465	      4906 ns/op	    3258 B/op	      24 allocs/op
Benchmark_github_com_demidovich_failure_full-12      	  215730	      7042 ns/op	    3271 B/op	      27 allocs/op
PASS
ok  	trace_errors/benchmark/fmt_discard	12.936s
```

[github.com/demidovich/go-error-trace-test](https://github.com/demidovich/go-error-trace-test)
