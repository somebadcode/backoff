package backoff

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrMaxAdverseEventsExceeded = errors.New("max adverse events exceeded")
)

// RetryWithTimeout is a wrapper of Retry which will abort the retries when timeout has been reached. This function only
// provides what you could do yourself by passing a context with a deadline.
func RetryWithTimeout(ctx context.Context, backoff Backoff, maxAdverseEvents uint, try TryFunc, timeout time.Duration) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return Retry(timeoutCtx, backoff, maxAdverseEvents, try)
}

// Retry will retry the TryFunc until it returns a nil error. Control the cancellation using the context that you pass in.
// The returned error is the state of the passed in context, not the error returned by TryFunc.
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
			return fmt.Errorf("%d out of %d maximum adverse events: %w", adverseEvents+1, maxAdverseEvents, ErrMaxAdverseEventsExceeded)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff.Delay(adverseEvents)):
		}
	}

	return ctx.Err()
}
