package eval

import (
	"github.com/sarumaj/edu-taschenrechner/pkg/memory"
	"github.com/sarumaj/edu-taschenrechner/pkg/runes"
)

// Evaluate processes the operator and operand and updates the input text and memory cell accordingly.
func Evaluate(operator string, text *runes.Input, memory memory.MemoryCell) {
	switch chain := Begin(text, memory); operator {
	case "<x":
		chain.Delete()

	case "()":
		chain.Brackets()

	case "AC":
		chain.Clear()

	case "÷":
		chain.Divide()

	case "-":
		chain.Minus()

	case "+":
		chain.Plus()

	case "×":
		chain.Times()

	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		chain.Digit([]rune(operator)[0])

	case ".":
		chain.DecimalPoint()

	case "=":
		chain.Equals()

	default:
		chain.Error()

	}
}
