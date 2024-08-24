package parser

import (
	"math"
	"math/big"
	"testing"
)

func TestExampleFor_Parser(t *testing.T) {
	x := WithVar("x", big.NewFloat(10.5))
	y := WithVar("y", big.NewFloat(5.2))
	pi := WithVar("PI", big.NewFloat(math.Pi))
	e := WithVar("e", big.NewFloat(math.E))
	max_2 := WithFunc("max_2", func(x, y *big.Float) (*big.Float, error) {
		if x.Cmp(y) >= 0 {
			return x, nil
		}
		return y, nil
	})
	sin := WithFunc("sin", func(x *big.Float) (*big.Float, error) {
		f, _ := x.Float64()
		return big.NewFloat(math.Sin(f)), nil
	})
	save := WithFunc("save", func(x *big.Float) (*big.Float, error) { return x, nil })

	type args struct {
		expr string
		opts []Option
	}

	for _, tt := range []struct {
		name string
		args args
		want *big.Float
	}{
		{"test#1", args{"(x+y)*2.5", []Option{x, y}}, big.NewFloat(39.25)},
		{"test#2", args{"x*y+(10.0-y)/2.0", []Option{x, y}}, big.NewFloat(57)},
		{"test#3", args{"max_2(x ,y)*3.5", []Option{x, y, max_2}}, big.NewFloat(36.75)},
		{"test#4", args{"x*-1", []Option{x}}, big.NewFloat(-10.5)},
		{"test#5", args{"6!", []Option{x}}, big.NewFloat(720)},
		{"test#6", args{"(3+1)!", []Option{x}}, big.NewFloat(24)},
		{"test#7", args{"√9", []Option{}}, big.NewFloat(3)},
		{"test#8", args{"√(1+3)", []Option{}}, big.NewFloat(2)},
		{"test#9", args{"PI*3", []Option{pi}}, big.NewFloat(0).Mul(big.NewFloat(math.Pi), big.NewFloat(3))},
		{"test#10", args{"3^2", []Option{e}}, big.NewFloat(9)},
		{"test#11", args{"2.5^2", []Option{}}, big.NewFloat(6.25)},
		{"test#12", args{"e^2", []Option{e}}, big.NewFloat(0).Mul(big.NewFloat(math.E), big.NewFloat(math.E))},
		{"test#13", args{"sin(PI/2)", []Option{pi, sin}}, big.NewFloat(math.Sin(math.Pi / 2))},
		{"test#14", args{"30°", []Option{}}, big.NewFloat(0).Quo(big.NewFloat(math.Pi), big.NewFloat(6))},
		{"test#15", args{"sin(90°)", []Option{sin}}, big.NewFloat(1)},
		{"test#16", args{"save(10)", []Option{save}}, big.NewFloat(10)},
		{"test#17", args{"save(10)+save(20)", []Option{save}}, big.NewFloat(30)},
		{"test#18", args{"6^-2", []Option{}}, big.NewFloat(0).Quo(big.NewFloat(1), big.NewFloat(36))},
		{"test#19", args{"6!°", []Option{}}, big.NewFloat(0).Mul(big.NewFloat(720), big.NewFloat(0).Quo(big.NewFloat(math.Pi), big.NewFloat(180)))},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewParser(tt.args.opts...).Parse(tt.args.expr)
			if err != nil {
				t.Errorf("Error parsing expression %q: %v", tt.args.expr, err)
			} else if got.Cmp(tt.want) != 0 {
				t.Errorf("Result of %q: %s, want %s", tt.args.expr, got.Text('f', -1), tt.want.Text('f', -1))
			}
		})
	}
}
