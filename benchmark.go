package goleakwatch

import (
	"fmt"
	"runtime"
	"time"
)

// BenchmarkResult holds the results of a benchmark
type BenchmarkResult struct {
	FunctionName     string
	BeforeGoroutines int
	AfterGoroutines  int
	LeakCount        int
	ExecutionTime    time.Duration
	MemoryUsage      uint64
}

// Benchmark runs a function and measures its performance and goroutine usage
func Benchmark(name string, fn func()) BenchmarkResult {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	beforeMem := m.Alloc

	before := runtime.NumGoroutine()
	start := time.Now()

	fn()

	executionTime := time.Since(start)
	after := runtime.NumGoroutine()

	runtime.ReadMemStats(&m)
	afterMem := m.Alloc

	var memUsage uint64
	if afterMem > beforeMem {
		memUsage = afterMem - beforeMem
	}

	return BenchmarkResult{
		FunctionName:     name,
		BeforeGoroutines: before,
		AfterGoroutines:  after,
		LeakCount:        after - before,
		ExecutionTime:    executionTime,
		MemoryUsage:      memUsage,
	}
}

// String returns a formatted string representation of the benchmark result
func (r BenchmarkResult) String() string {
	return fmt.Sprintf(
		"Function: %s | Goroutines: %d → %d (Δ: %d) | Time: %v | Memory: %d bytes",
		r.FunctionName,
		r.BeforeGoroutines,
		r.AfterGoroutines,
		r.LeakCount,
		r.ExecutionTime,
		r.MemoryUsage,
	)
}

// IsLeaking returns true if there are leaked goroutines
func (r BenchmarkResult) IsLeaking() bool {
	return r.LeakCount > 0
}
