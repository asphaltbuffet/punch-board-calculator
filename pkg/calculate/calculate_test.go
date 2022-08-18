// Package calculate contains calculators for envelope punch positions.
package calculate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateEnvelope(t *testing.T) {
	type args struct {
		length    float64
		width     float64
		isLoose   bool
		boardMini bool
	}
	tests := []struct {
		name              string
		args              args
		wantPaperSize     float64
		wantPunchLocation float64
		assertion         assert.ErrorAssertionFunc
	}{
		{
			name: "10x8 - loose",
			args: args{
				length:    10,
				width:     8,
				isLoose:   true,
				boardMini: false,
			},
			wantPaperSize:     15.73,
			wantPunchLocation: 7.16,
			assertion:         assert.NoError,
		},
		{
			name: "8x10 - loose",
			args: args{
				length:    8,
				width:     10,
				isLoose:   true,
				boardMini: false,
			},
			wantPaperSize:     15.73,
			wantPunchLocation: 7.16,
			assertion:         assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPaperSize, gotPunchLocation, err := CalculateEnvelope(tt.args.length, tt.args.width, tt.args.isLoose, tt.args.boardMini)

			tt.assertion(t, err)

			if err == nil {
				assert.InDelta(t, tt.wantPaperSize, gotPaperSize, 0.01)
				assert.InDelta(t, tt.wantPunchLocation, gotPunchLocation, 0.01)
			}
		})
	}
}

func TestParseDecimal(t *testing.T) {
	tests := []struct {
		name      string
		arg       string
		wantR     Rational
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "empty string",
			arg:       "",
			wantR:     Rational{},
			assertion: assert.NoError,
		},
		{
			name: "integer",
			arg:  "1",
			wantR: Rational{
				i: 1,
				n: 0,
				d: 0,
			},
			assertion: assert.NoError,
		},
		{
			name: "zero",
			arg:  "0",
			wantR: Rational{
				i: 0,
				n: 0,
				d: 0,
			},
			assertion: assert.NoError,
		},
		{
			name: "neg zero",
			arg:  "-0",
			wantR: Rational{
				i: 0,
				n: 0,
				d: 0,
			},
			assertion: assert.NoError,
		},
		{
			name: "int - no fraction",
			arg:  "1.0",
			wantR: Rational{
				i: 1,
				n: 0,
				d: 0,
			},
			assertion: assert.NoError,
		},
		{
			name: "int - fraction",
			arg:  "1.1",
			wantR: Rational{
				i: 1,
				n: 1,
				d: 10,
			},
			assertion: assert.NoError,
		},
		{
			name: "int - simplify fraction",
			arg:  "1.2",
			wantR: Rational{
				i: 1,
				n: 1,
				d: 5,
			},
			assertion: assert.NoError,
		},
		{
			name: "neg int - simplify fraction",
			arg:  "-1.2",
			wantR: Rational{
				i: -1,
				n: 1,
				d: 5,
			},
			assertion: assert.NoError,
		},
		{
			name: "neg integer",
			arg:  "-3",
			wantR: Rational{
				i: -3,
				n: 0,
				d: 0,
			},
			assertion: assert.NoError,
		},
		{
			name: "neg - no leading int",
			arg:  "-.2",
			wantR: Rational{
				i: 0,
				n: -1,
				d: 5,
			},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, err := ParseDecimal(tt.arg)

			tt.assertion(t, err)

			if err == nil {
				assert.Equal(t, tt.wantR, gotR)
			}
		})
	}
}

func Test_gcd(t *testing.T) {
	type args struct {
		a int64
		b int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "1,1",
			args: args{
				a: 1,
				b: 1,
			},
			want: 1,
		},
		{
			name: "6,27",
			args: args{
				a: 6,
				b: 27,
			},
			want: 3,
		},
		{
			name: "-6,27",
			args: args{
				a: -6,
				b: 27,
			},
			want: 3,
		},
		{
			name: "6,-27",
			args: args{
				a: 6,
				b: -27,
			},
			want: 3,
		},
		{
			name: "0,0",
			args: args{
				a: 0,
				b: 0,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := gcd(tt.args.a, tt.args.b); got != tt.want {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestAbsInt64(t *testing.T) {
	tests := []struct {
		name string
		arg  int64
		want int64
	}{
		{
			name: "zero",
			arg:  0,
			want: 0,
		},
		{
			name: "positive",
			arg:  3,
			want: 3,
		},
		{
			name: "negative",
			arg:  -3,
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AbsInt64(tt.arg); got != tt.want {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
