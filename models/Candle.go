// Package models contains data structures for stock and candlestick data
package models

import "time"

// Candle represents a single candlestick with OHLCV data
// This structure stores price and volume information for a specific time period
type Candle struct {
	Date   time.Time `json:"date"`   // Trading date for this candlestick
	Open   float64   `json:"open"`   // Opening price at the start of the period
	High   float64   `json:"high"`   // Highest price reached during the period
	Low    float64   `json:"low"`    // Lowest price reached during the period
	Close  float64   `json:"close"`  // Closing price at the end of the period
	Volume int64     `json:"volume"` // Total volume traded during the period
}

// CandleData represents a collection of candlesticks for analysis
// This structure is used to store multiple candlesticks for a single stock
type CandleData struct {
	Candles []Candle `json:"Candles"` // Array of candlesticks sorted by date (ascending)
}

// CandleResponse represents the raw API response from Alpha Vantage
// This structure is used to parse the JSON response before converting to CandleData
type CandleResponse struct {
	// MetaData contains information about the API response
	MetaData struct {
		Information   string `json:"1. Information"`    // Description of the data
		Symbol        string `json:"2. Symbol"`         // Stock symbol
		LastRefreshed string `json:"3. Last Refreshed"` // Last update timestamp
		OutputSize    string `json:"4. Output Size"`    // Size of the dataset
		TimeZone      string `json:"5. Time Zone"`      // Timezone information
	} `json:"Meta Data"`

	// TimeSeries contains the actual OHLCV data as strings from the API
	// Keys are date strings, values contain the price and volume data
	TimeSeries map[string]struct {
		Open   string `json:"1. open"`   // Opening price as string
		High   string `json:"2. high"`   // High price as string
		Low    string `json:"3. low"`    // Low price as string
		Close  string `json:"4. close"`  // Close price as string
		Volume string `json:"5. volume"` // Volume as string
	} `json:"Time Series (Daily)"`
}
