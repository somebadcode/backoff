package backoff_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/somebadcode/backoff"
)

var (
	errDummy = errors.New("dummy error")
)

func tryN(t *testing.T, n int, err error) backoff.TryFunc {
	t.Helper()

	return func(ctx context.Context) error {
		t.Helper()

		n -= 1
		if n > 0 {
			return err
		}

		return ctx.Err()
	}
}

func TestRetryWithTimeout(t *testing.T) {
	type args struct {
		backoff          backoff.Backoff
		try              backoff.TryFunc
		maxAdverseEvents uint
		timeout          time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "exponential_10ms_2_max100ms",
			args: args{
				backoff: backoff.Exponential{
					BaseDelay: 10 * time.Millisecond,
					Factor:    2,
					MaxDelay:  20 * time.Millisecond,
					MaxJitter: 5 * time.Millisecond,
				},
				try:     tryN(t, 5, errDummy),
				timeout: 100 * time.Millisecond,
			},
		},
		{
			name: "exponential_10ms_2_max100ms_timeout_10ms",
			args: args{
				backoff: backoff.Exponential{
					BaseDelay: 5 * time.Millisecond,
					Factor:    2,
					MaxDelay:  100 * time.Millisecond,
				},
				try:     tryN(t, 5, errDummy),
				timeout: 10 * time.Millisecond,
			},
			wantErr: true,
		},
		{
			name: "exponential_max_adverse_events",
			args: args{
				backoff:          backoff.Exponential{},
				try:              tryN(t, 5, errDummy),
				maxAdverseEvents: 1,
				timeout:          50 * time.Millisecond,
			},
			wantErr: true,
		},
		{
			name: "constant_10ms",
			args: args{
				backoff: backoff.Constant{
					ConstantDelay: 10 * time.Millisecond,
					MaxJitter:     5 * time.Millisecond,
				},
				try:     tryN(t, 2, errDummy),
				timeout: 30 * time.Millisecond,
			},
		},
		{
			name: "constant_10ms_timeout_20ms",
			args: args{
				backoff: backoff.Constant{
					ConstantDelay: 10 * time.Millisecond,
				},
				try:     tryN(t, 5, errDummy),
				timeout: 20 * time.Millisecond,
			},
			wantErr: true,
		},
		{
			name: "linear_10ms",
			args: args{
				backoff: backoff.Linear{
					BaseDelay: 10 * time.Millisecond,
					MaxJitter: 5 * time.Millisecond,
				},
				try:     tryN(t, 2, errDummy),
				timeout: 30 * time.Millisecond,
			},
		},
		{
			name: "linear_10ms_timeout_20ms",
			args: args{
				backoff: backoff.Linear{
					BaseDelay: 10 * time.Millisecond,
				},
				try:     tryN(t, 5, errDummy),
				timeout: 20 * time.Millisecond,
			},
			wantErr: true,
		},
	}

	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := backoff.RetryWithTimeout(context.Background(), tt.args.backoff, tt.args.maxAdverseEvents, tt.args.try, tt.args.timeout); (err != nil) != tt.wantErr {
				t.Errorf("RetryWithTimeout() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
