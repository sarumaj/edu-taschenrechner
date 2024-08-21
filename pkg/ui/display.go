package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/sarumaj/edu-taschenrechner/pkg/memory"
)

// Display is a custom label widget that extends the default label with a memory cell.
type Display struct {
	widget.Label
	memory.MemoryCell
}

// CreateRenderer returns a custom renderer for the Display widget.
func (m *Display) CreateRenderer() fyne.WidgetRenderer {
	return &DisplayRenderer{WidgetRenderer: m.Label.CreateRenderer()}
}

// DisplayRenderer is a custom renderer for the Display widget.
type DisplayRenderer struct{ fyne.WidgetRenderer }

// NewDisplay creates a new label widget with the given text and memory cell.
func NewDisplay(text string) *Display {
	display := &Display{
		Label: widget.Label{
			Alignment:  fyne.TextAlignCenter,
			Truncation: fyne.TextTruncateEllipsis,
			Wrapping:   fyne.TextWrapOff,
			Text:       text,
			TextStyle:  fyne.TextStyle{Monospace: true},
		},
		MemoryCell: memory.NewMemoryCell(),
	}
	display.ExtendBaseWidget(display)
	return display
}
