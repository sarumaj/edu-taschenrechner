package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// Custom button widget that extends the default button with an alternative text.
type Button struct {
	widget.Button
	AlternativeText string
	display         fyne.CanvasObject
}

// Invert inverts the text of the button with the alternative text.
func (b *Button) Invert() {
	if b.AlternativeText == "" {
		return
	}

	b.AlternativeText, b.Text = b.Text, b.AlternativeText
	b.Refresh()
}

// GetOnTapped returns a function that sets the text of the display widget to the button's text.
func (b *Button) GetOnTapped() func() {
	return func() {
		if v, ok := b.display.(*Display); ok {
			v.SetText(b.Text)
		}
	}
}

// SetAlternateText sets the alternative text of the button.
func (b *Button) SetAlternateText(alt string) *Button {
	b.AlternativeText = alt
	return b
}

// SetOnTapped sets the function that is called when the button is tapped.
func (b *Button) SetOnTapped(fn func()) *Button {
	b.Button.OnTapped = fn
	return b
}

// Create a new button widget with the given text and tapped function.
func NewButton(text string, display fyne.CanvasObject) *Button {
	btn := &Button{
		Button:  widget.Button{Text: text},
		display: display,
	}
	btn.OnTapped = btn.GetOnTapped()
	btn.ExtendBaseWidget(btn)
	return btn
}
