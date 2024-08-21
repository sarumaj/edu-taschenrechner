package eval

import (
	"fmt"
	"math"
	"strings"

	"github.com/maja42/goval"
	"github.com/sarumaj/edu-taschenrechner/pkg/memory"
	"github.com/sarumaj/edu-taschenrechner/pkg/runes"
)

// Chain is a chain of operations.
// It defines the basic operations of a calculator.
type Chain interface {
	Brackets()
	Clear()
	DecimalPoint()
	Delete()
	Digit(digit rune)
	Divide()
	Error()
	Equals()
	Minus()
	Plus()
	Times()
}

// chain is a chain of operations.
// Currently, it only supports the basic operations of a calculator.
type chain struct {
	text   *runes.Input
	memory memory.MemoryCell
}

// Brackets adds brackets to the input text.
// If the last character is an opening bracket, a closing bracket is added.
// If the last character is a digit or a closing bracket, a multiplication operator and an opening bracket are added.
// If the last character is an operator, a closing bracket is added.
func (c *chain) Brackets() {
	switch {
	case runes.IsAnyOf(c.text.Last(), "(+-×÷"), !runes.IsValid(c.text.Last()):
		c.text.Append("(") // just open

	case runes.HowManyOpen(c.text) > 0 && (runes.IsDigit(c.text.Last()) || c.text.EndsWith(")")):
		c.text.Append(")") // just close

	case c.text.EndsWith(")") || runes.IsDigit(c.text.Last()) || c.text.Equals("ANS"):
		c.text.Append("×(") // multiply and open

	}

	c.text.Append("_")
}

// Clear clears the input text.
func (c *chain) Clear() {
	c.text.Clear()
	c.text.Append("_")
}

// DecimalPoint adds a decimal point to the input text.
// Only one decimal point is allowed per number.
func (c *chain) DecimalPoint() {
	if !runes.IsDotted(c.text) && runes.IsDigit(c.text.Last()) && !c.text.Equals("ANS") {
		c.text.Append(".")
	}

	c.text.Append("_")
}

// Delete removes the last character from the input text.
// If the input text is "ANS", it is cleared.
// Otherwise, the last character is removed.
func (c *chain) Delete() {
	if c.text.Equals("ANS") {
		c.text.Clear()
	} else {
		c.text.Backspace()
	}

	c.text.Append("_")
}

// Digit adds a digit to the input text.
// If the last character is a closing bracket or the result of a previous calculation, a multiplication operator is added.
// If input is not a digit, an error is displayed.
func (c *chain) Digit(digit rune) {
	if !runes.IsDigit(digit) {
		c.Error()
	}

	// multiply if behind closing bracket or memory cell value
	if c.text.EndsWith(")") || c.text.Equals("ANS") {
		c.text.Append("×")
	}

	c.text.Append(string(digit))
	c.text.Append("_")
}

// Divide adds a divide operator to the input text.
// If the last character is a decimal point or an operator, it is removed.
// If the last character is a digit or the result of a previous calculation, a divide operator is added.
// If the last character is a closing bracket, a divide operator is added.
func (c *chain) Divide() {
	// undo last operation or remove decimal point
	if c.text.Last() == '.' || runes.IsAnyOf(c.text.Last(), "×÷+-") {
		c.text.Backspace()
	}

	// set operator
	if runes.IsDigit(c.text.Last()) || c.text.Equals("ANS") || c.text.EndsWith(")") {
		c.text.Append("÷")
	}

	c.text.Append("_")
}

// Error displays NaN in the input text.
func (c *chain) Error() {
	c.text.Clear()
	c.text.Append(fmt.Sprintf("%g", float32(math.NaN())))
}

// Evaluate evaluates the input text and updates the memory cell accordingly.
// It uses the goval package to evaluate the input text.
// The result is saved to the memory cell and displayed in the input text.
// If an error occurs during evaluation, the input text is cleared and NaN is displayed.
func (c *chain) Equals() {
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
		c.text.Append("_")
		return
	}

	// evaluate input text
	result, err := goval.NewEvaluator().Evaluate(
		// replace operators
		strings.NewReplacer("×", "*", "÷", "/").Replace(
			"save("+ // save to memory cell
				"1.0*"+ // enforce decimal result
				c.text.String()+
				")",
		),
		// define variables
		map[string]any{"ANS": c.memory.Get()},
		// define functions
		map[string]func(args ...any) (any, error){
			"save": func(args ...any) (any, error) {
				if err := c.memory.Set(args[0]); err != nil {
					return nil, err
				}
				return c.memory.Get(), nil
			},
		},
	)

	// handle error
	if err != nil {
		c.Error()
		return
	}

	// display result
	c.text.Clear()
	c.text.Append(fmt.Sprint(result))
}

// Minus adds a minus operator to the input text.
// If the last character is a minus operator, it is flipped to a plus operator.
// If the last character is a digit or a closing bracket, a minus operator is added.
// If the last character is a plus, times, or divide operator, the minus operator is added.
// If the input text begins with a plus operator, it is removed.
// If the input text ends with a previously flipped minus operator, it is removed.
func (c *chain) Minus() {
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

	c.text.Append("_")
}

// Plus adds a plus operator to the input text.
// If the last character is a decimal point or an operator, it is removed.
// If the last character is a digit or the result of a previous calculation, a plus operator is added.
func (c *chain) Plus() {
	// undo last operation or remove decimal point
	if c.text.Last() == '.' || runes.IsAnyOf(c.text.Last(), "×÷+-") {
		c.text.Backspace()
	}

	// set operator
	if runes.IsDigit(c.text.Last()) || c.text.Equals("ANS") || c.text.EndsWith(")") {
		c.text.Append("+")
	}

	c.text.Append("_")
}

// Times adds a times operator to the input text.
// If the last character is a decimal point or an operator, it is removed.
// If the last character is a digit or the result of a previous calculation, a times operator is added.
// If the last character is a closing bracket, a times operator is added.
func (c *chain) Times() {
	// undo last operation or remove decimal point
	if c.text.Last() == '.' || runes.IsAnyOf(c.text.Last(), "×÷+-") {
		c.text.Backspace()
	}

	// set operator
	if runes.IsDigit(c.text.Last()) || c.text.Equals("ANS") || c.text.EndsWith(")") {
		c.text.Append("×")
	}

	c.text.Append("_")
}

// Begin prepares the input text for a new calculation.
// If the input text ends with a cursor, it is removed.
// If the input text equals NaN, the screen is cleared.
// Otherwise, the result from the memory cell is reused.
func Begin(text *runes.Input, memory memory.MemoryCell) Chain {
	if text.EndsWith("_") { // remove cursor
		text.Backspace()
	} else if text.Equals(fmt.Sprint(math.NaN())) { // clear screen
		text.Clear()
	} else { // reuse result from memory cell
		text.Clear()
		text.Append("ANS")
	}

	return &chain{
		text:   text,
		memory: memory,
	}
}
