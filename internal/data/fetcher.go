// Package data provides data fetching and loading functionality for the SAPAN strategy
package data

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sapan/models"
	"sort"
	"strconv"
	"time"
)

// StockDataFetcher handles fetching stock data from external APIs
// This struct encapsulates the API key and URL, providing methods to fetch historical stock data
type StockDataFetcher struct {
	apiKey string // Alpha Vantage API key for authentication
	apiURL string // Alpha Vantage API base URL
}

// NewStockDataFetcher creates a new stock data fetcher with the provided API key and URL
// The API key and URL are required for authenticating requests to the Alpha Vantage API
func NewStockDataFetcher(apiKey, apiURL string) *StockDataFetcher {
	return &StockDataFetcher{
		apiKey: apiKey, // Store the API key for use in HTTP requests
		apiURL: apiURL, // Store the API URL for constructing requests
	}
}

// FetchStockData fetches historical stock data for a given symbol from Alpha Vantage API
// This method constructs the API URL, makes the HTTP request, and processes the response
// Returns CandleData containing sorted candlesticks or an error if the request fails
func (f *StockDataFetcher) FetchStockData(symbol string, outputSize int) (models.CandleData, error) {
	// Construct the API URL with the required parameters using the configured base URL
	url := fmt.Sprintf(
		"%s?function=TIME_SERIES_DAILY&symbol=%s&outputsize=%d&apikey=%s",
		f.apiURL, symbol, outputSize, f.apiKey,
	)

	// Make HTTP GET request to the Alpha Vantage API
	resp, err := http.Get(url)
	if err != nil {
		return models.CandleData{}, fmt.Errorf("failed to fetch data: %v", err)
	}
	defer resp.Body.Close() // Ensure response body is closed

	// Read the entire response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.CandleData{}, fmt.Errorf("failed to read response: %v", err)
	}

	// Parse the JSON response into our CandleResponse structure
	var avResponse models.CandleResponse
	if err = json.Unmarshal(body, &avResponse); err != nil {
		return models.CandleData{}, fmt.Errorf("failed to parse JSON: %v", err)
	}

	// Handle API errors (rate limits, invalid symbols, etc.)
	if len(avResponse.TimeSeries) == 0 {
		var errorResp map[string]interface{}
		if err := json.Unmarshal(body, &errorResp); err == nil {
			// Check for rate limit message
			if note, ok := errorResp["Note"]; ok {
				return models.CandleData{}, fmt.Errorf("API rate limit: %v", note)
			}
			// Check for error message
			if errorMsg, ok := errorResp["Error Message"]; ok {
				return models.CandleData{}, fmt.Errorf("API error: %v", errorMsg)
			}
		}

		return models.CandleData{}, fmt.Errorf("invalid API response")
	}

	// Convert the raw API response to our CandleData structure
	candles := f.convertToCandles(avResponse.TimeSeries)
	return models.CandleData{Candles: candles}, nil
}

// convertToCandles converts the raw API response to our Candle models
// This method parses string values from the API response and converts them to proper data types
// It also sorts the candles by date in ascending order for proper chronological analysis
func (f *StockDataFetcher) convertToCandles(timeSeries map[string]struct {
	Open   string `json:"1. open"`
	High   string `json:"2. high"`
	Low    string `json:"3. low"`
	Close  string `json:"4. close"`
	Volume string `json:"5. volume"`
}) []models.Candle {
	// Pre-allocate slice with capacity to avoid reallocations
	candles := make([]models.Candle, 0, len(timeSeries))

	// Iterate through each date in the time series
	for dateStr, data := range timeSeries {
		// Parse the date string (format: "2006-01-02")
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue // Skip invalid dates
		}

		// Parse opening price from string to float64
		open, err := strconv.ParseFloat(data.Open, 64)
		if err != nil {
			continue // Skip if parsing fails
		}

		// Parse high price from string to float64
		high, err := strconv.ParseFloat(data.High, 64)
		if err != nil {
			continue // Skip if parsing fails
		}

		// Parse low price from string to float64
		low, err := strconv.ParseFloat(data.Low, 64)
		if err != nil {
			continue // Skip if parsing fails
		}

		// Parse closing price from string to float64
		closePrice, err := strconv.ParseFloat(data.Close, 64)
		if err != nil {
			continue // Skip if parsing fails
		}

		// Parse volume from string to int64
		volume, err := strconv.ParseInt(data.Volume, 10, 64)
		if err != nil {
			continue // Skip if parsing fails
		}

		// Create a new Candle with the parsed data
		candles = append(candles, models.Candle{
			Date:   date,       // Trading date
			Open:   open,       // Opening price
			High:   high,       // Highest price
			Low:    low,        // Lowest price
			Close:  closePrice, // Closing price
			Volume: volume,     // Trading volume
		})
	}

	// Sort candles by date in ascending order (oldest first)
	// This is crucial for proper technical analysis as indicators depend on chronological order
	sort.Slice(candles, func(i, j int) bool {
		return candles[i].Date.Before(candles[j].Date)
	})

	return candles
}
