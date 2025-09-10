// Package strategy provides the core SAPAN trading strategy implementation
// This package contains the main strategy logic, pattern detection, and validation methods
package strategy

import "sapan/models"

// CandlestickPatternDetector handles candlestick pattern detection for the SAPAN strategy
// This struct provides methods to detect various reversal patterns including 2-candlestick and pinbar patterns
type CandlestickPatternDetector struct{}

// NewCandlestickPatternDetector creates a new candlestick pattern detector instance
// This constructor initializes the detector for identifying trading patterns
func NewCandlestickPatternDetector() *CandlestickPatternDetector {
	return &CandlestickPatternDetector{}
}

// PatternType represents the type of pattern detected by the pattern detector
// This enum helps identify which specific pattern was found during analysis
type PatternType int

const (
	NoPattern                 PatternType = iota // No valid pattern detected
	Long2CandlestickReversal                     // 2-candlestick bullish reversal pattern
	Short2CandlestickReversal                    // 2-candlestick bearish reversal pattern
	LongPinbarReversal                           // Bullish pinbar reversal pattern
	ShortPinbarReversal                          // Bearish pinbar reversal pattern
)

// DetectAllPatterns detects all possible patterns (long and short, 1 and 2 candlestick)
func (c *CandlestickPatternDetector) DetectAllPatterns(candles []models.Candle, ema20, ema50, ema100, ema200 float64) PatternType {
	if len(candles) < 3 {
		return NoPattern
	}

	// Check for 2-candlestick patterns first
	if c.DetectLong2CandlestickReversal(candles, ema20, ema50, ema100, ema200) {
		return Long2CandlestickReversal
	}

	if c.DetectShort2CandlestickReversal(candles, ema20, ema50, ema100, ema200) {
		return Short2CandlestickReversal
	}

	// Check for 1-candlestick pinbar patterns
	if c.DetectLongPinbarReversal(candles, ema20, ema50, ema100, ema200) {
		return LongPinbarReversal
	}

	if c.DetectShortPinbarReversal(candles, ema20, ema50, ema100, ema200) {
		return ShortPinbarReversal
	}

	return NoPattern
}

// DetectLong2CandlestickReversal detects long 2-candlestick reversal pattern
func (c *CandlestickPatternDetector) DetectLong2CandlestickReversal(candles []models.Candle, ema20, ema50, ema100, ema200 float64) bool {
	if len(candles) < 3 {
		return false
	}

	// Get the last 3 candles (we need at least 2 bearish + 1 bullish confirmation)
	lastCandle := candles[len(candles)-1]   // Confirmation candle
	secondCandle := candles[len(candles)-2] // Reversal candle
	firstCandle := candles[len(candles)-3]  // Previous bear candle

	// Rule A: Reversal candle body should be above EMA support
	if !c.isReversalBodyAboveSupport(secondCandle, ema20, ema50, ema100, ema200) {
		return false
	}

	// Rule B: Reversal candle tail should pierce EMA support and previous bear candle low
	if !c.isTailPiercingSupport(secondCandle, firstCandle, ema20, ema50, ema100, ema200) {
		return false
	}

	// Rule C: After reversal candle, we need rising lows and bullish confirmation
	if !c.isBullishConfirmation(lastCandle, secondCandle) {
		return false
	}

	return true
}

// DetectShort2CandlestickReversal detects short 2-candlestick reversal pattern
func (c *CandlestickPatternDetector) DetectShort2CandlestickReversal(candles []models.Candle, ema20, ema50, ema100, ema200 float64) bool {
	if len(candles) < 3 {
		return false
	}

	// Get the last 3 candles (we need at least 2 bullish + 1 bearish confirmation)
	lastCandle := candles[len(candles)-1]   // Confirmation candle
	secondCandle := candles[len(candles)-2] // Reversal candle
	firstCandle := candles[len(candles)-3]  // Previous bull candle

	// Rule A: Reversal candle body should be below EMA resistance
	if !c.isReversalBodyBelowResistance(secondCandle, ema20, ema50, ema100, ema200) {
		return false
	}

	// Rule B: Reversal candle tail should pierce EMA resistance and previous bull candle high
	if !c.isTailPiercingResistance(secondCandle, firstCandle, ema20, ema50, ema100, ema200) {
		return false
	}

	// Rule C: After reversal candle, we need falling highs and bearish confirmation
	if !c.isBearishConfirmation(lastCandle, secondCandle) {
		return false
	}

	return true
}

// DetectLongPinbarReversal detects long pinbar reversal pattern
func (c *CandlestickPatternDetector) DetectLongPinbarReversal(candles []models.Candle, ema20, ema50, ema100, ema200 float64) bool {
	if len(candles) < 3 {
		return false
	}

	// Get the last 2 candles (pinbar + confirmation)
	pinbar := candles[len(candles)-2]       // Pinbar candle
	confirmation := candles[len(candles)-1] // Confirmation candle

	// Check if it's a bullish pinbar (small body, long lower wick)
	if !c.isBullishPinbar(pinbar) {
		return false
	}

	// Rule A: Pinbar body should be above EMA support
	emaSupport := c.getLowestEMA(ema20, ema50, ema100, ema200)
	pinbarBody := (pinbar.Open + pinbar.Close) / 2
	if pinbarBody <= emaSupport {
		return false
	}

	// Rule B: Pinbar tail should pierce EMA support
	if pinbar.Low >= emaSupport {
		return false
	}

	// Rule C: Confirmation candle should be bullish and close above pinbar high
	if !c.isBullishConfirmation(confirmation, pinbar) {
		return false
	}

	return true
}

// DetectShortPinbarReversal detects short pinbar reversal pattern
func (c *CandlestickPatternDetector) DetectShortPinbarReversal(candles []models.Candle, ema20, ema50, ema100, ema200 float64) bool {
	if len(candles) < 3 {
		return false
	}

	// Get the last 2 candles (pinbar + confirmation)
	pinbar := candles[len(candles)-2]       // Pinbar candle
	confirmation := candles[len(candles)-1] // Confirmation candle

	// Check if it's a bearish pinbar (small body, long upper wick)
	if !c.isBearishPinbar(pinbar) {
		return false
	}

	// Rule A: Pinbar body should be below EMA resistance
	emaResistance := c.getHighestEMA(ema20, ema50, ema100, ema200)
	pinbarBody := (pinbar.Open + pinbar.Close) / 2
	if pinbarBody >= emaResistance {
		return false
	}

	// Rule B: Pinbar tail should pierce EMA resistance
	if pinbar.High <= emaResistance {
		return false
	}

	// Rule C: Confirmation candle should be bearish and close below pinbar low
	if !c.isBearishConfirmation(confirmation, pinbar) {
		return false
	}

	return true
}

// isReversalBodyAboveSupport checks if reversal candle body is above EMA support
func (c *CandlestickPatternDetector) isReversalBodyAboveSupport(candle models.Candle, ema20, ema50, ema100, ema200 float64) bool {
	// We use the lowest EMA as support level
	emaSupport := c.getLowestEMA(ema20, ema50, ema100, ema200)

	// Check if reversal candle body is above EMA support
	reversalBody := (candle.Open + candle.Close) / 2
	return reversalBody > emaSupport
}

// isTailPiercingSupport checks if tail pierces support levels
func (c *CandlestickPatternDetector) isTailPiercingSupport(reversalCandle, previousCandle models.Candle, ema20, ema50, ema100, ema200 float64) bool {
	emaSupport := c.getLowestEMA(ema20, ema50, ema100, ema200)
	reversalLow := reversalCandle.Low
	previousBearLow := previousCandle.Low

	// Tail should pierce both EMA support and previous bear candle low
	return reversalLow < emaSupport && reversalLow < previousBearLow
}

// isBullishConfirmation checks for bullish confirmation pattern
func (c *CandlestickPatternDetector) isBullishConfirmation(confirmationCandle, reversalCandle models.Candle) bool {
	// Confirmation candle should close above reversal candle high
	if confirmationCandle.Close <= reversalCandle.High {
		return false
	}

	// Additional check: Confirmation candle should be bullish (green)
	if confirmationCandle.Close <= confirmationCandle.Open {
		return false
	}

	// Check for rising lows (confirmation candle low should be higher than reversal candle low)
	return confirmationCandle.Low > reversalCandle.Low
}

// getLowestEMA returns the lowest EMA value
func (c *CandlestickPatternDetector) getLowestEMA(ema20, ema50, ema100, ema200 float64) float64 {
	emaSupport := ema20
	if ema50 < emaSupport {
		emaSupport = ema50
	}
	if ema100 < emaSupport {
		emaSupport = ema100
	}
	if ema200 < emaSupport {
		emaSupport = ema200
	}
	return emaSupport
}

// getHighestEMA returns the highest EMA value
func (c *CandlestickPatternDetector) getHighestEMA(ema20, ema50, ema100, ema200 float64) float64 {
	emaResistance := ema20
	if ema50 > emaResistance {
		emaResistance = ema50
	}
	if ema100 > emaResistance {
		emaResistance = ema100
	}
	if ema200 > emaResistance {
		emaResistance = ema200
	}
	return emaResistance
}

// isReversalBodyBelowResistance checks if reversal candle body is below EMA resistance
func (c *CandlestickPatternDetector) isReversalBodyBelowResistance(candle models.Candle, ema20, ema50, ema100, ema200 float64) bool {
	emaResistance := c.getHighestEMA(ema20, ema50, ema100, ema200)
	reversalBody := (candle.Open + candle.Close) / 2
	return reversalBody < emaResistance
}

// isTailPiercingResistance checks if tail pierces resistance levels
func (c *CandlestickPatternDetector) isTailPiercingResistance(reversalCandle, previousCandle models.Candle, ema20, ema50, ema100, ema200 float64) bool {
	emaResistance := c.getHighestEMA(ema20, ema50, ema100, ema200)
	reversalHigh := reversalCandle.High
	previousBullHigh := previousCandle.High

	// Tail should pierce both EMA resistance and previous bull candle high
	return reversalHigh > emaResistance && reversalHigh > previousBullHigh
}

// isBearishConfirmation checks for bearish confirmation pattern
func (c *CandlestickPatternDetector) isBearishConfirmation(confirmationCandle, reversalCandle models.Candle) bool {
	// Confirmation candle should close below reversal candle low
	if confirmationCandle.Close >= reversalCandle.Low {
		return false
	}

	// Additional check: Confirmation candle should be bearish (red)
	if confirmationCandle.Close >= confirmationCandle.Open {
		return false
	}

	// Check for falling highs (confirmation candle high should be lower than reversal candle high)
	return confirmationCandle.High < reversalCandle.High
}

// isBullishPinbar checks if candle is a bullish pinbar
func (c *CandlestickPatternDetector) isBullishPinbar(candle models.Candle) bool {
	bodySize := abs(candle.Close - candle.Open)
	totalRange := candle.High - candle.Low

	// Small body relative to total range
	if bodySize/totalRange > 0.3 {
		return false
	}

	// Long lower wick (at least 60% of total range)
	lowerWick := min(candle.Open, candle.Close) - candle.Low
	return lowerWick/totalRange >= 0.6
}

// isBearishPinbar checks if candle is a bearish pinbar
func (c *CandlestickPatternDetector) isBearishPinbar(candle models.Candle) bool {
	bodySize := abs(candle.Close - candle.Open)
	totalRange := candle.High - candle.Low

	// Small body relative to total range
	if bodySize/totalRange > 0.3 {
		return false
	}

	// Long upper wick (at least 60% of total range)
	upperWick := candle.High - max(candle.Open, candle.Close)
	return upperWick/totalRange >= 0.6
}

// Helper functions
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
