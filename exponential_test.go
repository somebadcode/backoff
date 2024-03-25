package backoff_test

import (
	"testing"
	"time"

	"github.com/somebadcode/backoff"
)

func TestExponential_Delay(t *testing.T) {
	type fields struct {
		BaseDelay time.Duration
		MaxDelay  time.Duration
		Jitter    time.Duration
		Factor    int64
	}

	type args struct {
		adverseEvents int64
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   time.Duration
	}{
		{
			fields: fields{
				BaseDelay: 1 * time.Second,
				Factor:    2,
				MaxDelay:  15 * time.Second,
				Jitter:    500 * time.Millisecond,
			},
			args: args{
				adverseEvents: 1,
			},
			want: 1 * time.Second,
		},
		{
			fields: fields{
				BaseDelay: 1 * time.Second,
				Factor:    2,
				MaxDelay:  15 * time.Second,
				Jitter:    500 * time.Millisecond,
			},
			args: args{
				adverseEvents: 2,
			},
			want: 2 * time.Second,
		},
		{
			fields: fields{
				BaseDelay: 1 * time.Second,
				Factor:    2,
				MaxDelay:  15 * time.Second,
				Jitter:    500 * time.Millisecond,
			},
			args: args{
				adverseEvents: 3,
			},
			want: 4 * time.Second,
		},
		{
			fields: fields{
				BaseDelay: 1 * time.Second,
				Factor:    2,
				MaxDelay:  15 * time.Second,
				Jitter:    500 * time.Millisecond,
			},
			args: args{
				adverseEvents: 4,
			},
			want: 8 * time.Second,
		},
		{
			fields: fields{
				BaseDelay: 1 * time.Second,
				Factor:    2,
				MaxDelay:  15 * time.Second,
				Jitter:    500 * time.Millisecond,
			},
			args: args{
				adverseEvents: 5,
			},
			want: 15 * time.Second,
		},
	}

	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := backoff.Exponential{
				BaseDelay: tt.fields.BaseDelay,
				MaxDelay:  tt.fields.MaxDelay,
				MaxJitter: tt.fields.Jitter,
				Factor:    tt.fields.Factor,
			}
			if got := e.Delay(tt.args.adverseEvents); !((got >= tt.want-e.MaxJitter) && (got < tt.want+e.MaxJitter)) {
				t.Errorf("Delay() = %v, want %v-%v", got, tt.want-e.MaxJitter, tt.want+e.MaxJitter)
			}
		})
	}
}
