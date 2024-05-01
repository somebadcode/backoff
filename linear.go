package backoff

import (
	"time"
)

// LinearBackoff satisfies [Backoff] and will cause a linearly increasing delay.
type LinearBackoff struct{}

var _ Backoff = (*LinearBackoff)(nil)

func Linear() LinearBackoff {
	return LinearBackoff{}
}

func (l LinearBackoff) Delay(baseDelay time.Duration, adverseEvents int64) time.Duration {
	return max(baseDelay*time.Duration(adverseEvents), baseDelay)
}
