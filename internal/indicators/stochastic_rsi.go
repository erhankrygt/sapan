// Package indicators provides technical analysis indicators for the SAPAN strategy
package indicators

// StochasticRSICalculator handles Stochastic RSI calculations
// Stochastic RSI applies the Stochastic oscillator formula to RSI values instead of prices
// This creates a more sensitive momentum indicator that oscillates between 0 and 100
type StochasticRSICalculator struct {
	rsiCalculator *RSICalculator // RSI calculator for computing RSI values
}

// NewStochasticRSICalculator creates a new Stochastic RSI calculator instance
// This constructor initializes the calculator with an RSI calculator for internal use
func NewStochasticRSICalculator() *StochasticRSICalculator {
	return &StochasticRSICalculator{
		rsiCalculator: NewRSICalculator(), // Initialize RSI calculator
	}
}

// StochasticRSIResult contains the result of Stochastic RSI calculation
// This structure holds the %K and %D lines along with crossover information
type StochasticRSIResult struct {
	K         float64 // %K line (fast stochastic of RSI)
	D         float64 // %D line (smoothed %K, typically 3-period SMA of %K)
	Crossover bool    // True if %K crossed above %D (bullish crossover)
}

// Calculate calculates Stochastic RSI and returns K, D values and crossover signal
// This method applies the Stochastic oscillator formula to RSI values
// Formula: %K = ((RSI - Lowest RSI) / (Highest RSI - Lowest RSI)) * 100
// %D is typically a 3-period SMA of %K values
func (s *StochasticRSICalculator) Calculate(prices []float64, rsiPeriod, stochKPeriod, stochDPeriod int) StochasticRSIResult {
	if len(prices) < rsiPeriod+stochKPeriod+stochDPeriod {
		return StochasticRSIResult{}
	}

	// Calculate RSI values
	rsiValues := make([]float64, 0)
	for i := rsiPeriod; i < len(prices); i++ {
		rsi := s.rsiCalculator.Calculate(prices[i-rsiPeriod:i+1], rsiPeriod)
		rsiValues = append(rsiValues, rsi)
	}

	if len(rsiValues) < stochKPeriod+stochDPeriod {
		return StochasticRSIResult{}
	}

	// Calculate Stochastic K values
	stochKValues := make([]float64, 0)
	for i := stochKPeriod - 1; i < len(rsiValues); i++ {
		periodStart := i - stochKPeriod + 1
		if periodStart < 0 {
			periodStart = 0
		}

		highestRSI := rsiValues[periodStart]
		lowestRSI := rsiValues[periodStart]

		for j := periodStart; j <= i; j++ {
			if rsiValues[j] > highestRSI {
				highestRSI = rsiValues[j]
			}
			if rsiValues[j] < lowestRSI {
				lowestRSI = rsiValues[j]
			}
		}

		currentRSI := rsiValues[i]
		if highestRSI == lowestRSI {
			stochKValues = append(stochKValues, 50)
		} else {
			stochK := ((currentRSI - lowestRSI) / (highestRSI - lowestRSI)) * 100
			stochKValues = append(stochKValues, stochK)
		}
	}

	// Calculate Stochastic D (SMA of K values)
	if len(stochKValues) < stochDPeriod {
		return StochasticRSIResult{}
	}

	currentK := stochKValues[len(stochKValues)-1]

	// Calculate D as SMA of last stochDPeriod K values
	sum := 0.0
	for i := len(stochKValues) - stochDPeriod; i < len(stochKValues); i++ {
		sum += stochKValues[i]
	}
	currentD := sum / float64(stochDPeriod)

	// Check for crossover (K crossing above D from below 30)
	var crossover bool
	if len(stochKValues) >= 2 {
		prevK := stochKValues[len(stochKValues)-2]
		prevD := 0.0
		if len(stochKValues) >= stochDPeriod+1 {
			sum = 0.0
			for i := len(stochKValues) - stochDPeriod - 1; i < len(stochKValues)-1; i++ {
				sum += stochKValues[i]
			}
			prevD = sum / float64(stochDPeriod)
		}

		// Crossover: K was below D and now above D, and K was below 30
		crossover = prevK < prevD && currentK > currentD && prevK < 30
	}

	return StochasticRSIResult{
		K:         currentK,
		D:         currentD,
		Crossover: crossover,
	}
}

// IsOversoldWithCrossover checks if Stochastic RSI is oversold with crossover signal
// This method is used for Long scenario validation in the SAPAN strategy
// Returns true if %K is below 30 (oversold) and there's a bullish crossover
func (s *StochasticRSICalculator) IsOversoldWithCrossover(prices []float64, rsiPeriod, stochKPeriod, stochDPeriod int) bool {
	result := s.Calculate(prices, rsiPeriod, stochKPeriod, stochDPeriod)
	return result.K < 30 && result.Crossover // Oversold + bullish crossover
}

// IsOverboughtWithCrossover checks if Stochastic RSI is overbought with crossover signal
// This method is used for Short scenario validation in the SAPAN strategy
// Returns true if %K is above 70 (overbought) and there's a bullish crossover
func (s *StochasticRSICalculator) IsOverboughtWithCrossover(prices []float64, rsiPeriod, stochKPeriod, stochDPeriod int) bool {
	result := s.Calculate(prices, rsiPeriod, stochKPeriod, stochDPeriod)
	return result.K > 70 && result.Crossover // Overbought + bullish crossover
}
