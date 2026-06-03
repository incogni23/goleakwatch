package goleakwatch

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/incogni23/goleakwatch/internal/errors"
	"github.com/incogni23/goleakwatch/internal/logger"
	"github.com/incogni23/goleakwatch/internal/snapshot"
	"github.com/incogni23/goleakwatch/internal/utils"
)

// Config holds configuration for the leak checker
type Config struct {
	Threshold    int           // Max allowed goroutine difference
	Wait         time.Duration // Wait time after function runs
	EnableTrace  bool          // Dump goroutine trace if leak suspected
	Out          io.Writer     // Where to write pprof dump (default: os.Stderr)
	Timeout      time.Duration // Timeout for the entire check operation
	Logger       logger.Logger // Custom logger (optional)
	FunctionName string        // Name of the function being tested
}

// Check runs the given function and checks for goroutine leaks
func Check(fn func(), cfg *Config) error {
	return CheckWithContext(context.Background(), fn, cfg)
}

// CheckWithContext runs the given function with context and checks for goroutine leaks
func CheckWithContext(ctx context.Context, fn func(), cfg *Config) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if cfg == nil {
		cfg = &Config{
			Threshold: 1,
			Wait:      200 * time.Millisecond,
		}
	}

	// Apply Timeout from config if set and context has no deadline
	if cfg.Timeout > 0 {
		if _, hasDeadline := ctx.Deadline(); !hasDeadline {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, cfg.Timeout)
			defer cancel()
		}
	}

	// Use custom logger or default
	log := cfg.Logger
	if log == nil {
		log = logger.GetGlobalLogger()
	}

	// Create timer for performance tracking
	timer := utils.NewTimer()

	before := runtime.NumGoroutine()
	log.Debug("Starting leak check", logger.F("before_goroutines", before), logger.F("function", cfg.FunctionName))

	// Run function in a goroutine so we can control it with context
	done := make(chan struct{})
	go func() {
		defer close(done)
		fn()
	}()

	// Wait for function completion or timeout
	select {
	case <-done:
		log.Debug("Function completed normally")
	case <-ctx.Done():
		log.Warn("Leak check cancelled", logger.F("error", ctx.Err()))
		return fmt.Errorf("leak check cancelled: %v", ctx.Err())
	case <-time.After(cfg.Wait):
		log.Debug("Wait time exceeded")
	}

	after := runtime.NumGoroutine()
	elapsed := timer.Elapsed()

	log.Debug("Leak check completed",
		logger.F("after_goroutines", after),
		logger.F("elapsed", elapsed),
		logger.F("function", cfg.FunctionName))

	diff := after - before
	if diff > cfg.Threshold {
		// Capture stack trace if enabled
		var stackTrace string
		if cfg.EnableTrace {
			traceOut := cfg.Out
			if traceOut == nil {
				traceOut = os.Stderr
			}
			fmt.Fprintf(traceOut, "\nDumping goroutine stack trace:\n")
			if err := pprof.Lookup("goroutine").WriteTo(traceOut, 2); err != nil {
				log.Warn("Failed to write stack trace", logger.F("error", err))
			}

			// Capture stack trace as string for error
			var buf strings.Builder
			if err := pprof.Lookup("goroutine").WriteTo(&buf, 2); err != nil {
				log.Warn("Failed to capture stack trace", logger.F("error", err))
			} else {
				stackTrace = buf.String()
			}
		}

		// Create detailed error
		leakErr := errors.NewLeakError(before, after, cfg.Threshold, cfg.Wait, stackTrace, cfg.FunctionName)
		_ = leakErr.WithInfo("elapsed_time", elapsed)
		_ = leakErr.WithInfo("context_cancelled", ctx.Err())

		log.Error("Goroutine leak detected",
			logger.F("leak_count", diff),
			logger.F("threshold", cfg.Threshold),
			logger.F("function", cfg.FunctionName))

		return leakErr
	}

	log.Info("No leaks detected",
		logger.F("goroutine_delta", diff),
		logger.F("threshold", cfg.Threshold),
		logger.F("function", cfg.FunctionName))

	return nil
}

// DefaultCheck runs leak check with sane defaults
func DefaultCheck(fn func()) error {
	cfg := &Config{
		Threshold:   1,
		Wait:        200 * time.Millisecond,
		EnableTrace: true,
		Out:         os.Stderr,
		Timeout:     5 * time.Second,
		Logger:      logger.GetGlobalLogger(),
	}
	return Check(fn, cfg)
}

// WithTest wraps test logic and reports errors via t.Errorf
func WithTest(t interface{ Errorf(string, ...interface{}) }, fn func()) {
	err := DefaultCheck(fn)
	if err != nil {
		t.Errorf(err.Error())
	}
}

// SnapshotCheck uses the snapshot system for more detailed analysis
func SnapshotCheck(fn func(), cfg *Config) error {
	manager := snapshot.NewSnapshotManager()

	// Take before snapshot
	manager.TakeSnapshot("before")

	// Run the function
	err := Check(fn, cfg)

	// Take after snapshot
	manager.TakeSnapshot("after")

	// Compare snapshots
	diff, compareErr := manager.CompareSnapshots("before", "after")
	if compareErr != nil {
		return fmt.Errorf("failed to compare snapshots: %v", compareErr)
	}

	// Log snapshot comparison
	log := cfg.Logger
	if log == nil {
		log = logger.GetGlobalLogger()
	}

	log.Info("Snapshot comparison",
		logger.F("diff", diff.String()),
		logger.F("is_leak", diff.IsLeak(cfg.Threshold)))

	return err
}

// SetLogger sets the global logger for the package
func SetLogger(l logger.Logger) {
	logger.SetGlobalLogger(l)
}

// GetLogger returns the current global logger
func GetLogger() logger.Logger {
	return logger.GetGlobalLogger()
}
