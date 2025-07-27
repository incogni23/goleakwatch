package errors

import (
	"fmt"
	"strings"
	"time"
)

// LeakError represents a goroutine leak with detailed information
type LeakError struct {
	BeforeCount    int                    // Number of goroutines before
	AfterCount     int                    // Number of goroutines after
	LeakCount      int                    // Number of leaked goroutines
	Threshold      int                    // Maximum allowed difference
	WaitDuration   time.Duration          // Time waited after function
	StackTrace     string                 // Goroutine stack trace
	Timestamp      time.Time              // When the leak was detected
	FunctionName   string                 // Name of the function being tested
	AdditionalInfo map[string]interface{} // Additional context
}

// Error implements the error interface
func (e *LeakError) Error() string {
	var parts []string

	parts = append(parts, fmt.Sprintf("🚨 Goroutine leak detected: +%d goroutines", e.LeakCount))
	parts = append(parts, fmt.Sprintf("Before: %d, After: %d (threshold: %d)", e.BeforeCount, e.AfterCount, e.Threshold))

	if e.FunctionName != "" {
		parts = append(parts, fmt.Sprintf("Function: %s", e.FunctionName))
	}

	parts = append(parts, fmt.Sprintf("Wait time: %v", e.WaitDuration))
	parts = append(parts, fmt.Sprintf("Detected at: %v", e.Timestamp.Format(time.RFC3339)))

	if e.StackTrace != "" {
		parts = append(parts, "\nStack trace:")
		parts = append(parts, e.StackTrace)
	}

	return strings.Join(parts, "\n")
}

// IsLeakError checks if an error is a LeakError
func IsLeakError(err error) bool {
	_, ok := err.(*LeakError)
	return ok
}

// GetLeakError extracts LeakError from an error
func GetLeakError(err error) (*LeakError, bool) {
	if leakErr, ok := err.(*LeakError); ok {
		return leakErr, true
	}
	return nil, false
}

// NewLeakError creates a new LeakError with current timestamp
func NewLeakError(before, after, threshold int, waitDuration time.Duration, stackTrace, functionName string) *LeakError {
	return &LeakError{
		BeforeCount:    before,
		AfterCount:     after,
		LeakCount:      after - before,
		Threshold:      threshold,
		WaitDuration:   waitDuration,
		StackTrace:     stackTrace,
		Timestamp:      time.Now(),
		FunctionName:   functionName,
		AdditionalInfo: make(map[string]interface{}),
	}
}

// WithInfo adds additional context to the error
func (e *LeakError) WithInfo(key string, value interface{}) *LeakError {
	e.AdditionalInfo[key] = value
	return e
}

// GetInfo retrieves additional context
func (e *LeakError) GetInfo(key string) (interface{}, bool) {
	value, exists := e.AdditionalInfo[key]
	return value, exists
}

// Summary returns a brief summary of the leak
func (e *LeakError) Summary() string {
	return fmt.Sprintf("Leak: +%d goroutines (threshold: %d)", e.LeakCount, e.Threshold)
}

// IsSignificant returns true if the leak exceeds threshold significantly
func (e *LeakError) IsSignificant(multiplier float64) bool {
	return float64(e.LeakCount) > float64(e.Threshold)*multiplier
}
