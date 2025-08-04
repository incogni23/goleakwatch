package goleakwatch

import (
	"testing"
	"time"
)

func TestNoLeak(t *testing.T) {
	WithTest(t, func() {
		go func() {
			time.Sleep(10 * time.Millisecond)
		}()
	})
}

func TestLeak(t *testing.T) {
	err := Check(func() {
		go func() {
			select {} // leak
		}()
	}, &Config{
		Threshold:   0,
		Wait:        100 * time.Millisecond,
		EnableTrace: false,
	})
	if err == nil {
		t.Errorf("expected leak but got none")
	}
}
