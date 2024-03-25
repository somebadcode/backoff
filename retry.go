package backoff

import (
	"context"
	"fmt"
	"time"
)

type MaxAdverseEventsReachedError struct {
	Max    uint
	Actual int64
}

func (e MaxAdverseEventsReachedError) Error() string {
	return fmt.Sprintf("%d out of %d maximum adverse events", e.Actual, e.Max)
}

// RetryWithTimeout is a wrapper of Retry which will abort the retries when timeout has been reached. This function only
// provides what you could do yourself by passing a context with a deadline.
func RetryWithTimeout(ctx context.Context, backoff Backoff, maxAdverseEvents uint, try TryFunc, timeout time.Duration) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return Retry(timeoutCtx, backoff, maxAdverseEvents, try)
}

// Retry will retry the TryFunc until it returns a nil error. Control the cancellation using the context that you pass in.
// The returned error is the state of the passed in context, not the error returned by TryFunc. Setting maximum adverse
// events to 0 will cause it to retry until the context is cancelled.
func Retry(ctx context.Context, backoff Backoff, maxAdverseEvents uint, try TryFunc) error {
	var adverseEvents int64
	var err error

	for ; ctx.Err() == nil; adverseEvents++ {
		err = try(ctx)
		if err == nil {
			return nil
		}

		// Last attempt failed, see if maximum adverse events has been reached.
		if int64(maxAdverseEvents) == adverseEvents+1 {
			return MaxAdverseEventsReachedError{
				Max:    maxAdverseEvents,
				Actual: adverseEvents + 1,
			}
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff.Delay(adverseEvents)):
		}
	}

	return ctx.Err()
}
