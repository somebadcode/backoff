package backoff

import (
	"context"
	"fmt"
	"math/rand/v2"
	"time"
)

type MaxAdverseEventsReachedError struct {
	Max    uint
	Actual int64
}

func (e MaxAdverseEventsReachedError) Error() string {
	return fmt.Sprintf("%d out of %d maximum adverse events", e.Actual, e.Max)
}

type DoFunc func(ctx context.Context) error

type Retryable struct {
	Backoff          Backoff
	MaxAdverseEvents uint
	BaseDelay        time.Duration
	MaxDelay         time.Duration
	Jitter           time.Duration
}

func (r *Retryable) Retry(ctx context.Context, do DoFunc) error {
	for adverseEvents := int64(0); ctx.Err() == nil; adverseEvents++ {
		err := do(ctx)
		if err == nil {
			return nil
		}

		// Last attempt failed, see if maximum adverse events has been reached.
		if int64(r.MaxAdverseEvents) == adverseEvents+1 {
			return MaxAdverseEventsReachedError{
				Max:    r.MaxAdverseEvents,
				Actual: adverseEvents + 1,
			}
		}

		var jitter int64

		if r.Jitter != 0 {
			jitter = rand.Int64N(int64(r.Jitter<<1)) - int64(r.Jitter)
		}

		delay := r.Backoff.Delay(r.BaseDelay, adverseEvents) + time.Duration(jitter)

		if r.BaseDelay < r.MaxDelay {
			delay = max(delay, r.MaxDelay)
		}

		select {
		case <-ctx.Done():
			break
		case <-time.After(delay):
		}
	}

	return ctx.Err()
}

// RetryWithTimeout is a wrapper of Retry which will abort the retries when timeout has been reached. This function only
// provides what you could do yourself by passing a context with a deadline.
func (r *Retryable) RetryWithTimeout(ctx context.Context, timeout time.Duration, do DoFunc) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return r.Retry(timeoutCtx, do)
}
