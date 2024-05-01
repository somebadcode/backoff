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

func tryN(t *testing.T, n int, err error) backoff.DoFunc {
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

func TestRetrier_RetryWithTimeout(t *testing.T) {
	type fields struct {
		Backoff          backoff.Backoff
		MaxAdverseEvents uint
		BaseDelay        time.Duration
		MaxDelay         time.Duration
		Jitter           time.Duration
	}

	type args struct {
		ctx     context.Context
		do      backoff.DoFunc
		timeout time.Duration
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "exponential_10ms_2_max100ms",
			fields: fields{
				Backoff:   backoff.Exponential(2),
				BaseDelay: 10 * time.Millisecond,
				MaxDelay:  15 * time.Millisecond,
				Jitter:    5 * time.Millisecond,
			},
			args: args{
				do:      tryN(t, 5, errDummy),
				timeout: 100 * time.Millisecond,
			},
		},
		{
			name: "exponential_10ms_2_max100ms_timeout_10ms",
			fields: fields{
				Backoff:   backoff.Exponential(2),
				BaseDelay: 5 * time.Millisecond,
				MaxDelay:  100 * time.Millisecond,
			},
			args: args{
				do:      tryN(t, 5, errDummy),
				timeout: 10 * time.Millisecond,
			},
			wantErr: true,
		},
		{
			name: "exponential_max_adverse_events",
			fields: fields{
				Backoff:          backoff.Exponential(2),
				MaxAdverseEvents: 1,
			},
			args: args{
				do:      tryN(t, 5, errDummy),
				timeout: 50 * time.Millisecond,
			},
			wantErr: true,
		},
		{
			name: "constant_10ms",
			fields: fields{
				Backoff:   backoff.Constant(),
				BaseDelay: 10 * time.Millisecond,
				Jitter:    5 * time.Millisecond,
			},
			args: args{
				do:      tryN(t, 2, errDummy),
				timeout: 30 * time.Millisecond,
			},
		},
		{
			name: "constant_deadline_exceeded",
			fields: fields{
				Backoff:   backoff.Constant(),
				BaseDelay: 1 * time.Millisecond,
			},
			args: args{
				do:      tryN(t, 100, errDummy),
				timeout: 5 * time.Millisecond,
			},
			wantErr: true,
		},
		{
			name: "constant_10ms_timeout_20ms",
			fields: fields{
				Backoff:   backoff.Constant(),
				BaseDelay: 10 * time.Millisecond,
			},
			args: args{
				do:      tryN(t, 5, errDummy),
				timeout: 20 * time.Millisecond,
			},
			wantErr: true,
		},
		{
			name: "linear_10ms",
			fields: fields{
				Backoff:   backoff.Linear(),
				BaseDelay: 10 * time.Millisecond,
				Jitter:    5 * time.Millisecond,
			},
			args: args{
				do:      tryN(t, 2, errDummy),
				timeout: 30 * time.Millisecond,
			},
		},
		{
			name: "linear_10ms_timeout_20ms",
			fields: fields{
				Backoff:   backoff.Linear(),
				BaseDelay: 10 * time.Millisecond,
				MaxDelay:  100 * time.Millisecond,
			},
			args: args{
				do:      tryN(t, 5, errDummy),
				timeout: 20 * time.Millisecond,
			},
			wantErr: true,
		},
	}

	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &backoff.Retryable{
				Backoff:          tt.fields.Backoff,
				MaxAdverseEvents: tt.fields.MaxAdverseEvents,
				BaseDelay:        tt.fields.BaseDelay,
				MaxDelay:         tt.fields.MaxDelay,
				Jitter:           tt.fields.Jitter,
			}

			var ctx context.Context

			if tt.args.ctx != nil {
				ctx = tt.args.ctx
			} else {
				ctx = context.Background()
			}

			if err := r.RetryWithTimeout(ctx, tt.args.timeout, tt.args.do); (err != nil) != tt.wantErr {
				t.Errorf("RetryWithTimeout() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
