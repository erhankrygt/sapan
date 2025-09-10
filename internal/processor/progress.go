// Package processor provides concurrent stock processing functionality for the SAPAN strategy
// This package handles parallel processing of multiple stocks with worker pools and progress tracking
package processor

import (
	"fmt"
	"sync/atomic"
	"time"
)

// ProgressTracker tracks progress of concurrent processing
// This struct provides thread-safe progress tracking using atomic operations
type ProgressTracker struct {
	total     int32     // Total number of items to process
	processed int32     // Number of items processed so far
	valid     int32     // Number of valid SAPAN setups found
	errors    int32     // Number of errors encountered
	startTime time.Time // Start time for calculating elapsed time
}

// NewProgressTracker creates a new progress tracker instance
// This constructor initializes the tracker with the total number of items to process
func NewProgressTracker(total int) *ProgressTracker {
	return &ProgressTracker{
		total:     int32(total), // Set the total number of items
		processed: 0,            // Initialize processed count
		valid:     0,            // Initialize valid count
		errors:    0,            // Initialize error count
		startTime: time.Now(),   // Record start time
	}
}

// UpdateProgress updates the progress counters atomically
// This method is thread-safe and can be called from multiple goroutines
func (p *ProgressTracker) UpdateProgress(success, valid bool) {
	atomic.AddInt32(&p.processed, 1) // Increment processed count
	if success {
		if valid {
			atomic.AddInt32(&p.valid, 1) // Increment valid count if setup is valid
		}
	} else {
		atomic.AddInt32(&p.errors, 1) // Increment error count if processing failed
	}
}

// GetProgress returns current progress information atomically
// This method provides thread-safe access to progress counters and calculates percentage
func (p *ProgressTracker) GetProgress() (processed, valid, errors int32, percentage float64) {
	processed = atomic.LoadInt32(&p.processed) // Get current processed count
	valid = atomic.LoadInt32(&p.valid)         // Get current valid count
	errors = atomic.LoadInt32(&p.errors)       // Get current error count

	// Calculate percentage completion
	if p.total > 0 {
		percentage = float64(processed) / float64(p.total) * 100
	}

	return processed, valid, errors, percentage
}

// PrintProgress prints current progress with real-time statistics
// This method displays progress information including percentage, valid setups, errors, and elapsed time
func (p *ProgressTracker) PrintProgress() {
	processed, valid, errors, percentage := p.GetProgress()
	elapsed := time.Since(p.startTime) // Calculate elapsed time

	fmt.Printf("\r🔄 Progress: %d/%d (%.1f%%) | ✅ Valid: %d | ❌ Errors: %d | ⏱️  %v",
		processed, p.total, percentage, valid, errors, elapsed.Round(time.Second))
}

// IsComplete checks if processing is complete
// This method returns true when all items have been processed
func (p *ProgressTracker) IsComplete() bool {
	return atomic.LoadInt32(&p.processed) >= p.total
}
