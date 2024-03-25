package backoff

import (
	"testing"
)

func Test_pow(t *testing.T) {
	type args[T integer] struct {
		a T
		b T
	}
	type testCase[T integer] struct {
		name string
		args args[T]
		want T
	}
	tests := []testCase[int]{
		{
			args: args[int]{a: 2, b: 8},
			want: 256,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pow(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("pow() = %v, want %v", got, tt.want)
			}
		})
	}
}
