package backoff

import (
	"math/rand/v2"
	"time"
)

// Constant will cause a constant delay after each adverse event.
// Do not specify negative duration for any of the fields.
type Constant struct {
	// ConstantDelay is the delay after each adverse event.
	ConstantDelay time.Duration

	// MaxJitter is the maximum jitter of calculated delay. Leave at 0 to disable jitter.
	MaxJitter time.Duration
}

var _ Backoff = (*Constant)(nil)

func (c Constant) Delay(_ int64) time.Duration {
	var jitter int64

	// Calculate a random jitter if MaxJitter is not zero.
	if c.MaxJitter != 0 {
		c.MaxJitter = c.MaxJitter.Abs()
		jitter = rand.Int64N(int64(c.MaxJitter<<1)) - int64(c.MaxJitter)
	}

	return c.ConstantDelay.Abs() + time.Duration(jitter)
}
