// Package models contains data structures for stock and candlestick data
package models

// Stock represents a single stock with its basic information
// This structure is used to store stock metadata from the stocks.json file
type Stock struct {
	Symbol   string `json:"symbol"`   // Stock ticker symbol (e.g., "AAPL", "GOOGL")
	Name     string `json:"name"`     // Full company name
	Sector   string `json:"sector"`   // Business sector (e.g., "Technology", "Healthcare")
	Industry string `json:"industry"` // Specific industry within the sector
}

// StockData represents a collection of stocks
// This structure is used to parse the entire stocks.json file
type StockData struct {
	Stocks []Stock `json:"Stocks"` // Array of all stocks to be analyzed
}
