package backoff

import (
	"math/rand/v2"
	"time"
)

type Linear struct {
	// BaseDelay is the delay of the initial backoff.
	BaseDelay time.Duration

	// MaxDelay is the absolute maximum of a backoff delay, which includes the jitter.
	MaxDelay time.Duration

	// MaxJitter is the maximum jitter of calculated delay. Leave at 0 to disable jitter.
	MaxJitter time.Duration
}

var _ Backoff = (*Linear)(nil)

func (l Linear) Delay(adverseEvents int64) time.Duration {
	var jitter int64

	// Calculate a random jitter if MaxJitter is not zero.
	if l.MaxJitter != 0 {
		l.MaxJitter = l.MaxJitter.Abs()
		jitter = rand.Int64N(int64(l.MaxJitter<<1)) - int64(l.MaxJitter)
	}

	return min((l.BaseDelay*time.Duration(adverseEvents)).Abs()+time.Duration(jitter), l.MaxDelay.Abs())
}
