package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/widget"
)

var _ interface {
	desktop.Cursorable
	mobile.Touchable
	fyne.Tappable
} = (*ButtonDropDown)(nil)

// ButtonDropDown is a custom dropdown widget that extends the default select widget with buttons.
// The dropdown is used to select a button from a list of buttons.
type ButtonDropDown struct {
	widget.Select
	buttons []*Button
}

// Cursor returns the pointer cursor.
func (*ButtonDropDown) Cursor() desktop.Cursor { return desktop.PointerCursor }

// GetOnChanged returns a function that changes the selected index of the dropdown.
func (b *ButtonDropDown) GetOnChanged() func(string) {
	return func(s string) {
		defer b.FocusLost() // close the dropdown and lose focus

		for _, btn := range b.buttons {
			if btn.Text != s || btn.OnTapped == nil {
				continue
			}

			btn.OnTapped()
			break
		}
	}
}

// Update updates the dropdown widget with the given buttons.
func (b *ButtonDropDown) Update() {
	options := make([]string, 0)
	for _, btn := range b.buttons {
		options = append(options, btn.Text)
		btn.Refresh()
	}

	if len(options) == 0 {
		options = []string{""}
	}

	selected := b.SelectedIndex()
	if selected < 0 || selected >= len(options) {
		selected = 0
	}

	b.Options = options
	b.PlaceHolder = options[selected]

	// temporarily disable the onChanged function to prevent it from being called
	var onChanged func(string)
	onChanged, b.Select.OnChanged = b.Select.OnChanged, onChanged

	// set the selected index and restore the onChanged function
	b.Select.SetSelected(options[selected])
	b.SetOnChanged(onChanged)

	b.Refresh()
}

// SetOnChanged sets the function that changes the selected index of the dropdown.
func (b *ButtonDropDown) SetOnChanged(fn func(string)) *ButtonDropDown {
	b.Select.OnChanged = fn
	return b
}

// Tapped selects the button that was tapped and loses focus.
func (b *ButtonDropDown) Tapped(e *fyne.PointEvent) {
	defer b.FocusLost() // close the dropdown and lose focus
	b.Select.Tapped(e)
}

// TouchCancel cancels the touch event of the dropdown.
func (*ButtonDropDown) TouchCancel(*mobile.TouchEvent) {}

// TouchDown triggers the touch event of the dropdown.
func (b *ButtonDropDown) TouchDown(e *mobile.TouchEvent) {
	scale := fyne.CurrentApp().Driver().CanvasForObject(b).Scale()
	b.Tapped(&fyne.PointEvent{
		AbsolutePosition: fyne.NewPos(e.AbsolutePosition.X/scale, e.AbsolutePosition.Y/scale),
		Position:         fyne.NewPos(e.Position.X/scale, e.Position.Y/scale),
	})
}

// TouchUp cancels the touch event of the dropdown.
func (*ButtonDropDown) TouchUp(*mobile.TouchEvent) {}

// NewButtonDropDown creates a new dropdown widget with the given buttons.
func NewButtonDropDown(buttons []*Button) *ButtonDropDown {
	btn := &ButtonDropDown{
		Select:  widget.Select{Options: []string{}},
		buttons: buttons,
	}
	btn.SetOnChanged(btn.GetOnChanged())
	btn.Update()
	btn.ExtendBaseWidget(btn)
	return btn
}
