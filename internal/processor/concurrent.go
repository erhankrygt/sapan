// Package processor provides concurrent stock processing functionality for the SAPAN strategy
// This package handles parallel processing of multiple stocks with worker pools and progress tracking
package processor

import (
	"fmt"
	"log"
	"sapan/internal/data"
	"sapan/internal/strategy"
	"sapan/internal/watcher"
	"sapan/models"
	"sync"
	"time"
)

// StockProcessor handles concurrent stock processing with worker pools
// This struct manages parallel processing of multiple stocks using goroutines and channels
type StockProcessor struct {
	stockFetcher     *data.StockDataFetcher    // Data fetcher for retrieving stock information
	sapanStrategy    *strategy.SAPANStrategy   // SAPAN strategy for validation
	watchListManager *watcher.WatchListManager // Watch list manager for storing results
	workerCount      int                       // Number of concurrent workers
	requestDelay     time.Duration             // Delay between API requests per worker
}

// NewStockProcessor creates a new stock processor instance
// This constructor initializes the processor with all required dependencies and configuration
func NewStockProcessor(
	stockFetcher *data.StockDataFetcher,
	sapanStrategy *strategy.SAPANStrategy,
	watchListManager *watcher.WatchListManager,
	workerCount int,
	requestDelay time.Duration,
) *StockProcessor {
	return &StockProcessor{
		stockFetcher:     stockFetcher,     // Initialize data fetcher
		sapanStrategy:    sapanStrategy,    // Initialize SAPAN strategy
		watchListManager: watchListManager, // Initialize watch list manager
		workerCount:      workerCount,      // Set worker count
		requestDelay:     requestDelay,     // Set request delay
	}
}

// ProcessingResult contains the result of processing a single stock
// This structure holds all information about the processing outcome for a single stock
type ProcessingResult struct {
	Symbol       string // Stock symbol that was processed
	Success      bool   // Whether the processing was successful (no errors)
	Error        error  // Error that occurred during processing (if any)
	IsValid      bool   // Whether any valid SAPAN setup was found
	IsLongValid  bool   // Whether a valid Long setup was found
	IsShortValid bool   // Whether a valid Short setup was found
	Message      string // Detailed message about the processing result
	Processed    bool   // Whether the stock was actually processed
}

// ProcessStocksConcurrently processes multiple stocks concurrently using worker pools
// This method creates channels, starts workers, and coordinates the processing of all stocks
func (p *StockProcessor) ProcessStocksConcurrently(stocks []models.Stock) {
	// Create channels for communication
	stockChan := make(chan models.Stock, len(stocks))
	resultChan := make(chan ProcessingResult, len(stocks))

	// Create progress tracker
	progressTracker := NewProgressTracker(len(stocks))

	// Start progress monitor
	go p.monitorProgress(progressTracker)

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < p.workerCount; i++ {
		wg.Add(1)
		go p.worker(i, stockChan, resultChan, progressTracker, &wg)
	}

	// Send stocks to workers
	go func() {
		defer close(stockChan)
		for _, stock := range stocks {
			stockChan <- stock
		}
	}()

	// Close result channel when all workers are done
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results
	p.collectResults(resultChan, progressTracker)
}

// worker processes stocks from the input channel
func (p *StockProcessor) worker(workerID int, stockChan <-chan models.Stock, resultChan chan<- ProcessingResult, progressTracker *ProgressTracker, wg *sync.WaitGroup) {
	defer wg.Done()

	for stock := range stockChan {
		result := p.processStock(stock)
		resultChan <- result

		// Update progress
		progressTracker.UpdateProgress(result.Success, result.IsValid)

		// Add delay between requests to respect API limits
		if p.requestDelay > 0 {
			time.Sleep(p.requestDelay)
		}
	}
}

// processStock processes a single stock
func (p *StockProcessor) processStock(stock models.Stock) ProcessingResult {
	result := ProcessingResult{
		Symbol:    stock.Symbol,
		Processed: true,
	}

	// Fetch stock data
	candleData, err := p.stockFetcher.FetchStockData(stock.Symbol, 200)
	if err != nil {
		result.Error = err
		result.Success = false
		log.Printf("Worker: Failed to fetch data for %s: %v", stock.Symbol, err)
		return result
	}

	// Validate SAPAN Long strategy first (priority)
	longResult := p.sapanStrategy.ValidateLongSetup(stock.Symbol, candleData.Candles)

	// Validate SAPAN Short strategy only if Long is not valid
	var shortResult strategy.ValidationResult
	if !longResult.IsValid {
		shortResult = p.sapanStrategy.ValidateShortSetup(stock.Symbol, candleData.Candles)
	}

	// Set results based on priority (Long has priority over Short)
	result.IsLongValid = longResult.IsValid
	result.IsShortValid = !longResult.IsValid && shortResult.IsValid
	result.Success = true
	result.IsValid = longResult.IsValid || shortResult.IsValid

	// Create message based on selected scenario
	if longResult.IsValid {
		result.Message = longResult.ValidationMessage
		// Add to Long watch list only
		p.watchListManager.AddToLongWatchList(stock.Symbol)
	} else if shortResult.IsValid {
		result.Message = shortResult.ValidationMessage
		// Add to Short watch list only
		p.watchListManager.AddToShortWatchList(stock.Symbol)
	} else {
		result.Message = "No valid SAPAN setups detected"
	}

	return result
}

// collectResults collects and processes results from workers
func (p *StockProcessor) collectResults(resultChan <-chan ProcessingResult, progressTracker *ProgressTracker) {
	successCount := 0
	errorCount := 0
	validCount := 0
	longCount := 0
	shortCount := 0

	log.Println("Processing results...")

	for result := range resultChan {
		if result.Success {
			successCount++
			if result.IsValid {
				validCount++
			}
			if result.IsLongValid {
				longCount++
			}
			if result.IsShortValid {
				shortCount++
			}
		} else {
			errorCount++
		}

		// Log detailed results
		if result.Success {
			if result.IsValid {
				log.Printf("✅ %s: %s", result.Symbol, result.Message)
			} else {
				log.Printf("❌ %s: %s", result.Symbol, result.Message)
			}
		} else {
			log.Printf("⚠️  %s: Error - %v", result.Symbol, result.Error)
		}
	}

	// Print final progress
	fmt.Println() // New line after progress indicator

	// Print summary (Long and Short are mutually exclusive)
	log.Printf("\n📊 Processing Summary:")
	log.Printf("   Total processed: %d", successCount+errorCount)
	log.Printf("   Successful: %d", successCount)
	log.Printf("   Errors: %d", errorCount)
	log.Printf("   Valid SAPAN setups: %d", validCount)
	log.Printf("   Long setups: %d", longCount)
	log.Printf("   Short setups: %d", shortCount)
	log.Printf("   Note: Each stock can only be either Long OR Short (mutually exclusive)")
}

// monitorProgress monitors and displays progress
func (p *StockProcessor) monitorProgress(progressTracker *ProgressTracker) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		progressTracker.PrintProgress()
		if progressTracker.IsComplete() {
			return
		}
	}
}
