package goleakwatch

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

// Config holds configuration for the leak checker
type Config struct {
	Threshold   int           // Max allowed goroutine difference
	Wait        time.Duration // Wait time after function runs
	EnableTrace bool          // Dump goroutine trace if leak suspected
	Out         io.Writer     // Where to write pprof dump (default: os.Stderr)
	Timeout     time.Duration // Timeout for the entire check operation
}

// Check runs the given function and checks for goroutine leaks
func Check(fn func(), cfg Config) error {
	return CheckWithContext(context.Background(), fn, cfg)
}

// CheckWithContext runs the given function with context and checks for goroutine leaks
func CheckWithContext(ctx context.Context, fn func(), cfg Config) error {
	if ctx == nil {
		ctx = context.Background()
	}

	before := runtime.NumGoroutine()

	// Run function in a goroutine so we can control it with context
	done := make(chan struct{})
	go func() {
		defer close(done)
		fn()
	}()

	// Wait for function completion or timeout
	select {
	case <-done:
		// Function completed normally
	case <-ctx.Done():
		// Context was cancelled or timed out
		return fmt.Errorf("leak check cancelled: %v", ctx.Err())
	case <-time.After(cfg.Wait):
		// Wait time exceeded
	}

	after := runtime.NumGoroutine()

	diff := after - before
	if diff > cfg.Threshold {
		msg := fmt.Sprintf("\u26a0\ufe0f Potential goroutine leak: +%d goroutines (before: %d, after: %d)", diff, before, after)
		if cfg.EnableTrace {
			traceOut := cfg.Out
			if traceOut == nil {
				traceOut = os.Stderr
			}
			fmt.Fprintf(traceOut, "\nDumping goroutine stack trace:\n")
			pprof.Lookup("goroutine").WriteTo(traceOut, 2)
		}
		return fmt.Errorf(msg)
	}
	return nil
}

// DefaultCheck runs leak check with sane defaults
func DefaultCheck(fn func()) error {
	return Check(fn, Config{
		Threshold:   1,
		Wait:        200 * time.Millisecond,
		EnableTrace: true,
		Out:         os.Stderr,
		Timeout:     5 * time.Second,
	})
}

// WithTest wraps test logic and reports errors via t.Errorf
func WithTest(t interface{ Errorf(string, ...interface{}) }, fn func()) {
	err := DefaultCheck(fn)
	if err != nil {
		t.Errorf(err.Error())
	}
}
