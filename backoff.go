package backoff

import (
	"time"
)

type Backoff interface {
	// Delay takes two arguments which represent the base delay and the number of adverse events and returns a backoff delay.
	// Jitter will be applied to the returned delay.
	Delay(time.Duration, int64) time.Duration
}
