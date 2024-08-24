package parser

import (
	"fmt"
	"math/big"
	"regexp"
	"strings"
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
	Replacements map[string]string
	Variables    map[string]*big.Float
	Functions    map[string]func(...*big.Float) (*big.Float, error)
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
	tokens, err := Tokenize(opts.Replace(expr))
	if err != nil {
		return nil, err
	}

	root, err := tokens.Tree()
	if err != nil {
		return nil, err
	}

	return root.Evaluate(opts)
}

// Replace replaces variables and functions in the expression
func (opts *parser) Replace(expr string) string {
	for k, v := range opts.Replacements {
		expr = regexp.
			MustCompile(fmt.Sprintf(`(\b|[^a-zA-Z0-9_])%s(\b|[^a-zA-Z0-9_])`, k)).
			ReplaceAllStringFunc(
				expr,
				func(match string) string {
					if index := strings.Index(match, k); index >= 0 {
						return match[:index] + v + match[index+len(k):]
					}
					return match
				},
			)
	}

	return expr
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
		Variables:    make(map[string]*big.Float),
		Functions:    make(map[string]func(...*big.Float) (*big.Float, error)),
		Replacements: make(map[string]string),
	}

	return p.ApplyOptions(opts...)
}

// WithFunc returns an option to set a function
func WithFunc[
	F interface {
		~func(...*big.Float) (*big.Float, error) |
			~func(*big.Float) (*big.Float, error) |
			~func(*big.Float, *big.Float) (*big.Float, error) |
			~func(float64) float64 |
			~func(float64) (float64, error) |
			~func(float64, float64) float64 |
			~func(float64, float64) (float64, error)
	},
](name string, fn F, checks ...F) func(*parser) {
	return func(p *parser) {
		switch fn := any(fn).(type) {
		case func(...*big.Float) (*big.Float, error):
			p.Functions[name] = func(f ...*big.Float) (*big.Float, error) {
				return fn(f...)
			}

		case func(*big.Float) (*big.Float, error):
			p.Functions[name] = func(args ...*big.Float) (*big.Float, error) {
				if len(args) != 1 {
					return nil, fmt.Errorf("%s function requires exactly 1 argument", name)
				}
				return fn(args[0])
			}

		case func(*big.Float, *big.Float) (*big.Float, error):
			p.Functions[name] = func(args ...*big.Float) (*big.Float, error) {
				if len(args) != 2 {
					return nil, fmt.Errorf("%s function requires exactly 2 arguments", name)
				}
				return fn(args[0], args[1])
			}

		case func(float64) float64:
			p.Functions[name] = func(args ...*big.Float) (*big.Float, error) {
				if len(args) != 1 {
					return nil, fmt.Errorf("%s function requires exactly 1 argument", name)
				}
				f, _ := args[0].Float64()
				return big.NewFloat(fn(f)), nil
			}

		case func(float64) (float64, error):
			p.Functions[name] = func(args ...*big.Float) (*big.Float, error) {
				if len(args) != 1 {
					return nil, fmt.Errorf("%s function requires exactly 1 argument", name)
				}
				f, _ := args[0].Float64()
				r, err := fn(f)
				if err != nil {
					return nil, err
				}
				return big.NewFloat(r), nil
			}

		case func(float64, float64) float64:
			p.Functions[name] = func(args ...*big.Float) (*big.Float, error) {
				if len(args) != 2 {
					return nil, fmt.Errorf("%s function requires exactly 2 arguments", name)
				}
				f1, _ := args[0].Float64()
				f2, _ := args[1].Float64()
				return big.NewFloat(fn(f1, f2)), nil
			}

		case func(float64, float64) (float64, error):
			p.Functions[name] = func(args ...*big.Float) (*big.Float, error) {
				if len(args) != 2 {
					return nil, fmt.Errorf("%s function requires exactly 2 arguments", name)
				}
				f1, _ := args[0].Float64()
				f2, _ := args[1].Float64()
				r, err := fn(f1, f2)
				if err != nil {
					return nil, err
				}
				return big.NewFloat(r), nil
			}

		}
	}
}

func WithReplacement(name, value string) func(*parser) {
	return func(p *parser) {
		p.Replacements[name] = value
	}
}

func WithReplacements(replacements ...string) func(*parser) {
	return func(p *parser) {
		for i := 0; i < len(replacements); i += 2 {
			p.Replacements[replacements[i]] = replacements[i+1]
		}
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
