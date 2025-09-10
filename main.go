// Package main provides the main entry point for the SAPAN trading strategy application
// This application processes stocks concurrently to identify Long and Short trading setups
package main

import (
	"log"
	"sapan/internal/config"
	"sapan/internal/data"
	"sapan/internal/processor"
	"sapan/internal/strategy"
	"sapan/internal/watcher"
	"time"
)

// main is the entry point of the SAPAN trading strategy application
// This function initializes all components, loads stock data, and processes stocks concurrently
func main() {
	// Load configuration from environment variables
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize all required components using dependency injection
	stockFetcher := data.NewStockDataFetcher(cfg.APIKey, cfg.APIURL) // Initialize data fetcher with API key and URL
	stockLoader := data.NewStockListLoader()                         // Initialize stock list loader
	watchListManager := watcher.NewWatchListManager()                // Initialize watch list manager
	sapanStrategy := strategy.NewSAPANStrategy()                     // Initialize SAPAN strategy

	// Load stock list
	log.Println("📈 Loading stock list...")
	stockData, err := stockLoader.LoadStocksFromFile(cfg.StocksFile)
	if err != nil {
		log.Fatal("Failed to load stocks:", err)
	}

	log.Printf("📊 Loaded %d stocks for analysis", len(stockData.Stocks))

	// Create concurrent processor
	stockProcessor := processor.NewStockProcessor(
		stockFetcher,
		sapanStrategy,
		watchListManager,
		cfg.GetOptimalWorkerCount(),
		cfg.RequestDelay,
	)

	// Process stocks concurrently
	log.Printf("🚀 Starting concurrent processing with %d workers...", cfg.GetOptimalWorkerCount())
	startTime := time.Now()

	stockProcessor.ProcessStocksConcurrently(stockData.Stocks)

	processingTime := time.Since(startTime)
	log.Printf("⏱️  Total processing time: %v", processingTime)

	// Print final results
	log.Println("\n🎯 Final Results:")
	watchListManager.PrintWatchList()

	log.Println("\n✅ SAPAN Strategy analysis completed!")
	time.Sleep(time.Minute * 1)
}
