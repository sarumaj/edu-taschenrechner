/*
Package parser provides a generic parser for mathematical expressions.
It can be used to parse and evaluate mathematical expressions with variables and functions.

Example:

	p := parser.NewParser(
		parser.WithConst("pi", math.Pi),
		parser.WithVar("x", func() float64 { return 42 }),
		parser.WithFunc("add", func(a, b float64) float64 { return a + b }),
	)

	result, err := p.Parse(context.Background(), "add(pi, x)")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(result) // prints 45
*/
package parser

import (
	"context"
	"fmt"
	"math/big"
	"regexp"
	"strings"
)

var _ Parser = (*parser)(nil)

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

// Parser specifies the generic interface for the parser
type Parser = ParserInterface[*parser]

// ParserInterface is a generic interface for the parser
type ParserInterface[T any] interface {
	ApplyOptions(opts ...Option) T
	LookupConst(name string) (*big.Float, bool)
	LookupFunc(name string) (func(...*big.Float) (*big.Float, error), bool)
	LookupVariable(name string) (func() *big.Float, bool)
	Parse(ctx context.Context, expr string) (*big.Float, error)
}

// parser is the implementation of the ParserInterface
type parser struct {
	constants    map[string]*big.Float
	functions    map[string]func(...*big.Float) (*big.Float, error)
	replacements map[string]string
	variables    map[string]func() *big.Float
}

// replace replaces variables and functions in the expression
func (opts *parser) replace(expr string) string {
	for k, v := range opts.replacements {
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

// LookupConst returns the value of a constant
func (opts *parser) LookupConst(name string) (*big.Float, bool) {
	v, ok := opts.constants[name]
	return v, ok
}

// LookupFunc returns the function with the given name
func (opts *parser) LookupFunc(name string) (func(...*big.Float) (*big.Float, error), bool) {
	f, ok := opts.functions[name]
	return f, ok
}

// LookupVariable returns the value of a variable
func (opts *parser) LookupVariable(name string) (func() *big.Float, bool) {
	v, ok := opts.variables[name]
	return v, ok
}

// Parse parses the expression and returns the result
func (opts *parser) Parse(ctx context.Context, expr string) (*big.Float, error) {
	tokens, err := Tokenize(opts.replace(expr))
	if err != nil {
		return nil, err
	}

	root, err := tokens.Tree()
	if err != nil {
		return nil, err
	}

	return root.Evaluate(ctx, opts)
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
		constants:    make(map[string]*big.Float),
		functions:    make(map[string]func(...*big.Float) (*big.Float, error)),
		replacements: make(map[string]string),
		variables:    make(map[string]func() *big.Float),
	}

	return p.ApplyOptions(opts...)
}

// WithConst returns an option to set a constant
func WithConst[N number](name string, value N) func(*parser) {
	return func(p *parser) {
		v, ok := ConvertToBigFloat(value)
		if !ok {
			return
		}
		p.constants[name] = v
	}
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
			p.functions[name] = func(f ...*big.Float) (*big.Float, error) {
				return fn(f...)
			}

		case func(*big.Float) (*big.Float, error):
			p.functions[name] = func(args ...*big.Float) (*big.Float, error) {
				if len(args) != 1 {
					return nil, fmt.Errorf("%s function requires exactly 1 argument", name)
				}
				return fn(args[0])
			}

		case func(*big.Float, *big.Float) (*big.Float, error):
			p.functions[name] = func(args ...*big.Float) (*big.Float, error) {
				if len(args) != 2 {
					return nil, fmt.Errorf("%s function requires exactly 2 arguments", name)
				}
				return fn(args[0], args[1])
			}

		case func(float64) float64:
			p.functions[name] = func(args ...*big.Float) (*big.Float, error) {
				if len(args) != 1 {
					return nil, fmt.Errorf("%s function requires exactly 1 argument", name)
				}
				f, _ := args[0].Float64()
				return big.NewFloat(fn(f)), nil
			}

		case func(float64) (float64, error):
			p.functions[name] = func(args ...*big.Float) (*big.Float, error) {
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
			p.functions[name] = func(args ...*big.Float) (*big.Float, error) {
				if len(args) != 2 {
					return nil, fmt.Errorf("%s function requires exactly 2 arguments", name)
				}
				f1, _ := args[0].Float64()
				f2, _ := args[1].Float64()
				return big.NewFloat(fn(f1, f2)), nil
			}

		case func(float64, float64) (float64, error):
			p.functions[name] = func(args ...*big.Float) (*big.Float, error) {
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
		p.replacements[name] = value
	}
}

func WithReplacements(replacements ...string) func(*parser) {
	return func(p *parser) {
		for i := 0; i < len(replacements); i += 2 {
			p.replacements[replacements[i]] = replacements[i+1]
		}
	}
}

// WithVar returns an option to set a variable
func WithVar[N number](name string, value func() N) func(*parser) {
	return func(p *parser) {
		p.variables[name] = func() *big.Float {
			v, ok := ConvertToBigFloat(value())
			if !ok {
				return nil
			}

			return v
		}
	}
}
