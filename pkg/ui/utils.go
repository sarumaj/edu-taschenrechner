package ui

import (
	"fmt"
	"math"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/maja42/goval"
	"github.com/sarumaj/edu-taschenrechner/pkg/memory"
	"github.com/sarumaj/edu-taschenrechner/pkg/runes"
)

// doBtnClick changes the text of the label based on the button type.
func doBtnClick(label *widget.Label, memory memory.MemoryCell, btnId string) func() {
	return func() {
		text := runes.NewInput(label.Text)
		modifyInput(text, memory, btnId)
		label.SetText(text.String())
	}
}

// getButtons provides a selector to select buttons for a map based on their type.
func getButtons(objects map[string]fyne.CanvasObject) func(in ...string) (out []fyne.CanvasObject) {
	return func(in ...string) (out []fyne.CanvasObject) {
		for _, c := range in {
			if v, ok := objects[c].(*widget.Button); ok {
				out = append(out, v)
			}
		}

		return
	}
}

// modifyInput interacts with the input field based on the ID of the button
func modifyInput(text *runes.Input, memory memory.MemoryCell, btnId string) {
	if text.EndsWith("_") { // remove cursor
		text.Backspace()
	} else if text.Equals(fmt.Sprint(math.NaN())) { // clear screen
		text.Clear()
	} else { // reuse result from memory cell
		text.Clear()
		text.Append("ANS")
	}

	switch btnId {
	case "()": // open or close brackets
		switch {
		case runes.IsAnyOf(text.Last(), "(+-×÷"), !runes.IsValid(text.Last()):
			text.Append("(") // just open

		case runes.HowManyOpen(text) > 0 && (runes.IsDigit(text.Last()) || text.EndsWith(")")):
			text.Append(")") // just close

		case text.EndsWith(")") || runes.IsDigit(text.Last()) || text.Equals("ANS"):
			text.Append("×(") // multiply and open

		}

		text.Append("_")

	case "<x": // remove trailing
		if text.Equals("ANS") {
			text.Clear()
		} else {
			text.Backspace()
		}

		text.Append("_")

	case ".": // set decimal point
		if !runes.IsDotted(text) && runes.IsDigit(text.Last()) && !text.Equals("ANS") {
			text.Append(".")
		}

		text.Append("_")

	case "AC": // clear screen
		text.Clear()
		text.Append("_")

	case "+", "×", "÷":
		// undo last operation or remove decimal point
		if text.Last() == '.' || runes.IsAnyOf(text.Last(), "×÷+-") {
			text.Backspace()
		}

		// set operator
		if runes.IsDigit(text.Last()) || text.Equals("ANS") || text.EndsWith(")") {
			text.Append(btnId)
		}

		text.Append("_")

	case "-": // handle +- conversions and set operator
		if text.EndsWith("-") { // flip sign
			text.Backspace()
			text.Append("+")

		} else {
			text.Append("-")
		}

		if text.BeginsWith("+") { // remove if leading
			text.Delete()
		}

		if runes.IsAnyOf(text.Shift().Last(), "+×÷") && text.EndsWith("+") { // remove if flipped too much
			text.Backspace()
		}

		text.Append("_")

	case "=": // calculate
		// abort incomplete operation
		for runes.IsAnyOf(text.Last(), "(+-×÷") {
			text.Backspace()
		}

		// close all opened brackets
		for i, o := 0, runes.HowManyOpen(text); i < o; i++ {
			text.Append(")")
		}

		// do nothing if the begin is invalid
		if !runes.IsValid(text.First()) {
			text.Append("_")
			return
		}

		// evaluate
		eval := goval.NewEvaluator()
		result, err := eval.Evaluate(
			// replace operators
			strings.NewReplacer("×", "*", "÷", "/").Replace(
				"save("+ // save to memory cell
					"1.0*"+ // enforce decimal result
					text.String()+
					")",
			),
			// define variables
			map[string]any{"ANS": memory.Get()},
			// define functions
			map[string]func(args ...any) (any, error){
				"save": func(args ...any) (any, error) {
					memory.Set(args[0])
					return memory.Get(), nil
				},
			},
		)

		// handle error
		if err != nil {
			text.Clear()
			text.Append(fmt.Sprintf("%g", math.NaN()))
			return
		}

		// display result
		text.Clear()
		text.Append(fmt.Sprint(result))

	default: // handle digits
		// multiply if behind closing bracket or memory cell value
		if text.EndsWith(")") || text.Equals("ANS") {
			text.Append("×")
		}

		text.Append(btnId + "_")

	}
}
