// Package data provides data fetching and loading functionality for the SAPAN strategy
package data

import (
	"encoding/json"
	"os"
	"sapan/models"
)

// StockListLoader handles loading stock lists from JSON files
// This struct provides methods to load stock symbols and metadata from local files
type StockListLoader struct{}

// NewStockListLoader creates a new stock list loader instance
// This constructor initializes the loader for reading stock data from files
func NewStockListLoader() *StockListLoader {
	return &StockListLoader{}
}

// LoadStocksFromFile loads stock symbols and metadata from a JSON file
// This method opens the file, uses a JSON decoder to parse the content, and returns StockData
// The file should contain a JSON structure with a "Stocks" array containing stock information
func (l *StockListLoader) LoadStocksFromFile(filename string) (models.StockData, error) {
	// Open the file for reading
	file, err := os.Open(filename)
	if err != nil {
		return models.StockData{}, err // Return empty StockData and error if file cannot be opened
	}
	defer file.Close() // Ensure file is closed when function returns

	// Create a JSON decoder and decode the file content into StockData structure
	var stocks models.StockData
	err = json.NewDecoder(file).Decode(&stocks)
	if err != nil {
		return models.StockData{}, err // Return empty StockData and error if JSON parsing fails
	}

	// Return the successfully parsed stock data
	return stocks, nil
}
