// Package indicators provides technical analysis indicators for the SAPAN strategy
package indicators

// MACDCalculator handles Moving Average Convergence Divergence (MACD) calculations
// MACD is a trend-following momentum indicator that shows the relationship between two EMAs
type MACDCalculator struct {
	emaCalculator *EMACalculator // EMA calculator for computing fast and slow EMAs
}

// NewMACDCalculator creates a new MACD calculator instance
// This constructor initializes the calculator with an EMA calculator for internal use
func NewMACDCalculator() *MACDCalculator {
	return &MACDCalculator{
		emaCalculator: NewEMACalculator(), // Initialize EMA calculator
	}
}

// MACDResult contains the result of MACD calculation
// This structure holds the MACD line, Signal line, and Histogram values
type MACDResult struct {
	MACD      float64 // MACD line (fast EMA - slow EMA)
	Signal    float64 // Signal line (EMA of MACD line)
	Histogram float64 // MACD histogram (MACD - Signal)
}

// Calculate calculates MACD for given prices and periods
// MACD formula: MACD = Fast EMA - Slow EMA
// Signal line is typically a 9-period EMA of the MACD line
// Histogram = MACD - Signal line
func (m *MACDCalculator) Calculate(prices []float64, fastPeriod, slowPeriod, signalPeriod int) MACDResult {
	if len(prices) < slowPeriod {
		return MACDResult{}
	}

	// Calculate EMAs
	emaFast := m.emaCalculator.Calculate(prices, fastPeriod)
	emaSlow := m.emaCalculator.Calculate(prices, slowPeriod)
	macd := emaFast - emaSlow

	// Calculate MACD line values over time for signal line
	macdValues := make([]float64, 0)
	for i := slowPeriod; i < len(prices); i++ {
		fastEMA := m.emaCalculator.Calculate(prices[:i+1], fastPeriod)
		slowEMA := m.emaCalculator.Calculate(prices[:i+1], slowPeriod)
		macdValue := fastEMA - slowEMA
		macdValues = append(macdValues, macdValue)
	}

	// Calculate signal line (EMA of MACD values)
	var signal float64
	if len(macdValues) >= signalPeriod {
		signal = m.emaCalculator.Calculate(macdValues, signalPeriod)
	} else {
		signal = macd * 0.9 // Fallback
	}

	histogram := macd - signal

	return MACDResult{
		MACD:      macd,
		Signal:    signal,
		Histogram: histogram,
	}
}

// IsBullMarket checks if MACD is in bull market
// IsBullMarket checks if MACD indicates a bull market
// Returns true if MACD line is above the Signal line, indicating bullish momentum
func (m *MACDCalculator) IsBullMarket(prices []float64, fastPeriod, slowPeriod, signalPeriod int) bool {
	result := m.Calculate(prices, fastPeriod, slowPeriod, signalPeriod)
	return result.MACD > result.Signal // Bull market when MACD > Signal
}

// IsBearMarketAcceptable checks if bear market duration is acceptable (≤ 5 candlesticks)
func (m *MACDCalculator) IsBearMarketAcceptable(prices []float64, fastPeriod, slowPeriod, signalPeriod int) bool {
	result := m.Calculate(prices, fastPeriod, slowPeriod, signalPeriod)

	// If in bull market, it's acceptable
	if result.MACD > result.Signal {
		return true
	}

	// Bear market - check duration
	bearishCount := 0
	for j := len(prices) - 1; j >= 0 && bearishCount < 6; j-- {
		if j < 1 {
			break
		}
		// Calculate MACD for this point
		subCloses := prices[:j+1]
		if len(subCloses) >= slowPeriod {
			subResult := m.Calculate(subCloses, fastPeriod, slowPeriod, signalPeriod)
			if subResult.MACD <= subResult.Signal {
				bearishCount++
			} else {
				break
			}
		}
	}

	// If bearish for 5 or fewer candlesticks, it's acceptable
	return bearishCount <= 5
}

// IsBullMarketAcceptable checks if bull market duration is acceptable (≤ 5 candlesticks)
func (m *MACDCalculator) IsBullMarketAcceptable(prices []float64, fastPeriod, slowPeriod, signalPeriod int) bool {
	result := m.Calculate(prices, fastPeriod, slowPeriod, signalPeriod)

	// If in bear market, it's acceptable
	if result.MACD < result.Signal {
		return true
	}

	// Bull market - check duration
	bullishCount := 0
	for j := len(prices) - 1; j >= 0 && bullishCount < 6; j-- {
		if j < 1 {
			break
		}
		// Calculate MACD for this point
		subCloses := prices[:j+1]
		if len(subCloses) >= slowPeriod {
			subResult := m.Calculate(subCloses, fastPeriod, slowPeriod, signalPeriod)
			if subResult.MACD >= subResult.Signal {
				bullishCount++
			} else {
				break
			}
		}
	}

	// If bullish for 5 or fewer candlesticks, it's acceptable
	return bullishCount <= 5
}
