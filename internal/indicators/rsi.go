// Package indicators provides technical analysis indicators for the SAPAN strategy
package indicators

// RSICalculator handles Relative Strength Index (RSI) calculations
// RSI is a momentum oscillator that measures the speed and magnitude of price changes
// It ranges from 0 to 100 and is used to identify overbought and oversold conditions
type RSICalculator struct{}

// NewRSICalculator creates a new RSI calculator instance
// This constructor initializes the calculator for performing RSI calculations
func NewRSICalculator() *RSICalculator {
	return &RSICalculator{}
}

// Calculate calculates the Relative Strength Index for given prices and period
// RSI formula: RSI = 100 - (100 / (1 + RS))
// where RS = Average Gain / Average Loss
// Uses Wilder's smoothing method for more accurate RSI calculation
// Returns 0 if there's insufficient data for the specified period
func (r *RSICalculator) Calculate(prices []float64, period int) float64 {
	// Check if we have enough data points for the specified period
	if len(prices) < period+1 {
		return 0 // Return 0 if insufficient data
	}

	// Calculate price changes and separate gains from losses
	var gains, losses []float64
	for i := 1; i < len(prices); i++ {
		change := prices[i] - prices[i-1] // Calculate price change
		if change > 0 {
			gains = append(gains, change) // Store positive change as gain
			losses = append(losses, 0)    // No loss for this period
		} else {
			gains = append(gains, 0)         // No gain for this period
			losses = append(losses, -change) // Store negative change as loss (make positive)
		}
	}

	// Check if we have enough gain/loss data for the period
	if len(gains) < period {
		return 0 // Return 0 if insufficient data
	}

	// Calculate initial average gain and loss using simple average
	avgGain := 0.0
	avgLoss := 0.0
	for i := 0; i < period; i++ {
		avgGain += gains[i]  // Sum gains for initial period
		avgLoss += losses[i] // Sum losses for initial period
	}
	avgGain /= float64(period) // Calculate average gain
	avgLoss /= float64(period) // Calculate average loss

	// Apply Wilder's smoothing method for more accurate RSI calculation
	// This method gives more weight to recent data while maintaining stability
	for i := period; i < len(gains); i++ {
		// Wilder's smoothing: new_avg = (old_avg * (period-1) + new_value) / period
		avgGain = (avgGain*float64(period-1) + gains[i]) / float64(period)
		avgLoss = (avgLoss*float64(period-1) + losses[i]) / float64(period)
	}

	// Handle edge case where average loss is zero (all gains)
	if avgLoss == 0 {
		return 100 // RSI is 100 when there are no losses
	}

	// Calculate Relative Strength (RS) and RSI
	rs := avgGain / avgLoss       // Relative Strength
	rsi := 100 - (100 / (1 + rs)) // RSI formula

	return rsi
}
