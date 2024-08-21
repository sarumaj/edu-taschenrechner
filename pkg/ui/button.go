package ui

import (
	"fyne.io/fyne/driver/mobile"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/sarumaj/edu-taschenrechner/pkg/eval"
	"github.com/sarumaj/edu-taschenrechner/pkg/runes"
)

// Custom button widget that implements the Touchable interface for future usage on mobile devices.
type Button struct{ widget.Button }

// CreateRenderer returns a custom renderer for the Button widget.
func (b *Button) CreateRenderer() fyne.WidgetRenderer {
	return &buttonRenderer{WidgetRenderer: b.Button.CreateRenderer()}
}

// Implement Touchable interface for Button.
var _ mobile.Touchable = &Button{}

func (b *Button) TouchCancel(*mobile.TouchEvent) {}
func (b *Button) TouchDown(*mobile.TouchEvent)   { b.Button.OnTapped() }
func (b *Button) TouchUp(*mobile.TouchEvent)     {}

// buttonRenderer is a custom renderer for the Button widget.
// Required to implement the Touchable interface.
type buttonRenderer struct{ fyne.WidgetRenderer }

// Create a new button widget with the given text and tapped function.
func NewButton(text string, display *Display) *Button {
	btn := &Button{Button: widget.Button{
		Text: text,
		OnTapped: func() {
			input := runes.NewInput(display.Text)
			eval.Evaluate(text, input, display.MemoryCell)
			display.SetText(input.String())
		},
	}}
	btn.ExtendBaseWidget(btn)
	return btn
}
