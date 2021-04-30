package billablesbuddy

import (
	"math"
)

type Status int

const (
	gracePeriodInHours        = 0.25
	StatusOnTrack      Status = iota
	StatusAhead
	StatusBehind
	StatusOver
)

func getCurrentStatus(actual float64, expected float64, expectedTotal float64) Status {
	difference := actual - expected
	isBehindExpectedHours := difference < 0
	isDifferenceZero := difference == 0.0
	isWithinGracePeriod := math.Abs(difference) <= gracePeriodInHours

	switch {
	case actual > expectedTotal:
		// If the total expected hours for the week are exceeded
		return StatusOver
	case isDifferenceZero || isWithinGracePeriod:
		// The user is on track or within the grace period
		return StatusOnTrack
	case isBehindExpectedHours:
		// User is behind
		return StatusBehind
	case !isBehindExpectedHours:
		// User is ahead
		return StatusAhead
	default:
		// Logical error
		return -1
	}
}
