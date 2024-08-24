package ui

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/sarumaj/edu-taschenrechner/pkg/cursor"
	"github.com/sarumaj/edu-taschenrechner/pkg/memory"
	"github.com/sarumaj/edu-taschenrechner/pkg/runes"
)

// Display is a custom label widget that extends the default label with a memory cell.
type Display struct {
	widget.Entry
	memory.MemoryCell
	window fyne.Window
}

// SetText sets the text of the display widget.
// It moves the cursor to the end of the text and checks the state of the cursor.
// If the cursor is in an invalid state, it shows an error dialog.
func (m *Display) SetText(text string) {
	// move cursor by submitting the text to the cursor
	state := cursor.Do(strings.TrimSpace(text), runes.NewSequence(m.Text), m.MemoryCell)
	m.Entry.SetText(state.String())

	// set cursor to the end of the text
	m.Entry.CursorRow = len(m.Text) - 1
	m.Entry.CursorColumn = len(m.Text)
	m.Entry.Refresh()

	// check the state of the cursor
	if err := state.Check(); err != nil {
		dialog.ShowError(err, m.window)
	}
}

// NewDisplay creates a new label widget with the given text and memory cell.
func NewDisplay(text string, window fyne.Window) *Display {
	display := &Display{
		Entry: widget.Entry{
			Wrapping:  fyne.TextWrapOff,
			Text:      text,
			TextStyle: fyne.TextStyle{Monospace: true},
		},
		MemoryCell: memory.NewMemoryCell(),
		window:     window,
	}
	display.ExtendBaseWidget(display)
	return display
}
