// Package config provides configuration management for the SAPAN strategy application
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all configuration parameters for the SAPAN strategy application
// This structure centralizes all configurable settings to make the application easily tunable
type Config struct {
	APIKey       string        // Alpha Vantage API key for fetching stock data
	APIURL       string        // Alpha Vantage API base URL
	WorkerCount  int           // Number of concurrent workers for processing stocks
	RequestDelay time.Duration // Delay between API requests per worker (to respect rate limits)
	StocksFile   string        // Path to the JSON file containing stock symbols to analyze
	OutputSize   int           // Number of days of historical data to fetch from API
}

// LoadConfig loads configuration from environment variables with fallback defaults
// This function reads environment variables and provides sensible defaults for missing values
func LoadConfig() (*Config, error) {
	config := &Config{}

	// Load API key from environment (required)
	apiKey := os.Getenv("ALPHA_VANTAGE_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("ALPHA_VANTAGE_API_KEY environment variable is required")
	}
	config.APIKey = apiKey

	// Load API URL from environment (optional, default: Alpha Vantage URL)
	apiURL := os.Getenv("ALPHA_VANTAGE_API_URL")
	if apiURL != "" {
		config.APIURL = apiURL
	} else {
		config.APIURL = "https://www.alphavantage.co/query" // Default Alpha Vantage URL
	}

	// Load worker count from environment (optional, default: 5)
	workerCountStr := os.Getenv("WORKER_COUNT")
	if workerCountStr != "" {
		workerCount, err := strconv.Atoi(workerCountStr)
		if err != nil {
			return nil, fmt.Errorf("invalid WORKER_COUNT value: %v", err)
		}
		config.WorkerCount = workerCount
	} else {
		config.WorkerCount = 5 // Default value
	}

	// Load request delay from environment (optional, default: 2 seconds)
	requestDelayStr := os.Getenv("REQUEST_DELAY_SECONDS")
	if requestDelayStr != "" {
		requestDelay, err := strconv.Atoi(requestDelayStr)
		if err != nil {
			return nil, fmt.Errorf("invalid REQUEST_DELAY_SECONDS value: %v", err)
		}
		config.RequestDelay = time.Duration(requestDelay) * time.Second
	} else {
		config.RequestDelay = time.Second * 2 // Default value
	}

	// Load stocks file path from environment (optional, default: dist/Stocks.json)
	stocksFile := os.Getenv("STOCKS_FILE")
	if stocksFile != "" {
		config.StocksFile = stocksFile
	} else {
		config.StocksFile = "dist/Stocks.json" // Default value
	}

	// Load output size from environment (optional, default: 200)
	outputSizeStr := os.Getenv("OUTPUT_SIZE")
	if outputSizeStr != "" {
		outputSize, err := strconv.Atoi(outputSizeStr)
		if err != nil {
			return nil, fmt.Errorf("invalid OUTPUT_SIZE value: %v", err)
		}
		config.OutputSize = outputSize
	} else {
		config.OutputSize = 200 // Default value
	}

	return config, nil
}

// GetOptimalWorkerCount calculates the optimal number of workers based on request delay
// This method ensures we don't exceed API rate limits while maximizing throughput
// With 2 second delay, 5 workers = 1 request every 0.4 seconds, which is safe for most APIs
func (c *Config) GetOptimalWorkerCount() int {
	// Cap the worker count to prevent overwhelming the API
	if c.WorkerCount > 10 {
		return 10 // Maximum 10 workers to stay within reasonable limits
	}
	if c.WorkerCount < 1 {
		return 1 // Minimum 1 worker to ensure processing continues
	}
	return c.WorkerCount // Return the configured worker count if within bounds
}
