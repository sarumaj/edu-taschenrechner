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
		{"test#6", "6!", []string{"6", "!"}},
		{"test#7", "( 0.3 + 2.7 )!", []string{"(", "0.3", "+", "2.7", ")", "!"}},
		{"test#8", "√7", []string{"√", "7"}},
		{"test#9", "√( 1+3 )", []string{"√", "(", "1", "+", "3", ")"}},
		{"test#10", "pi * 3", []string{"pi", "*", "3"}},
		{"test#11", "3 ^ 2", []string{"3", "^", "2"}},
		{"test#12", "2.5 ^ 2", []string{"2.5", "^", "2"}},
		{"test#13", "e ^ 2", []string{"e", "^", "2"}},
		{"test#14", "sin( pi / 2 )", []string{"sin", "(", "pi", "/", "2", ")"}},
		{"test#15", "sin( pi / 2 )!", []string{"sin", "(", "pi", "/", "2", ")", "!"}},
		{"test#16", "sin( 30° )! + 1", []string{"sin", "(", "30", "°", ")", "!", "+", "1"}},
		{"test#17", "6 ^ - 2", []string{"6", "^", "-", "2"}},
		{"test#18", "6!°", []string{"6", "!", "°"}},
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
