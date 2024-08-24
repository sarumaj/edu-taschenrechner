package calc

import (
	"math/big"
	"testing"
)

func TestFactorial(t *testing.T) {
	type args struct {
		n    int
		step int
	}

	for _, tt := range []struct {
		name string
		args args
		want int
	}{
		{"test#1", args{0, 1}, 1},
		{"test#2", args{1, 1}, 1},
		{"test#3", args{2, 1}, 2},
		{"test#4", args{3, 1}, 6},
		{"test#5", args{4, 1}, 24},
		{"test#6", args{5, 1}, 120},
		{"test#7", args{6, 2}, 48},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := Factorial(big.NewFloat(float64(tt.args.n)), tt.args.step); err != nil {
				t.Errorf("Error calculating factorial of %d: %v", tt.args.n, err)
			} else if got.Cmp(big.NewFloat(float64(tt.want))) != 0 {
				t.Errorf("Factorial(%d, %d) = %v, want %d", tt.args.n, tt.args.step, got, tt.want)
			}
		})
	}
}

func TestPow(t *testing.T) {
	type args struct {
		x, y float64
	}

	for _, tt := range []struct {
		name string
		args args
		want float64
	}{
		{"test#1", args{3, 2}, 9},
		{"test#2", args{2.5, 2}, 6.25},
		{"test#3", args{5, -2}, 1.0 / 25},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := Pow(big.NewFloat(tt.args.x), big.NewFloat(tt.args.y)); err != nil {
				t.Errorf("Error calculating %f^%f: %v", tt.args.x, tt.args.y, err)
			} else if got.Cmp(big.NewFloat(tt.want)) != 0 {
				t.Errorf("Pow(%f, %f) = %v, want %f", tt.args.x, tt.args.y, got, tt.want)
			}
		})
	}
}
