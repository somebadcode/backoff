package backoff

import (
	"time"
)

// ExponentialBackoff satisfies [Backoff] and will cause exponentially increasing delay.
//
// Formula: x(n) = b × fⁿ⁻¹ ± j, where n is the number of adverse events.
type ExponentialBackoff struct {
	// Factor is used to calculate the exponentially increasing the backoff delay.
	Factor int64
}

func Exponential(factor int64) ExponentialBackoff {
	return ExponentialBackoff{
		Factor: factor,
	}
}

var _ Backoff = (*ExponentialBackoff)(nil)

func (e ExponentialBackoff) Delay(baseDelay time.Duration, adverseEvents int64) time.Duration {
	multiplier := pow(max(2, e.Factor), max(0, adverseEvents-1))

	return max(baseDelay*time.Duration(multiplier), baseDelay)
}
