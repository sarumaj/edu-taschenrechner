package parser

import (
	"fmt"
	"math/big"
	"testing"
)

func TestExampleFor_Parser(t *testing.T) {
	x := WithVar("x", big.NewFloat(10.5))
	y := WithVar("y", big.NewFloat(5.2))
	max_2 := WithFunc("max_2", func(args ...*big.Float) (*big.Float, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("max function requires exactly 2 arguments")
		}
		if args[0].Cmp(args[1]) >= 0 {
			return args[0], nil
		}
		return args[1], nil
	})

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
