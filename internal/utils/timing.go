package utils

import (
	"time"
)

// Timer provides utility functions for timing operations
type Timer struct {
	start time.Time
}

// NewTimer creates a new timer
func NewTimer() *Timer {
	return &Timer{
		start: time.Now(),
	}
}

// Elapsed returns the elapsed time since the timer was created
func (t *Timer) Elapsed() time.Duration {
	return time.Since(t.start)
}

// Reset resets the timer
func (t *Timer) Reset() {
	t.start = time.Now()
}
