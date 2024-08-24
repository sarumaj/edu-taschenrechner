package cursor

import (
	"fmt"
	"math"
	"math/big"

	"github.com/sarumaj/edu-taschenrechner/pkg/memory"
	"github.com/sarumaj/edu-taschenrechner/pkg/parser"
	"github.com/sarumaj/edu-taschenrechner/pkg/runes"
)

// Cursor is an actor performing operations.
// It defines the basic operations of a calculator.
type Cursor interface {
	Divide() Cursor
	Minus() Cursor
	Plus() Cursor
	Times() Cursor

	Arccos() Cursor
	Arcsin() Cursor
	Arctan() Cursor
	Cos() Cursor
	Factorial() Cursor
	Ln() Cursor
	Log() Cursor
	Power() Cursor
	Sin() Cursor
	Square() Cursor
	SquareRoot() Cursor
	Tan() Cursor

	Euler() Cursor
	Pi() Cursor
	Zero() Cursor
	One() Cursor
	Two() Cursor
	Three() Cursor
	Four() Cursor
	Five() Cursor
	Six() Cursor
	Seven() Cursor
	Eight() Cursor
	Nine() Cursor

	Degrees() Cursor
	Radians() Cursor

	Brackets() Cursor
	Clear() Cursor
	DecimalPoint() Cursor
	Delete() Cursor
	Error(error) Cursor
	Equals() Cursor

	Check() error
	String() string
}

// cursor is a cursor of operations.
// Currently, it only supports the basic operations of a calculator.
type cursor struct {
	err    error
	char   rune
	ready  bool
	text   *runes.Sequence
	memory memory.MemoryCell
}

// abortFunction aborts the function if the last character is an opening bracket.
// It is used to prevent incomplete functions.
func (c *cursor) abortFunction() (aborted bool) {
	if !c.text.EndsWith("(") || c.text.Shift().EndsWith("(") {
		return false
	}

	for !runes.IsAnyOf(c.text.Last(), "+-×÷") && runes.IsValid(c.text.Last()) {
		c.text.Backspace()
		aborted = true
	}

	if aborted && runes.IsValid(c.text.Last()) {
		c.text.Backspace()
	}

	return
}

// abortOperation aborts a binary operation if the last character is an operator.
func (c *cursor) abortOperation() (aborted bool) {
	for runes.IsAnyOf(c.text.Last(), "+-×÷.^") {
		c.text.Backspace()
		aborted = true
	}

	return
}

// binary adds a binary operator to the input text.
// If the last character is a decimal point or an operator, it is removed.
// If the last character is a digit, a memory cell value, or a closing bracket, the operator is added.
// If the operator is a power operator, it is added only if the last character is not a minus or a power operator.
// If the operator is a minus operator, it is added only if the last character is not a minus operator.
// If the operator is a minus operator, it flips the sign of the number.
func (c *cursor) binary(op rune, opts ...rune) Cursor {
	c.prepare()
	defer c.exhaust()

	if !runes.IsAnyOf(op, "×÷+-^") {
		return c.Error(fmt.Errorf("unsupported operator: %c", op))
	}

	if c.text.EndsWith(".") { // remove decimal point
		c.text.Backspace()
	}

	switch op {
	case '-':
		if c.text.EndsWith("-") { // flip sign
			c.text.Backspace()
			c.text.Append("+")

		} else {
			if c.text.EndsWith("+") { // remove if flipped too much
				c.text.Backspace()
			}
			c.text.Append("-")
		}

		if c.text.BeginsWith("+") { // remove if leading
			c.text.Delete()
		}

		if runes.IsAnyOf(c.text.Shift().Last(), "+×÷") && c.text.EndsWith("+") { // remove if flipped too much
			c.text.Backspace()
		}

	default:
		for runes.IsAnyOf(c.text.Last(), "-×÷+^") { // abort operation
			c.text.Backspace()
		}
	}

	// set operator
	if runes.IsDigit(c.text.Last()) || c.text.Equals("ANS") || runes.IsAnyOf(c.text.Last(), ")πe!°") {
		c.text.Append(string(op) + string(opts))
	}

	c.text.Append(string(c.char))
	return c
}

// character adds a character to the input text.
// If the last character is a closing bracket or the result of a previous calculation, a multiplication operator is added.
func (c *cursor) character(v rune, opts ...rune) Cursor {
	c.prepare()
	defer c.exhaust()

	// multiply if behind closing bracket, memory cell value, constant, factorial, or degree
	if c.text.Equals("ANS") || runes.IsAnyOf(c.text.Last(), ")πe!°") || (runes.IsDigit(c.text.Last()) && !runes.IsDigit(v)) {
		c.text.Append("×")
	}

	c.text.Append(string(v) + string(opts) + string(c.char))
	return c
}

// closeBrackets closes all opened brackets.
// It is used to prevent incomplete brackets.
func (c *cursor) closeBrackets() (closed bool) {
	removed := false
	for c.text.EndsWith("(") {
		c.text.Backspace()
		removed = true
	}

	if removed {
		for !runes.IsAnyOf(c.text.Last(), "+-×÷") && runes.IsValid(c.text.Last()) {
			c.text.Backspace()
		}
		c.text.Backspace()
	}

	for i, o := 0, runes.HowManyOpen(c.text); i < o; i++ {
		c.text.Append(")")
		closed = true
	}

	return
}

// exhaust disables the cursor and returns it.
// It is used to prevent further operations on the cursor.
// The cursor is disabled after an operation is completed.
// To re-enable the cursor, a prepare operation is required.
func (c *cursor) exhaust() { c.ready = false }

// function adds a function to the input text.
// If the last character is a closing bracket or the result of a previous calculation, a multiplication operator is added.
func (c *cursor) function(name string) Cursor {
	c.prepare()
	defer c.exhaust()

	// multiply if behind closing bracket or memory cell value
	if c.text.Equals("ANS") || runes.IsAnyOf(c.text.Last(), ")πe!°") || runes.IsDigit(c.text.Last()) {
		c.text.Append("×")
	}

	c.text.Append(name + "(" + string(c.char))
	return c
}

// prepare prepares the input text for a new calculation.
// If the input text ends with a cursor, it is removed.
// If the input text equals NaN, the screen is cleared.
// Otherwise, the result from the memory cell is reused.
func (c *cursor) prepare() {
	if c.ready {
		return
	}

	if c.text.EndsWith(string(c.char)) { // remove cursor
		c.text.Backspace()
	} else if c.text.Equals(fmt.Sprint(math.NaN())) { // clear screen
		c.text.Clear()
	} else { // reuse result from memory cell
		c.text.Clear()
		c.text.Append("ANS")
	}

	c.ready = true
}

// unit adds a unit to the input text.
// If the last character is a digit, a memory cell value, or a closing bracket, the unit is added.
func (c *cursor) unit(u string) Cursor {
	c.prepare()
	defer c.exhaust()

	c.abortOperation()
	backup := c.text.Copy()
	if c.abortFunction() {
		*c.text = backup
		c.text.Append(string(c.char))
		return c
	}

	if runes.IsValid(c.text.First()) &&
		!c.text.EndsWith(string(u)) &&
		(runes.IsDigit(c.text.Last()) || c.text.Equals("ANS") || runes.IsAnyOf(c.text.Last(), ")!°")) {

		c.text.Append(u)
	}

	c.text.Append(string(c.char))
	return c
}

// Arccos adds an arccos function to the input text.
// It uses the function method to add the function.
func (c *cursor) Arccos() Cursor { return c.function("arccos") }

// Arcsin adds an arcsin function to the input text.
// It uses the function method to add the function.
func (c *cursor) Arcsin() Cursor { return c.function("arcsin") }

// Arctan adds an arctan function to the input text.
// It uses the function method to add the function.
func (c *cursor) Arctan() Cursor { return c.function("arctan") }

// Brackets adds brackets to the input text.
// If the last character is an opening bracket, a closing bracket is added.
// If the last character is a digit or a closing bracket, a multiplication operator and an opening bracket are added.
// If the last character is an operator, a closing bracket is added.
func (c *cursor) Brackets() Cursor {
	c.prepare()
	defer c.exhaust()

	switch {
	case // just open
		runes.IsAnyOf(c.text.Last(), "(+-×÷√^"),
		!runes.IsValid(c.text.Last()):

		c.text.Append("(")

	case // just close
		runes.IsValid(c.text.First()) &&
			runes.HowManyOpen(c.text) > 0 &&
			(runes.IsDigit(c.text.Last()) || runes.IsAnyOf(c.text.Last(), ")πe!°")):

		c.text.Append(")")

	case // multiply and open
		runes.IsValid(c.text.First()) &&
			(runes.IsAnyOf(c.text.Last(), ")!°") ||
				runes.IsDigit(c.text.Last()) ||
				c.text.Equals("ANS")):

		c.text.Append("×(")

	}

	c.text.Append(string(c.char))
	return c
}

// Check returns the last error.
func (c *cursor) Check() error { return c.err }

// Clear clears the input text.
func (c *cursor) Clear() Cursor {
	c.text.Clear()
	c.text.Append(string(c.char))
	return c
}

// Cos adds a cos function to the input text.
// It uses the function method to add the function.

// DecimalPoint adds a decimal point to the input text.
// Only one decimal point is allowed per number.
func (c *cursor) DecimalPoint() Cursor {
	c.prepare()
	defer c.exhaust()

	if runes.IsValid(c.text.First()) && runes.IsDigit(c.text.Last()) && !runes.IsDotted(c.text) {
		c.text.Append(".")
	}

	c.text.Append(string(c.char))
	return c
}

// Delete removes the last character from the input text.
// If the input text is "ANS", it is cleared.
// If the input was a function call, the whole function is removed.
// Otherwise, the last character is removed.
func (c *cursor) Delete() Cursor {
	c.prepare()
	defer c.exhaust()

	switch {
	case c.abortFunction(): // aborted function, nothing to do

	case c.text.Equals("ANS"): // clear screen
		c.text.Clear()

	default: // remove last character
		c.text.Backspace()

	}

	c.text.Append(string(c.char))
	return c
}

// Error displays NaN in the input text.
// Check() can be used to retrieve the error.
func (c *cursor) Error(err error) Cursor {
	c.err = err
	c.text.Clear()
	c.text.Append(fmt.Sprintf("%g", float32(math.NaN())))
	return c
}

// Evaluate evaluates the input text and updates the memory cell accordingly.
// It uses the parser to evaluate the input text.
func (c *cursor) Equals() Cursor {
	c.prepare()
	defer c.exhaust()

	// complete operation if the last character is an operator
	// return for confirmation if completion was needed
	if c.abortOperation() || c.abortFunction() || c.closeBrackets() {
		c.text.Append(string(c.char))
		return c
	}

	// do nothing if the begin is invalid
	if !runes.IsValid(c.text.First()) {
		c.text.Append(string(c.char))
		return c
	}

	// define options
	options := []parser.Option{
		parser.WithVar("ANS", c.memory.Get()),
		parser.WithVar("pi", big.NewFloat(math.Pi)),
		parser.WithVar("e", big.NewFloat(math.E)),
		parser.WithFunc("save", func(arg *big.Float) (*big.Float, error) {
			if err := c.memory.Set(arg); err != nil {
				return nil, err
			}

			return c.memory.Get(), nil
		}),
		parser.WithReplacements("×", "*", "÷", "/", "π", "pi"),
	}
	for name, fn := range map[string]func(float64) (float64, error){
		"sin": func(f float64) (float64, error) { return math.Sin(f), nil },
		"cos": func(f float64) (float64, error) { return math.Cos(f), nil },
		"tan": func(f float64) (float64, error) { return math.Tan(f), nil },
		"arcsin": func(f float64) (float64, error) {
			if f < -1 || f > 1 {
				return 0, fmt.Errorf("arcsin(%g) is undefined", f)
			}
			return math.Asin(f), nil
		},
		"arccos": func(f float64) (float64, error) {
			if f < -1 || f > 1 {
				return 0, fmt.Errorf("arccos(%g) is undefined", f)
			}
			return math.Acos(f), nil
		},
		"arctan": func(f float64) (float64, error) { return math.Atan(f), nil },
		"log": func(f float64) (float64, error) {
			if f <= 0 {
				return 0, fmt.Errorf("log(%g) is undefined", f)
			}
			return math.Log10(f), nil
		},
		"ln": func(f float64) (float64, error) {
			if f <= 0 {
				return 0, fmt.Errorf("ln(%g) is undefined", f)
			}
			return math.Log(f), nil
		},
	} {
		options = append(options, parser.WithFunc(name, fn))
	}

	// evaluate input text
	result, err := parser.NewParser(options...).Parse("save(" + c.text.String() + ")")
	if err != nil {
		return c.Error(err)
	}

	if result == nil {
		return c.Error(fmt.Errorf("no result"))
	}

	// display result
	c.text.Clear()
	c.text.Append(result.Text('f', -1))
	return c
}

// String returns the input text as a string.
func (c *cursor) String() string { return c.text.String() }

/*
Units and Unary Operators
*/
func (c *cursor) Degrees() Cursor   { return c.unit("°") }
func (c *cursor) Factorial() Cursor { return c.unit("!") }
func (c *cursor) Radians() Cursor   { return c.unit("÷1°") }

/*
Binary Operators
*/
func (c *cursor) Divide() Cursor { return c.binary('÷') }
func (c *cursor) Minus() Cursor  { return c.binary('-') }
func (c *cursor) Plus() Cursor   { return c.binary('+') }
func (c *cursor) Times() Cursor  { return c.binary('×') }

/*
Functions
*/
func (c *cursor) Cos() Cursor        { return c.function("cos") }
func (c *cursor) Ln() Cursor         { return c.function("ln") }
func (c *cursor) Log() Cursor        { return c.function("log") }
func (c *cursor) Power() Cursor      { return c.binary('^') }
func (c *cursor) Sin() Cursor        { return c.function("sin") }
func (c *cursor) Square() Cursor     { return c.binary('^', '2') }
func (c *cursor) SquareRoot() Cursor { return c.character('√') }
func (c *cursor) Tan() Cursor        { return c.function("tan") }

/*
Numbers and Constants
*/
func (c *cursor) Euler() Cursor { return c.character('e') }
func (c *cursor) Pi() Cursor    { return c.character('π') }
func (c *cursor) One() Cursor   { return c.character('1') }
func (c *cursor) Zero() Cursor  { return c.character('0') }
func (c *cursor) Two() Cursor   { return c.character('2') }
func (c *cursor) Three() Cursor { return c.character('3') }
func (c *cursor) Four() Cursor  { return c.character('4') }
func (c *cursor) Five() Cursor  { return c.character('5') }
func (c *cursor) Six() Cursor   { return c.character('6') }
func (c *cursor) Seven() Cursor { return c.character('7') }
func (c *cursor) Eight() Cursor { return c.character('8') }
func (c *cursor) Nine() Cursor  { return c.character('9') }

// New creates new cursor.
func New(text *runes.Sequence, memory memory.MemoryCell) Cursor {
	return &cursor{
		char:   '_',
		text:   text,
		memory: memory,
	}
}

// Do processes the given operator on the input text.
// It uses the cursor to perform the operation.
func Do(operator string, text *runes.Sequence, memory memory.MemoryCell) Cursor {
	c := New(text, memory)
	if fn, ok := map[string]func() Cursor{
		"<x":    c.Delete,
		"()":    c.Brackets,
		"AC":    c.Clear,
		"°":     c.Degrees,
		"1/°":   c.Radians,
		"÷":     c.Divide,
		"-":     c.Minus,
		"+":     c.Plus,
		"×":     c.Times,
		"0":     c.Zero,
		"1":     c.One,
		"2":     c.Two,
		"3":     c.Three,
		"4":     c.Four,
		"5":     c.Five,
		"6":     c.Six,
		"7":     c.Seven,
		"8":     c.Eight,
		"9":     c.Nine,
		".":     c.DecimalPoint,
		"=":     c.Equals,
		"π":     c.Pi,
		"!":     c.Factorial,
		"√":     c.SquareRoot,
		"e":     c.Euler,
		"e^":    func() Cursor { return c.Euler().Power() },
		"^":     c.Power,
		"x²":    c.Square,
		"10^":   func() Cursor { return c.One().Zero().Power() },
		"sin⁻¹": c.Arcsin,
		"cos⁻¹": c.Arccos,
		"tan⁻¹": c.Arctan,
		"sin":   c.Sin,
		"cos":   c.Cos,
		"tan":   c.Tan,
		"log":   c.Log,
		"ln":    c.Ln,
	}[operator]; ok {
		return fn()
	}

	return c.Error(fmt.Errorf("unknown operator: %s", operator))
}
