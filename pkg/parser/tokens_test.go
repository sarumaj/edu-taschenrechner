package parser

import "testing"

func TestExampleFor_Tokenize(t *testing.T) {
	for _, tt := range []struct {
		name string
		args string
		want []string
	}{
		{"test#1", "( x + y ) * 2.5", []string{"(", "x", "+", "y", ")", "*", "2.5"}},
		{"test#2", "x * y + ( 10.0 - y ) / 2.0", []string{"x", "*", "y", "+", "(", "10.0", "-", "y", ")", "/", "2.0"}},
		{"test#3", "max_2( x , y ) * 3.5", []string{"max_2", "(", "x", ",", "y", ")", "*", "3.5"}},
		{"test#4", "x2 * -1", []string{"x2", "*", "-", "1"}},
		{"test#5", "x2*-1", []string{"x2", "*", "-", "1"}},
	} {
		t.Run(tt.name, func(t *testing.T) {

			if tokens, err := Tokenize(tt.args); err != nil {
				t.Errorf("Error tokenizing expression %q: %v", tt.args, err)
			} else if !tokens.Compare(tt.want...) {
				t.Errorf("Tokens for %q: %v, want %v", tt.args, tokens, tt.want)
			}
		})
	}
}
