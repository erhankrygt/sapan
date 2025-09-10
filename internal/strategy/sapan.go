// Package strategy provides the core SAPAN trading strategy implementation
// This package contains the main strategy logic, pattern detection, and validation methods
package strategy

import (
	"sapan/internal/indicators"
	"sapan/models"
)

// SAPANStrategy implements the SAPAN trading strategy with both Long and Short scenarios
// This struct orchestrates all technical indicators and pattern detection to validate trading setups
type SAPANStrategy struct {
	emaCalculator           *indicators.EMACalculator           // EMA calculator for trend analysis
	stochasticRSICalculator *indicators.StochasticRSICalculator // Stochastic RSI calculator for momentum analysis
	macdCalculator          *indicators.MACDCalculator          // MACD calculator for trend confirmation
	patternDetector         *CandlestickPatternDetector         // Pattern detector for candlestick analysis
}

// NewSAPANStrategy creates a new SAPAN strategy instance with all required calculators
// This constructor initializes all technical indicators and pattern detectors
func NewSAPANStrategy() *SAPANStrategy {
	return &SAPANStrategy{
		emaCalculator:           indicators.NewEMACalculator(),           // Initialize EMA calculator
		stochasticRSICalculator: indicators.NewStochasticRSICalculator(), // Initialize Stochastic RSI calculator
		macdCalculator:          indicators.NewMACDCalculator(),          // Initialize MACD calculator
		patternDetector:         NewCandlestickPatternDetector(),         // Initialize pattern detector
	}
}

// ValidationResult contains the result of strategy validation for a single stock
// This structure holds all validation results and provides detailed feedback about the analysis
type ValidationResult struct {
	IsValid           bool        // Overall validation result (true if all conditions are met)
	EMATrendValid     bool        // EMA trend validation result
	StochasticValid   bool        // Stochastic RSI validation result
	MACDValid         bool        // MACD validation result
	PatternValid      bool        // Candlestick pattern validation result
	PatternType       PatternType // Type of pattern detected (if any)
	Symbol            string      // Stock symbol being analyzed
	ValidationMessage string      // Detailed message explaining the validation result
}

// ScenarioType represents the type of trading scenario being validated
// This enum helps distinguish between Long and Short trading setups
type ScenarioType int

const (
	LongScenario  ScenarioType = iota // Long (bullish) trading scenario
	ShortScenario                     // Short (bearish) trading scenario
)

// ValidateLongSetup validates if the given stock data meets SAPAN long setup criteria
// This method checks all conditions required for a bullish (long) trading setup
// Returns ValidationResult with detailed information about the validation
// Note: Long scenario has priority over Short scenario
func (s *SAPANStrategy) ValidateLongSetup(symbol string, candles []models.Candle) ValidationResult {
	return s.validateSetup(symbol, candles, LongScenario)
}

// ValidateShortSetup validates if the given stock data meets SAPAN short setup criteria
// This method checks all conditions required for a bearish (short) trading setup
// Returns ValidationResult with detailed information about the validation
// Note: Short scenario is only considered if Long scenario is not valid
func (s *SAPANStrategy) ValidateShortSetup(symbol string, candles []models.Candle) ValidationResult {
	return s.validateSetup(symbol, candles, ShortScenario)
}

// validateSetup validates setup for both long and short scenarios
// This is the core validation method that orchestrates all technical analysis checks
// It validates EMA trends, Stochastic RSI, MACD, and candlestick patterns based on the scenario
func (s *SAPANStrategy) validateSetup(symbol string, candles []models.Candle, scenario ScenarioType) ValidationResult {
	result := ValidationResult{
		Symbol: symbol,
	}

	// Extract closing prices
	closes := s.extractClosingPrices(candles)
	if len(closes) < 200 {
		result.ValidationMessage = "Insufficient data for analysis"
		return result
	}

	// Validate EMA trend based on scenario
	if scenario == LongScenario {
		result.EMATrendValid = s.validateEMATrend(closes)
		if !result.EMATrendValid {
			result.ValidationMessage = "EMA trend not in uptrend order (20 > 50 > 100 > 200)"
			return result
		}
	} else {
		result.EMATrendValid = s.validateEMADowntrend(closes)
		if !result.EMATrendValid {
			result.ValidationMessage = "EMA trend not in downtrend order (20 < 50 < 100 < 200)"
			return result
		}
	}

	// Validate Stochastic RSI based on scenario
	if scenario == LongScenario {
		result.StochasticValid = s.validateStochasticRSILong(closes)
		if !result.StochasticValid {
			result.ValidationMessage = "Stochastic RSI not in oversold region with crossover"
			return result
		}
	} else {
		result.StochasticValid = s.validateStochasticRSIShort(closes)
		if !result.StochasticValid {
			result.ValidationMessage = "Stochastic RSI not in overbought region with crossover"
			return result
		}
	}

	// Validate MACD based on scenario
	if scenario == LongScenario {
		result.MACDValid = s.validateMACDLong(closes)
		if !result.MACDValid {
			result.ValidationMessage = "MACD not in bull market or bear market exceeds 5 candlesticks"
			return result
		}
	} else {
		result.MACDValid = s.validateMACDShort(closes)
		if !result.MACDValid {
			result.ValidationMessage = "MACD not in bear market or bull market exceeds 5 candlesticks"
			return result
		}
	}

	// Validate candlestick pattern
	result.PatternType = s.patternDetector.DetectAllPatterns(candles,
		s.emaCalculator.Calculate(closes, 20),
		s.emaCalculator.Calculate(closes, 50),
		s.emaCalculator.Calculate(closes, 100),
		s.emaCalculator.Calculate(closes, 200))

	if scenario == LongScenario {
		result.PatternValid = (result.PatternType == Long2CandlestickReversal || result.PatternType == LongPinbarReversal)
		if !result.PatternValid {
			result.ValidationMessage = "Long reversal pattern not detected"
			return result
		}
	} else {
		result.PatternValid = (result.PatternType == Short2CandlestickReversal || result.PatternType == ShortPinbarReversal)
		if !result.PatternValid {
			result.ValidationMessage = "Short reversal pattern not detected"
			return result
		}
	}

	result.IsValid = true
	if scenario == LongScenario {
		result.ValidationMessage = "All SAPAN long strategy conditions met"
	} else {
		result.ValidationMessage = "All SAPAN short strategy conditions met"
	}
	return result
}

// validateEMATrend validates EMA trend according to SAPAN rules for Long scenario
// Checks if EMAs are in uptrend order: 20 > 50 > 100 > 200
func (s *SAPANStrategy) validateEMATrend(closes []float64) bool {
	return s.emaCalculator.ValidateTrend(closes)
}

// validateEMADowntrend validates EMA downtrend according to SAPAN rules for Short scenario
// Checks if EMAs are in downtrend order: 20 < 50 < 100 < 200
func (s *SAPANStrategy) validateEMADowntrend(closes []float64) bool {
	return s.emaCalculator.ValidateDowntrend(closes)
}

// validateStochasticRSILong validates Stochastic RSI for long scenario
// Checks if Stochastic RSI is oversold (< 30) with bullish crossover
func (s *SAPANStrategy) validateStochasticRSILong(closes []float64) bool {
	return s.stochasticRSICalculator.IsOversoldWithCrossover(closes, 5, 3, 3)
}

// validateStochasticRSIShort validates Stochastic RSI for short scenario
// Checks if Stochastic RSI is overbought (> 70) with bullish crossover
func (s *SAPANStrategy) validateStochasticRSIShort(closes []float64) bool {
	return s.stochasticRSICalculator.IsOverboughtWithCrossover(closes, 5, 3, 3)
}

// validateMACDLong validates MACD for long scenario
// Checks if in bull market OR bear market has lasted ≤ 5 candlesticks
func (s *SAPANStrategy) validateMACDLong(closes []float64) bool {
	return s.macdCalculator.IsBearMarketAcceptable(closes, 50, 100, 9)
}

// validateMACDShort validates MACD for short scenario
// Checks if in bear market OR bull market has lasted ≤ 5 candlesticks
func (s *SAPANStrategy) validateMACDShort(closes []float64) bool {
	return s.macdCalculator.IsBullMarketAcceptable(closes, 50, 100, 9)
}

// extractClosingPrices extracts closing prices from candles for technical analysis
// This helper method converts candle data to a slice of closing prices for indicator calculations
func (s *SAPANStrategy) extractClosingPrices(candles []models.Candle) []float64 {
	closes := make([]float64, len(candles))
	for i, candle := range candles {
		closes[i] = candle.Close // Extract closing price from each candle
	}
	return closes
}
