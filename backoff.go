package backoff

import (
	"context"
	"time"
)

type Backoff interface {
	// Delay takes one argument which represents the number of adverse events and returns a backoff delay.
	Delay(int64) time.Duration
}

type TryFunc func(ctx context.Context) error
