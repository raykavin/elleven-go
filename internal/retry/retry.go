// Package retry provides exponential backoff utilities for the elleven SDK.
package retry

import (
	"math"
	"time"
)

// Wait calculates the exponential backoff wait duration for a given attempt number.
// attempt starts at 1 for the first retry.
// The result is capped at max.
func Wait(attempt int, min, max time.Duration) time.Duration {
	if attempt <= 0 {
		return min
	}
	// exponential: min * 2^(attempt-1)
	exp := math.Pow(2, float64(attempt-1))
	wait := time.Duration(float64(min) * exp)
	if wait > max {
		wait = max
	}
	return wait
}
