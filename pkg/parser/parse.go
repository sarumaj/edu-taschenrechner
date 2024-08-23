package parser

import (
	"fmt"
	"math/big"
)

// number is a type constraint for numbers
type number interface {
	~float64 | ~float32 |
		~int64 | ~int32 | ~int16 | ~int8 | ~int |
		~uint64 | ~uint32 | ~uint16 | ~uint8 | ~uint |
		*big.Float | *big.Int |
		~string | ~[]byte | ~[]rune
}

// Option can be used to configure the parser
type Option func(*parser)

// options holds the configuration for the parser
type parser struct {
	Variables map[string]*big.Float
	Functions map[string]func(...*big.Float) (*big.Float, error)
}

// Apply applies the options to the parser
func (o *parser) ApplyOptions(opts ...Option) *parser {
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(o)
	}

	return o
}

// Parse parses the expression and returns the result
func (opts *parser) Parse(expr string) (*big.Float, error) {
	tokens, err := Tokenize(expr)
	if err != nil {
		return nil, err
	}

	root, err := tokens.Tree()
	if err != nil {
		return nil, err
	}

	return root.Evaluate(opts)
}

// ConvertToBigFloat converts a number to a big.Float
func ConvertToBigFloat[N number](n N) (*big.Float, bool) {
	switch n := any(n).(type) {
	case *big.Float:
		return n, n != nil

	case *big.Int:
		if n == nil {
			return nil, false
		}
		return big.NewFloat(0).SetInt(n), true

	case float64, float32:
		return big.NewFloat(0).SetString(fmt.Sprintf("%g", n))

	case int64, int32, int16, int8, int, uint64, uint32, uint16, uint8, uint:
		return big.NewFloat(0).SetString(fmt.Sprintf("%d.0", n))

	case string:
		return big.NewFloat(0).SetString(n)

	case []byte, []rune:
		return big.NewFloat(0).SetString(fmt.Sprintf("%s", n))

	default:
		return nil, false

	}
}

// NewParser returns a new parser.
// It can be configured with options.
// Options can be used to set variables and functions.
// All supplied options are applied to the parser upon creation.
func NewParser(opts ...Option) *parser {
	p := &parser{
		Variables: make(map[string]*big.Float),
		Functions: make(map[string]func(...*big.Float) (*big.Float, error)),
	}

	return p.ApplyOptions(opts...)
}

// WithFunc returns an option to set a function
func WithFunc(name string, fn func(...*big.Float) (*big.Float, error)) func(*parser) {
	return func(p *parser) {
		if fn == nil {
			return
		}

		p.Functions[name] = fn
	}
}

// WithVar returns an option to set a variable
func WithVar[N number](name string, value N) func(*parser) {
	return func(p *parser) {
		if n, ok := ConvertToBigFloat(value); ok {
			p.Variables[name] = n
		}
	}
}
