package backoff

import (
	"time"
)

// LinearBackoff satisfies [Backoff] and will cause a linearly increasing delay.
type LinearBackoff struct {
	// Slope is the multiplicative factor m in "f(n) = m × n + b" where n is the number of adverse effects and b is the
	// base delay. A slope factor of 1.0 will cause the base delay to be multiplied with the number of adverse events
	// and is therefore at a 45° angle.
	Slope float64
}

var _ Backoff = (*LinearBackoff)(nil)

func Linear(m float64) LinearBackoff {
	return LinearBackoff{
		Slope: m,
	}
}

func (l LinearBackoff) Delay(baseDelay time.Duration, adverseEvents int64) time.Duration {
	d := l.Slope*float64(adverseEvents) + baseDelay.Seconds()

	return max(time.Duration(d), baseDelay)
}
