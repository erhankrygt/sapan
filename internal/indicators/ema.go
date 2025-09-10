// Package indicators provides technical analysis indicators for the SAPAN strategy
package indicators

// EMACalculator handles Exponential Moving Average (EMA) calculations
// EMA gives more weight to recent prices, making it more responsive to price changes than SMA
type EMACalculator struct{}

// NewEMACalculator creates a new EMA calculator instance
// This constructor initializes the calculator for performing EMA calculations
func NewEMACalculator() *EMACalculator {
	return &EMACalculator{}
}

// Calculate calculates the Exponential Moving Average for given prices and period
// EMA formula: EMA = (Price * Multiplier) + (Previous EMA * (1 - Multiplier))
// where Multiplier = 2 / (Period + 1)
// Returns 0 if there's insufficient data for the specified period
func (e *EMACalculator) Calculate(prices []float64, period int) float64 {
	// Check if we have enough data points for the specified period
	if len(prices) < period {
		return 0 // Return 0 if insufficient data
	}

	// Calculate the smoothing factor (multiplier) for EMA
	// This determines how much weight to give to recent prices vs historical EMA
	multiplier := 2.0 / (float64(period) + 1.0)

	// Start with Simple Moving Average (SMA) for the first EMA value
	// This provides a stable starting point for the EMA calculation
	sum := 0.0
	for i := 0; i < period; i++ {
		sum += prices[i] // Sum the first 'period' prices
	}
	ema := sum / float64(period) // Calculate SMA as initial EMA value

	// Calculate EMA for remaining values using the standard EMA formula
	for i := period; i < len(prices); i++ {
		// EMA formula: new EMA = (current price * multiplier) + (previous EMA * (1 - multiplier))
		ema = (prices[i] * multiplier) + (ema * (1 - multiplier))
	}

	return ema
}

// ValidateTrend validates if EMAs are in uptrend order (20 > 50 > 100 > 200)
// This method checks if shorter-term EMAs are above longer-term EMAs, indicating an uptrend
// Used for Long scenario validation in the SAPAN strategy
func (e *EMACalculator) ValidateTrend(prices []float64) bool {
	// Calculate EMAs for different periods
	ema20 := e.Calculate(prices, 20)   // Short-term EMA (20 periods)
	ema50 := e.Calculate(prices, 50)   // Medium-term EMA (50 periods)
	ema100 := e.Calculate(prices, 100) // Long-term EMA (100 periods)
	ema200 := e.Calculate(prices, 200) // Very long-term EMA (200 periods)

	// Check if EMAs are in proper uptrend order (faster EMAs above slower ones)
	return ema20 > ema50 && ema50 > ema100 && ema100 > ema200
}

// ValidateDowntrend validates if EMAs are in downtrend order (20 < 50 < 100 < 200)
// This method checks if shorter-term EMAs are below longer-term EMAs, indicating a downtrend
// Used for Short scenario validation in the SAPAN strategy
func (e *EMACalculator) ValidateDowntrend(prices []float64) bool {
	// Calculate EMAs for different periods
	ema20 := e.Calculate(prices, 20)   // Short-term EMA (20 periods)
	ema50 := e.Calculate(prices, 50)   // Medium-term EMA (50 periods)
	ema100 := e.Calculate(prices, 100) // Long-term EMA (100 periods)
	ema200 := e.Calculate(prices, 200) // Very long-term EMA (200 periods)

	// Check if EMAs are in proper downtrend order (faster EMAs below slower ones)
	return ema20 < ema50 && ema50 < ema100 && ema100 < ema200
}
