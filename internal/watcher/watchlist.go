// Package watcher provides watch list management functionality for the SAPAN strategy
// This package handles thread-safe storage and retrieval of trading signals
package watcher

import (
	"fmt"
	"sync"
	"time"
)

// WatchListManager manages the watch list for trading signals
// This struct provides thread-safe operations for storing and retrieving Long and Short trading setups
type WatchListManager struct {
	longWatchList  map[time.Time]string // Map of Long setups with timestamps
	shortWatchList map[time.Time]string // Map of Short setups with timestamps
	mutex          sync.RWMutex         // Read-write mutex for thread-safe operations
}

// NewWatchListManager creates a new watch list manager instance
// This constructor initializes both Long and Short watch lists with thread-safe maps
func NewWatchListManager() *WatchListManager {
	return &WatchListManager{
		longWatchList:  make(map[time.Time]string), // Initialize Long watch list
		shortWatchList: make(map[time.Time]string), // Initialize Short watch list
	}
}

// AddToLongWatchList adds a symbol to the long watch list (thread-safe)
// This method stores a Long trading setup with the current timestamp
func (w *WatchListManager) AddToLongWatchList(symbol string) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.longWatchList[time.Now().UTC()] = symbol // Store with current UTC timestamp
	fmt.Printf("✅ SAPAN Long Setup detected for %s\n", symbol)
}

// GetLongWatchList returns the current long watch list (thread-safe)
// This method returns a copy of the Long watch list to avoid race conditions
func (w *WatchListManager) GetLongWatchList() map[time.Time]string {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	// Create a copy to avoid race conditions
	result := make(map[time.Time]string)
	for timestamp, symbol := range w.longWatchList {
		result[timestamp] = symbol // Copy each entry to the result map
	}
	return result
}

// PrintWatchList prints the current watch list (thread-safe)
// This method displays both Long and Short watch lists with timestamps
func (w *WatchListManager) PrintWatchList() {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	// Print Long Watch List
	fmt.Println("Current Long Watch List:")
	if len(w.longWatchList) == 0 {
		fmt.Println("  No valid SAPAN long setups found")
	} else {
		for timestamp, symbol := range w.longWatchList {
			fmt.Printf("  %s: %s\n", timestamp.Format("2006-01-02 15:04:05"), symbol)
		}
	}

	// Print Short Watch List
	fmt.Println("\nCurrent Short Watch List:")
	if len(w.shortWatchList) == 0 {
		fmt.Println("  No valid SAPAN short setups found")
	} else {
		for timestamp, symbol := range w.shortWatchList {
			fmt.Printf("  %s: %s\n", timestamp.Format("2006-01-02 15:04:05"), symbol)
		}
	}
}

// AddToShortWatchList adds a symbol to the short watch list (thread-safe)
// This method stores a Short trading setup with the current timestamp
func (w *WatchListManager) AddToShortWatchList(symbol string) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.shortWatchList[time.Now().UTC()] = symbol // Store with current UTC timestamp
	fmt.Printf("✅ SAPAN Short Setup detected for %s\n", symbol)
}

// GetShortWatchList returns the current short watch list (thread-safe)
// This method returns a copy of the Short watch list to avoid race conditions
func (w *WatchListManager) GetShortWatchList() map[time.Time]string {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	// Create a copy to avoid race conditions
	result := make(map[time.Time]string)
	for timestamp, symbol := range w.shortWatchList {
		result[timestamp] = symbol // Copy each entry to the result map
	}
	return result
}

// GetCount returns the total number of items in both watch lists (thread-safe)
// This method provides the combined count of Long and Short setups
func (w *WatchListManager) GetCount() int {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return len(w.longWatchList) + len(w.shortWatchList) // Total count of all setups
}

// GetLongCount returns the number of long items in the watch list (thread-safe)
// This method provides the count of Long trading setups
func (w *WatchListManager) GetLongCount() int {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return len(w.longWatchList) // Count of Long setups
}

// GetShortCount returns the number of short items in the watch list (thread-safe)
// This method provides the count of Short trading setups
func (w *WatchListManager) GetShortCount() int {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return len(w.shortWatchList) // Count of Short setups
}
