package backoff

import (
	"time"
)

// ConstantBackoff implements [Backoff] where the returned delay will always be the base delay.
type ConstantBackoff struct{}

var _ Backoff = (*ConstantBackoff)(nil)

func Constant() ConstantBackoff {
	return ConstantBackoff{}
}

// Delay returns the base delay given by the caller. The number of adverse events are ignored.
func (c ConstantBackoff) Delay(baseDelay time.Duration, _ int64) time.Duration {
	return baseDelay
}
