package ui

import (
	"fyne.io/fyne/driver/mobile"
	"fyne.io/fyne/v2/widget"
)

// Custom button widget that implements the Touchable interface for future usage on mobile devices.
type Button struct{ widget.Button }

// Implement Touchable interface for Button.
var _ mobile.Touchable = &Button{}

func (b *Button) TouchCancel(*mobile.TouchEvent) {}
func (b *Button) TouchDown(*mobile.TouchEvent)   { b.Button.OnTapped() }
func (b *Button) TouchUp(*mobile.TouchEvent)     {}

// Create a new button widget with the given text and tapped function.
func NewButton(text string, display *Display) *Button {
	btn := &Button{Button: widget.Button{
		Text:     text,
		OnTapped: func() { display.SetText(text) },
	}}
	btn.ExtendBaseWidget(btn)
	return btn
}
