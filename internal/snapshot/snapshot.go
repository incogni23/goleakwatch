package snapshot

import (
	"fmt"
	"runtime"
	"runtime/pprof"
	"strings"
	"sync"
	"time"
)

// Snapshot represents a point-in-time capture of goroutine state
type Snapshot struct {
	ID             string                 // Unique identifier
	Timestamp      time.Time              // When the snapshot was taken
	GoroutineCount int                    // Number of goroutines
	StackTraces    string                 // Full goroutine stack traces
	Metadata       map[string]interface{} // Additional metadata
	mu             sync.RWMutex           // Protects metadata
}

// SnapshotManager manages multiple snapshots
type SnapshotManager struct {
	snapshots map[string]*Snapshot
	mu        sync.RWMutex
}

// NewSnapshotManager creates a new snapshot manager
func NewSnapshotManager() *SnapshotManager {
	return &SnapshotManager{
		snapshots: make(map[string]*Snapshot),
	}
}

// TakeSnapshot captures the current goroutine state
func (sm *SnapshotManager) TakeSnapshot(id string) *Snapshot {
	// Capture stack traces
	var buf strings.Builder
	if err := pprof.Lookup("goroutine").WriteTo(&buf, 2); err != nil {
		// Log error but continue with empty stack trace
		buf.WriteString("Failed to capture stack trace: " + err.Error())
	}

	snapshot := &Snapshot{
		ID:             id,
		Timestamp:      time.Now(),
		GoroutineCount: runtime.NumGoroutine(),
		StackTraces:    buf.String(),
		Metadata:       make(map[string]interface{}),
	}

	sm.mu.Lock()
	sm.snapshots[id] = snapshot
	sm.mu.Unlock()

	return snapshot
}

// GetSnapshot retrieves a snapshot by ID
func (sm *SnapshotManager) GetSnapshot(id string) (*Snapshot, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	snapshot, exists := sm.snapshots[id]
	return snapshot, exists
}

// CompareSnapshots compares two snapshots and returns the difference
func (sm *SnapshotManager) CompareSnapshots(beforeID, afterID string) (*SnapshotDiff, error) {
	before, exists := sm.GetSnapshot(beforeID)
	if !exists {
		return nil, fmt.Errorf("snapshot '%s' not found", beforeID)
	}

	after, exists := sm.GetSnapshot(afterID)
	if !exists {
		return nil, fmt.Errorf("snapshot '%s' not found", afterID)
	}

	return &SnapshotDiff{
		Before:         before,
		After:          after,
		GoroutineDelta: after.GoroutineCount - before.GoroutineCount,
		TimeDelta:      after.Timestamp.Sub(before.Timestamp),
	}, nil
}

// SnapshotDiff represents the difference between two snapshots
type SnapshotDiff struct {
	Before         *Snapshot
	After          *Snapshot
	GoroutineDelta int
	TimeDelta      time.Duration
}

// IsLeak returns true if there's a potential leak
func (sd *SnapshotDiff) IsLeak(threshold int) bool {
	return sd.GoroutineDelta > threshold
}

// String returns a formatted string representation
func (sd *SnapshotDiff) String() string {
	return fmt.Sprintf(
		"Goroutines: %d → %d (Δ: %d) | Time: %v",
		sd.Before.GoroutineCount,
		sd.After.GoroutineCount,
		sd.GoroutineDelta,
		sd.TimeDelta,
	)
}

// GetMetadata retrieves metadata from a snapshot
func (s *Snapshot) GetMetadata(key string) (interface{}, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, exists := s.Metadata[key]
	return value, exists
}

// SetMetadata sets metadata on a snapshot
func (s *Snapshot) SetMetadata(key string, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Metadata[key] = value
}

// AutoSnapshot creates snapshots at regular intervals
type AutoSnapshot struct {
	manager  *SnapshotManager
	interval time.Duration
	stopChan chan struct{}
	running  bool
	mu       sync.Mutex
}

// NewAutoSnapshot creates a new auto-snapshot system
func NewAutoSnapshot(manager *SnapshotManager, interval time.Duration) *AutoSnapshot {
	return &AutoSnapshot{
		manager:  manager,
		interval: interval,
		stopChan: make(chan struct{}),
	}
}

// Start begins automatic snapshotting
func (as *AutoSnapshot) Start() {
	as.mu.Lock()
	defer as.mu.Unlock()

	if as.running {
		return
	}

	as.running = true
	go as.run()
}

// Stop stops automatic snapshotting
func (as *AutoSnapshot) Stop() {
	as.mu.Lock()
	defer as.mu.Unlock()

	if !as.running {
		return
	}

	as.running = false
	close(as.stopChan)
}

// run is the internal goroutine for automatic snapshotting
func (as *AutoSnapshot) run() {
	ticker := time.NewTicker(as.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			id := fmt.Sprintf("auto_%d", time.Now().Unix())
			as.manager.TakeSnapshot(id)
		case <-as.stopChan:
			return
		}
	}
}
