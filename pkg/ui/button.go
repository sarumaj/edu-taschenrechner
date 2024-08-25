package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Make sure the Button widget implements all necessary interfaces.
var _ interface {
	mobile.Touchable
	desktop.Cursorable
} = (*Button)(nil)

// Custom button widget that extends the default button with an alternative text.
type Button struct {
	widget.Button
	AlternativeText string
}

// CreateRenderer creates the renderer for the custom button, reusing the original Button renderer.
func (b *Button) CreateRenderer() fyne.WidgetRenderer {
	label := canvas.NewText(b.AlternativeText, color.Gray{Y: 100})
	label.TextSize = theme.TextSize() * 0.7 // 55% of the default text size
	label.Alignment = fyne.TextAlignCenter

	return &ButtonRenderer{
		WidgetRenderer: b.Button.CreateRenderer(),
		button:         b,
		altLabel:       label,
	}
}

// Cursor returns the pointer cursor.
func (*Button) Cursor() desktop.Cursor { return desktop.PointerCursor }

// Invert inverts the text of the button with the alternative text.
func (b *Button) Invert() *Button {
	if b.AlternativeText == "" {
		return b
	}

	b.AlternativeText, b.Text = b.Text, b.AlternativeText
	b.Refresh()
	return b
}

// GetOnTapped returns a function that sets the text of the display widget to the button's text.
func (b *Button) GetOnTapped(display *Display) func() { return func() { display.SetText(b.Text) } }

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

// TouchCancel cancels the touch event of the icon.
func (b *Button) TouchCancel(*mobile.TouchEvent) {}

// TouchDown triggers the touch event of the icon.
func (b *Button) TouchDown(e *mobile.TouchEvent) {
	scale := fyne.CurrentApp().Driver().CanvasForObject(b).Scale()
	b.Button.Tapped(&fyne.PointEvent{Position: fyne.NewPos(e.Position.X/scale, e.Position.Y/scale)})
}

// TouchUp cancels the touch event of the icon.
func (b *Button) TouchUp(*mobile.TouchEvent) {}

// ButtonRenderer extends the original Button renderer to include the alternative text.
type ButtonRenderer struct {
	fyne.WidgetRenderer
	button   *Button
	altLabel *canvas.Text
}

// Layout positions the components of the button, including the alternative text.
func (r *ButtonRenderer) Layout(size fyne.Size) {
	altLabelSize := r.altLabel.MinSize()
	padding := theme.Padding()

	// Position the alternative text in the bottom left corner
	r.altLabel.Move(fyne.NewPos(padding, size.Height-altLabelSize.Height-padding))

	r.altLabel.Resize(altLabelSize)
	r.WidgetRenderer.Layout(size)
}

// Refresh updates the state of the button and redraws it, including the alternative text.
func (r *ButtonRenderer) Refresh() {
	r.altLabel.Text = r.button.AlternativeText
	r.altLabel.Color = color.Gray{Y: 100}

	canvas.Refresh(r.button)
	r.WidgetRenderer.Refresh()
}

// Objects returns the objects that should be rendered for this button.
func (r *ButtonRenderer) Objects() []fyne.CanvasObject {
	return append(r.WidgetRenderer.Objects(), r.altLabel)
}

// Create a new button widget with the given text and tapped function.
func NewButton(text string, display *Display) *Button {
	return NewButtonWithIcon(text, nil, display)
}

// Create a new button widget with the given text, icon, and tapped function.
func NewButtonWithIcon(text string, icon fyne.Resource, display *Display) *Button {
	btn := &Button{
		Button: widget.Button{Text: text, Icon: icon},
	}
	btn.SetOnTapped(btn.GetOnTapped(display))
	btn.ExtendBaseWidget(btn)
	return btn
}
