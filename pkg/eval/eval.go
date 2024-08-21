package eval

import (
	"github.com/sarumaj/edu-taschenrechner/pkg/memory"
	"github.com/sarumaj/edu-taschenrechner/pkg/runes"
)

// Evaluate processes the given operator on the input text.
func Evaluate(operator string, text *runes.Input, memory memory.MemoryCell) {
	switch c := Begin(text, memory); operator {
	case "<x":
		c.Delete()

	case "()":
		c.Brackets()

	case "AC":
		c.Clear()

	case "÷":
		c.Divide()

	case "-":
		c.Minus()

	case "+":
		c.Plus()

	case "×":
		c.Times()

	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		c.Digit([]rune(operator)[0])

	case ".":
		c.DecimalPoint()

	case "=":
		c.Equals()

	default:
		c.Error()

	}
}
