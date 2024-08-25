/*
Package cursor provides a cursor for a calculator.
It mimics the behavior of a physical calculator.

Example:

	c := cursor.New(runes.Sequence("_"), 0)
	c.Do("1") // or c.One() adds 1 to the input text
	c.Do("+") // or c.Plus() adds + to the input text
	c.Do("2") // or c.Two() adds 2 to the input text
	c.Do("=") // or c.Equals() evaluates the input text
	c.String() // returns the result of evaluation: "3"
*/
package cursor

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/sarumaj/edu-taschenrechner/pkg/parser"
	"github.com/sarumaj/edu-taschenrechner/pkg/runes"
)

// make sure that the cursor type implements the Cursor interface
var _ Cursor = (*cursor)(nil)

// Cursor is an actor performing operations.
// It defines the basic operations of a calculator.
type Cursor = CursorInterface[*cursor]

// CursorInterface is a generic interface for the cursor.
type CursorInterface[T any] interface {
	Divide() T
	Minus() T
	Plus() T
	Times() T

	Arccos() T
	Arcsin() T
	Arctan() T
	Cos() T
	Factorial() T
	Gdc() *cursor
	Lcm() *cursor
	Ln() T
	Log() T
	Power() T
	Sin() T
	Square() T
	SquareRoot() T
	Tan() T

	Euler() T
	Pi() T
	Zero() T
	One() T
	Two() T
	Three() T
	Four() T
	Five() T
	Six() T
	Seven() T
	Eight() T
	Nine() T

	Degrees() T
	Radians() T

	Brackets() T
	Cancel()
	Clear() T
	Comma() T
	DecimalPoint() T
	Delete() T
	Do(string) T
	Error(error) T
	Equals() T
	EqualsWithFormat(format byte) T

	Check() error
	String() string
}

// cursor is a cursor of operations.
// Currently, it only supports the basic operations of a calculator.
type cursor struct {
	ctx    context.Context
	cancel context.CancelFunc
	err    error
	char   rune
	ready  bool
	parser parser.Parser
	text   *runes.Sequence
}

// abortFunction aborts the function if the last character is an opening bracket.
// It is used to prevent incomplete functions.
func (c *cursor) abortFunction() (aborted bool) {
	if !c.text.EndsWith(",") && (!c.text.EndsWith("(") || c.text.Shift().EndsWith("(")) {
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
func (c *cursor) binary(op rune, opts ...rune) *cursor {
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
func (c *cursor) character(v rune, opts ...rune) *cursor {
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
func (c *cursor) function(name string) *cursor {
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
func (c *cursor) unit(u string) *cursor {
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
		(runes.IsDigit(c.text.Last()) || c.text.Equals("ANS") || runes.IsAnyOf(c.text.Last(), ")!°πe")) {

		c.text.Append(u)
	}

	c.text.Append(string(c.char))
	return c
}

// Brackets adds brackets to the input text.
// If the last character is an opening bracket, a closing bracket is added.
// If the last character is a digit or a closing bracket, a multiplication operator and an opening bracket are added.
// If the last character is an operator, a closing bracket is added.
func (c *cursor) Brackets() *cursor {
	c.prepare()
	defer c.exhaust()

	switch {
	case // just open
		runes.IsAnyOf(c.text.Last(), "(+-×÷√^,"),
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

// Cancel cancels the cursor.
func (c *cursor) Cancel() {
	if c.cancel != nil {
		c.cancel()
	}
}

// Check returns the last error.
func (c *cursor) Check() error { return c.err }

// Clear clears the input text.
func (c *cursor) Clear() *cursor {
	c.text.Clear()
	c.text.Append(string(c.char))
	return c
}

// Comma adds a comma to the input text.
func (c *cursor) Comma() *cursor {
	c.prepare()
	defer c.exhaust()

	backup := c.text.Copy()
	for runes.IsDigit(c.text.Last()) || runes.IsAnyOf(c.text.Last(), "eπ.,") || (c.text.EndsWith(")") && runes.HowManyOpen(c.text) > 0) {
		c.text.Backspace()
	}

	if !c.abortFunction() {
		*c.text = backup
		c.text.Append(string(c.char))
		return c
	}

	*c.text = backup
	if runes.IsValid(c.text.First()) && runes.IsDigit(c.text.Last()) {
		c.text.Append(",")
	}

	c.text.Append(string(c.char))
	return c
}

// DecimalPoint adds a decimal point to the input text.
// Only one decimal point is allowed per number.
func (c *cursor) DecimalPoint() *cursor {
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
func (c *cursor) Delete() *cursor {
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

// Do processes the given operator on the input text.
// It uses the cursor to perform the operation.
func (c *cursor) Do(operator string) *cursor {
	if fn, ok := map[string]func() *cursor{
		"↩":     c.Delete,
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
		",":     c.Comma,
		"=":     c.Equals,
		"π":     c.Pi,
		"!":     c.Factorial,
		"√":     c.SquareRoot,
		"e":     c.Euler,
		"eⁿ":    func() *cursor { return c.character('e', '^') },
		"xⁿ":    c.Power,
		"x²":    c.Square,
		"10ⁿ":   func() *cursor { return c.character('1', '0', '^') },
		"sin⁻¹": c.Arcsin,
		"cos⁻¹": c.Arccos,
		"tan⁻¹": c.Arctan,
		"sin":   c.Sin,
		"cos":   c.Cos,
		"tan":   c.Tan,
		"log":   c.Log,
		"ln":    c.Ln,
		"gdc":   c.Gdc,
		"lcm":   c.Lcm,
	}[operator]; ok {
		return fn()
	}

	return c.Error(fmt.Errorf("unknown operator: %s", operator))
}

// Error displays NaN in the input text.
// Check() can be used to retrieve the error.
func (c *cursor) Error(err error) *cursor {
	c.err = err
	c.text.Clear()
	c.text.Append(fmt.Sprintf("%g", float32(math.NaN())))
	return c
}

// Equals evaluates the input text and displays the result.
// It uses the EqualsWithFormat method to display the result using the default format.
func (c *cursor) Equals() *cursor { return c.EqualsWithFormat('f') }

// EqualsWithFormat evaluates the input text and displays the result.
func (c *cursor) EqualsWithFormat(format byte) *cursor {
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

	// evaluate input text
	result, err := c.parser.Parse(c.ctx, "save("+c.text.String()+")")
	if err != nil {
		return c.Error(err)
	}

	if result == nil {
		return c.Error(fmt.Errorf("no result"))
	}

	// display result
	c.text.Clear()
	c.text.Append(result.Text(format, -1))
	return c
}

// String returns the input text as a string.
func (c *cursor) String() string { return c.text.String() }

/*
Units and Unary Operators
*/
func (c *cursor) Degrees() *cursor   { return c.unit("°") }
func (c *cursor) Factorial() *cursor { return c.unit("!") }
func (c *cursor) Radians() *cursor   { return c.unit("÷1°") }

/*
Binary Operators
*/
func (c *cursor) Divide() *cursor { return c.binary('÷') }
func (c *cursor) Minus() *cursor  { return c.binary('-') }
func (c *cursor) Plus() *cursor   { return c.binary('+') }
func (c *cursor) Times() *cursor  { return c.binary('×') }

/*
Functions
*/

func (c *cursor) Arccos() *cursor     { return c.function("arccos") }
func (c *cursor) Arcsin() *cursor     { return c.function("arcsin") }
func (c *cursor) Arctan() *cursor     { return c.function("arctan") }
func (c *cursor) Cos() *cursor        { return c.function("cos") }
func (c *cursor) Gdc() *cursor        { return c.function("gdc") }
func (c *cursor) Lcm() *cursor        { return c.function("lcm") }
func (c *cursor) Ln() *cursor         { return c.function("ln") }
func (c *cursor) Log() *cursor        { return c.function("log") }
func (c *cursor) Power() *cursor      { return c.binary('^') }
func (c *cursor) Sin() *cursor        { return c.function("sin") }
func (c *cursor) Square() *cursor     { return c.binary('^', '2') }
func (c *cursor) SquareRoot() *cursor { return c.character('√') }
func (c *cursor) Tan() *cursor        { return c.function("tan") }

/*
Numbers and Constants
*/
func (c *cursor) Euler() *cursor { return c.character('e') }
func (c *cursor) Pi() *cursor    { return c.character('π') }
func (c *cursor) One() *cursor   { return c.character('1') }
func (c *cursor) Zero() *cursor  { return c.character('0') }
func (c *cursor) Two() *cursor   { return c.character('2') }
func (c *cursor) Three() *cursor { return c.character('3') }
func (c *cursor) Four() *cursor  { return c.character('4') }
func (c *cursor) Five() *cursor  { return c.character('5') }
func (c *cursor) Six() *cursor   { return c.character('6') }
func (c *cursor) Seven() *cursor { return c.character('7') }
func (c *cursor) Eight() *cursor { return c.character('8') }
func (c *cursor) Nine() *cursor  { return c.character('9') }

// New creates new cursor.
func New(text *runes.Sequence, timeout time.Duration, parserOpts ...parser.Option) Cursor {
	c := cursor{
		char:   '_',
		text:   text,
		parser: parser.NewParser(parserOpts...),
	}

	if timeout > 0 {
		c.ctx, c.cancel = context.WithTimeout(context.Background(), timeout)
	} else {
		c.ctx, c.cancel = context.WithCancel(context.Background())
	}

	return &c
}
