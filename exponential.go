package backoff

import (
	"math/rand/v2"
	"time"
)

// Exponential implements an exponential backoff algorithm with jitter. Factor is at least 2.
//
// Formula: x(n) = b × fⁿ⁻¹ ± j, where n is the number of adverse events.
type Exponential struct {
	// BaseDelay is the delay of the initial backoff.
	BaseDelay time.Duration

	// MaxDelay is the absolute maximum of a backoff delay, which includes the jitter.
	MaxDelay time.Duration

	// MaxJitter is the maximum jitter of calculated delay. Leave at 0 to disable jitter.
	MaxJitter time.Duration

	// Factor is used for the exponentially increasing the backoff delay.
	Factor int64
}

var _ Backoff = (*Exponential)(nil)

func (e Exponential) Delay(adverseEvents int64) time.Duration {
	var jitter int64

	// Calculate a random jitter if MaxJitter is not zero.
	if e.MaxJitter != 0 {
		e.MaxJitter = e.MaxJitter.Abs()
		jitter = rand.Int64N(int64(e.MaxJitter<<1)) - int64(e.MaxJitter)
	}

	multiplier := pow(max(2, e.Factor), max(0, adverseEvents-1))

	return min(
		e.BaseDelay.Abs()*time.Duration(multiplier)+time.Duration(jitter),
		e.MaxDelay,
	).Abs()
}
