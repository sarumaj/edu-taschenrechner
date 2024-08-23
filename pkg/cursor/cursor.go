package cursor

import (
	"fmt"
	"math"
	"math/big"
	"strings"

	"github.com/sarumaj/edu-taschenrechner/pkg/memory"
	"github.com/sarumaj/edu-taschenrechner/pkg/parser"
	"github.com/sarumaj/edu-taschenrechner/pkg/runes"
)

// Cursor is an actor performing operations.
// It defines the basic operations of a calculator.
type Cursor interface {
	Brackets() Cursor
	Clear() Cursor
	DecimalPoint() Cursor
	Delete() Cursor
	Divide() Cursor
	Eight() Cursor
	Error() Cursor
	Equals() Cursor
	Five() Cursor
	Four() Cursor
	Minus() Cursor
	Nine() Cursor
	One() Cursor
	Plus() Cursor
	Seven() Cursor
	Six() Cursor
	String() string
	Three() Cursor
	Times() Cursor
	Two() Cursor
	Zero() Cursor
}

// cursor is a cursor of operations.
// Currently, it only supports the basic operations of a calculator.
type cursor struct {
	char   rune
	ready  bool
	text   *runes.Sequence
	memory memory.MemoryCell
}

// Construct adds a plus, times, or divide operator to the input text.
// If the last character is a decimal point or an operator, it is removed.
// If the last character is a digit or the result of a previous calculation, the operator is added.
func (c *cursor) addMultiplyDivide(op rune) Cursor {
	c.prepare()
	defer c.exhaust()

	if !runes.IsAnyOf(op, "×÷+") {
		return c.Error()
	}

	// undo last operation or remove decimal point
	if c.text.Last() == '.' || runes.IsAnyOf(c.text.Last(), "×÷+-") {
		c.text.Backspace()
	}

	// set operator
	if runes.IsDigit(c.text.Last()) || c.text.Equals("ANS") || c.text.EndsWith(")") {
		c.text.Append(string(op))
	}

	c.text.Append(string(c.char))
	return c
}

// digit adds a digit to the input text.
// If the last character is a closing bracket or the result of a previous calculation, a multiplication operator is added.
// If input is not a digit, an error is displayed.
func (c *cursor) digit(digit rune) Cursor {
	c.prepare()
	defer c.exhaust()

	if !runes.IsDigit(digit) {
		return c.Error()
	}

	// multiply if behind closing bracket or memory cell value
	if c.text.EndsWith(")") || c.text.Equals("ANS") {
		c.text.Append("×")
	}

	c.text.Append(string(digit))
	c.text.Append(string(c.char))
	return c
}

// exhaust disables the cursor and returns it.
// It is used to prevent further operations on the cursor.
// The cursor is disabled after an operation is completed.
// To re-enable the cursor, a prepare operation is required.
func (c *cursor) exhaust() { c.ready = false }

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

// Brackets adds brackets to the input text.
// If the last character is an opening bracket, a closing bracket is added.
// If the last character is a digit or a closing bracket, a multiplication operator and an opening bracket are added.
// If the last character is an operator, a closing bracket is added.
func (c *cursor) Brackets() Cursor {
	c.prepare()
	defer c.exhaust()

	switch {
	case runes.IsAnyOf(c.text.Last(), "(+-×÷"), !runes.IsValid(c.text.Last()):
		c.text.Append("(") // just open

	case runes.HowManyOpen(c.text) > 0 && (runes.IsDigit(c.text.Last()) || c.text.EndsWith(")")):
		c.text.Append(")") // just close

	case c.text.EndsWith(")") || runes.IsDigit(c.text.Last()) || c.text.Equals("ANS"):
		c.text.Append("×(") // multiply and open

	}

	c.text.Append(string(c.char))
	return c
}

// Clear clears the input text.
func (c *cursor) Clear() Cursor {
	c.text.Clear()
	c.text.Append(string(c.char))
	return c
}

// DecimalPoint adds a decimal point to the input text.
// Only one decimal point is allowed per number.
func (c *cursor) DecimalPoint() Cursor {
	c.prepare()
	defer c.exhaust()

	if !runes.IsDotted(c.text) && runes.IsDigit(c.text.Last()) && !c.text.Equals("ANS") {
		c.text.Append(".")
	}

	c.text.Append(string(c.char))
	return c
}

// Delete removes the last character from the input text.
// If the input text is "ANS", it is cleared.
// Otherwise, the last character is removed.
func (c *cursor) Delete() Cursor {
	c.prepare()
	defer c.exhaust()

	if c.text.Equals("ANS") {
		c.text.Clear()
	} else {
		c.text.Backspace()
	}

	c.text.Append(string(c.char))
	return c
}

// Divide adds a divide operator to the input text.
// It uses the addMultiplyDivide method to add the operator.
func (c *cursor) Divide() Cursor { return c.addMultiplyDivide('÷') }

// Eight adds an eight to the input text.
// It uses the digit method to add the digit.
func (c *cursor) Eight() Cursor { return c.digit('8') }

// Error displays NaN in the input text.
func (c *cursor) Error() Cursor {
	c.text.Clear()
	c.text.Append(fmt.Sprintf("%g", float32(math.NaN())))
	return c
}

// Evaluate evaluates the input text and updates the memory cell accordingly.
// It uses the goval package to evaluate the input text.
// The result is saved to the memory cell and displayed in the input text.
// If an error occurs during evaluation, the input text is cleared and NaN is displayed.
func (c *cursor) Equals() Cursor {
	c.prepare()
	defer c.exhaust()

	// abort incomplete operation
	for runes.IsAnyOf(c.text.Last(), "(+-×÷") {
		c.text.Backspace()
	}

	// close all opened brackets
	for i, o := 0, runes.HowManyOpen(c.text); i < o; i++ {
		c.text.Append(")")
	}

	// do nothing if the begin is invalid
	if !runes.IsValid(c.text.First()) {
		c.text.Append(string(c.char))
		return c
	}

	// evaluate input text
	result, err := parser.NewParser(
		parser.WithVar("ANS", c.memory.Get()),
		parser.WithFunc("save", func(args ...*big.Float) (*big.Float, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("save function requires exactly 1 argument")
			}

			if err := c.memory.Set(args[0]); err != nil {
				return nil, err
			}

			return c.memory.Get(), nil
		}),
	).Parse(strings.NewReplacer("×", "*", "÷", "/").Replace("save(" + c.text.String() + ")"))

	// handle error
	if err != nil {
		return c.Error()
	}

	// display result
	c.text.Clear()
	c.text.Append(result.Text('g', -1))
	return c
}

// Five adds a five to the input text.
// It uses the digit method to add the digit.
func (c *cursor) Five() Cursor { return c.digit('5') }

// Four adds a four to the input text.
// It uses the digit method to add the digit.
func (c *cursor) Four() Cursor { return c.digit('4') }

// Minus adds a minus operator to the input text.
// If the last character is a minus operator, it is flipped to a plus operator.
// If the last character is a digit or a closing bracket, a minus operator is added.
// If the last character is a plus, times, or divide operator, the minus operator is added.
// If the input text begins with a plus operator, it is removed.
// If the input text ends with a previously flipped minus operator, it is removed.
func (c *cursor) Minus() Cursor {
	c.prepare()
	defer c.exhaust()

	if c.text.EndsWith("-") { // flip sign
		c.text.Backspace()
		c.text.Append("+")

	} else {
		c.text.Append("-")
	}

	if c.text.BeginsWith("+") { // remove if leading
		c.text.Delete()
	}

	if runes.IsAnyOf(c.text.Shift().Last(), "+×÷") && c.text.EndsWith("+") { // remove if flipped too much
		c.text.Backspace()
	}

	c.text.Append(string(c.char))
	return c
}

// Nine adds a nine to the input text.
// It uses the digit method to add the digit.
func (c *cursor) Nine() Cursor { return c.digit('9') }

// One adds a one to the input text.
// It uses the digit method to add the digit.
func (c *cursor) One() Cursor { return c.digit('1') }

// Plus adds a plus operator to the input text.
// It uses the addMultiplyDivide method to add the operator.
func (c *cursor) Plus() Cursor { return c.addMultiplyDivide('+') }

// Seven adds a seven to the input text.
// It uses the digit method to add the digit.
func (c *cursor) Seven() Cursor { return c.digit('7') }

// Six adds a six to the input text.
// It uses the digit method to add the digit.
func (c *cursor) Six() Cursor { return c.digit('6') }

// String returns the input text as a string.
func (c *cursor) String() string { return c.text.String() }

// Three adds a three to the input text.
// It uses the digit method to add the digit.
func (c *cursor) Three() Cursor { return c.digit('3') }

// Times adds a times operator to the input text.
// It uses the addMultiplyDivide method to add the operator.
func (c *cursor) Times() Cursor { return c.addMultiplyDivide('×') }

// Two adds a two to the input text.
// It uses the digit method to add the digit.
func (c *cursor) Two() Cursor { return c.digit('2') }

// Zero adds a zero to the input text.
// It uses the digit method to add the digit.
func (c *cursor) Zero() Cursor { return c.digit('0') }

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
		"<x": c.Delete,
		"()": c.Brackets,
		"AC": c.Clear,
		"÷":  c.Divide,
		"-":  c.Minus,
		"+":  c.Plus,
		"×":  c.Times,
		"0":  c.Zero,
		"1":  c.One,
		"2":  c.Two,
		"3":  c.Three,
		"4":  c.Four,
		"5":  c.Five,
		"6":  c.Six,
		"7":  c.Seven,
		"8":  c.Eight,
		"9":  c.Nine,
		".":  c.DecimalPoint,
		"=":  c.Equals,
	}[operator]; ok {
		return fn()
	}

	return c.Error()
}
