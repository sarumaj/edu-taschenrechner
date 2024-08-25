package calc

import (
	"context"
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
			if got, err := Factorial(context.TODO(), big.NewFloat(float64(tt.args.n)), tt.args.step); err != nil {
				t.Errorf("Error calculating factorial of %d: %v", tt.args.n, err)
			} else if got.Cmp(big.NewFloat(float64(tt.want))) != 0 {
				t.Errorf("Factorial(%d, %d) = %v, want %d", tt.args.n, tt.args.step, got, tt.want)
			}
		})
	}
}

func TestGreatestCommonDivisor(t *testing.T) {
	type args struct {
		x, y int
	}
	for _, tt := range []struct {
		name string
		args args
		want int
	}{
		{"test#1", args{0, 0}, 0},
		{"test#2", args{0, 1}, 1},
		{"test#3", args{1, 0}, 1},
		{"test#4", args{1, 1}, 1},
		{"test#5", args{1, 2}, 1},
		{"test#6", args{2, 1}, 1},
		{"test#7", args{2, 2}, 2},
		{"test#8", args{2, 3}, 1},
		{"test#9", args{3, 2}, 1},
		{"test#10", args{3, 3}, 3},
		{"test#11", args{3, 4}, 1},
		{"test#12", args{4, 3}, 1},
		{"test#13", args{4, 4}, 4},
		{"test#14", args{4, 5}, 1},
		{"test#15", args{5, 4}, 1},
		{"test#16", args{5, 5}, 5},
		{"test#17", args{5, 6}, 1},
		{"test#18", args{6, 5}, 1},
		{"test#19", args{6, 6}, 6},
		{"test#20", args{6, 7}, 1},
		{"test#21", args{7, 6}, 1},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := GreatestCommonDivisor(context.TODO(), big.NewFloat(float64(tt.args.x)), big.NewFloat(float64(tt.args.y))); err != nil {
				t.Errorf("Error calculating gcd(%d, %d): %v", tt.args.x, tt.args.y, err)
			} else if got.Cmp(big.NewFloat(float64(tt.want))) != 0 {
				t.Errorf("GreatestCommonDivisor(%d, %d) = %v, want %d", tt.args.x, tt.args.y, got, tt.want)
			}
		})
	}
}

func TestLeastCommonMultiple(t *testing.T) {
	type args struct {
		x, y int
	}

	for _, tt := range []struct {
		name string
		args args
		want int
	}{
		{"test#1", args{0, 0}, 0},
		{"test#2", args{0, 1}, 0},
		{"test#3", args{1, 0}, 0},
		{"test#4", args{1, 1}, 1},
		{"test#5", args{1, 2}, 2},
		{"test#6", args{2, 1}, 2},
		{"test#7", args{2, 2}, 2},
		{"test#8", args{2, 3}, 6},
		{"test#9", args{3, 2}, 6},
		{"test#10", args{3, 3}, 3},
		{"test#11", args{3, 4}, 12},
		{"test#12", args{4, 3}, 12},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := LeastCommonMultiple(context.TODO(), big.NewFloat(float64(tt.args.x)), big.NewFloat(float64(tt.args.y))); err != nil {
				t.Errorf("Error calculating lcm(%d, %d): %v", tt.args.x, tt.args.y, err)
			} else if got.Cmp(big.NewFloat(float64(tt.want))) != 0 {
				t.Errorf("LeastCommonMultiple(%d, %d) = %v, want %d", tt.args.x, tt.args.y, got, tt.want)
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
			if got, err := Pow(context.TODO(), big.NewFloat(tt.args.x), big.NewFloat(tt.args.y)); err != nil {
				t.Errorf("Error calculating %f^%f: %v", tt.args.x, tt.args.y, err)
			} else if got.Cmp(big.NewFloat(tt.want)) != 0 {
				t.Errorf("Pow(%f, %f) = %v, want %f", tt.args.x, tt.args.y, got, tt.want)
			}
		})
	}
}
