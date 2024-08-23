package ui

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/sarumaj/edu-taschenrechner/pkg/cursor"
	"github.com/sarumaj/edu-taschenrechner/pkg/memory"
	"github.com/sarumaj/edu-taschenrechner/pkg/runes"
)

// Display is a custom label widget that extends the default label with a memory cell.
type Display struct {
	widget.Entry
	memory.MemoryCell
}

// SetText sets the text of the display widget.
func (m *Display) SetText(text string) {
	m.Entry.SetText(cursor.Do(strings.TrimSpace(text), runes.NewSequence(m.Text), m.MemoryCell).String())
}

// NewDisplay creates a new label widget with the given text and memory cell.
func NewDisplay(text string) *Display {
	display := &Display{
		Entry: widget.Entry{
			Wrapping:  fyne.TextWrapOff,
			Text:      text,
			TextStyle: fyne.TextStyle{Monospace: true},
		},
		MemoryCell: memory.NewMemoryCell(),
	}
	display.ExtendBaseWidget(display)
	return display
}
