# API Documentation

## Overview

`goleakwatch` provides utilities to detect goroutine leaks in Go applications. It's designed to be lightweight and easy to integrate into your development and testing workflow.

## Core Functions

### `Check(fn func(), cfg Config) error`

Runs the given function and checks for goroutine leaks.

**Parameters:**
- `fn`: The function to execute and monitor
- `cfg`: Configuration options for the leak detection

**Returns:**
- `error`: Returns an error if goroutine leaks are detected, nil otherwise

**Example:**
```go
err := goleakwatch.Check(func() {
    // Your code here
    go someFunction()
}, goleakwatch.Config{
    Threshold:   1,
    Wait:        200 * time.Millisecond,
    EnableTrace: true,
})
```

### `CheckWithContext(ctx context.Context, fn func(), cfg Config) error`

Runs the given function with context support for cancellation and timeout control.

**Parameters:**
- `ctx`: Context for cancellation and timeout
- `fn`: The function to execute and monitor
- `cfg`: Configuration options for the leak detection

**Returns:**
- `error`: Returns an error if goroutine leaks are detected or context is cancelled

### `DefaultCheck(fn func()) error`

Convenience function that runs leak check with sensible defaults.

**Parameters:**
- `fn`: The function to execute and monitor

**Returns:**
- `error`: Returns an error if goroutine leaks are detected

### `WithTest(t interface{ Errorf(string, ...interface{}) }, fn func())`

Test helper that integrates with Go's testing framework.

**Parameters:**
- `t`: Testing interface (usually `*testing.T`)
- `fn`: The function to execute and monitor

## Configuration

### `Config` struct

```go
type Config struct {
    Threshold   int           // Max allowed goroutine difference
    Wait        time.Duration // Wait time after function runs
    EnableTrace bool          // Dump goroutine trace if leak suspected
    Out         io.Writer     // Where to write pprof dump (default: os.Stderr)
    Timeout     time.Duration // Timeout for the entire check operation
}
```

**Fields:**
- `Threshold`: Maximum number of additional goroutines allowed (default: 1)
- `Wait`: Time to wait after function execution before counting goroutines (default: 200ms)
- `EnableTrace`: Whether to dump goroutine stack traces when leaks are detected (default: true)
- `Out`: Writer for stack trace output (default: os.Stderr)
- `Timeout`: Maximum time for the entire check operation (default: 5s)

## Benchmarking

### `Benchmark(name string, fn func()) BenchmarkResult`

Measures performance and goroutine usage of a function.

**Parameters:**
- `name`: Name of the function being benchmarked
- `fn`: The function to benchmark

**Returns:**
- `BenchmarkResult`: Detailed benchmark results

### `BenchmarkResult` struct

```go
type BenchmarkResult struct {
    FunctionName     string
    BeforeGoroutines int
    AfterGoroutines  int
    LeakCount        int
    ExecutionTime    time.Duration
    MemoryUsage      uint64
}
```

## Best Practices

1. **Use in Tests**: Integrate leak detection into your unit tests
2. **Set Appropriate Thresholds**: Allow for some variance in goroutine counts
3. **Enable Tracing**: Use stack traces to debug leaks in development
4. **Use Context**: Use `CheckWithContext` for long-running operations
5. **Benchmark Regularly**: Use the benchmarking utilities to track performance

## Examples

See the `examples/` directory for complete working examples. 