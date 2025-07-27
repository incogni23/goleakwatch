package goleakwatch

import (
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
}

// Check runs the given function and checks for goroutine leaks
func Check(fn func(), cfg Config) error {
	before := runtime.NumGoroutine()
	fn()
	time.Sleep(cfg.Wait)
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
	})
}

// WithTest wraps test logic and reports errors via t.Errorf
func WithTest(t interface{ Errorf(string, ...interface{}) }, fn func()) {
	err := DefaultCheck(fn)
	if err != nil {
		t.Errorf(err.Error())
	}
}
