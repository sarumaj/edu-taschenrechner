package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// ButtonDropDown is a custom dropdown widget that extends the default select widget with buttons.
// The dropdown is used to select a button from a list of buttons.
type ButtonDropDown struct {
	widget.Select
	buttons []fyne.CanvasObject
}

// GetOnChanged returns a function that changes the selected index of the dropdown.
func (b *ButtonDropDown) GetOnChanged() func(string) {
	return func(s string) {
		defer b.FocusLost() // close the dropdown and lose focus

		for _, btn := range b.buttons {
			if btn, ok := btn.(*Button); ok && btn.Text == s && btn.OnTapped != nil {
				btn.OnTapped()
				break
			}
		}
	}
}

// Update updates the dropdown widget with the given buttons.
func (b *ButtonDropDown) Update() {
	options := make([]string, 0)
	for _, btn := range b.buttons {
		if button, ok := btn.(*Button); ok {
			options = append(options, button.Text)
			button.Refresh()
		}
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

// NewButtonDropDown creates a new dropdown widget with the given buttons.
func NewButtonDropDown(buttons []fyne.CanvasObject) *ButtonDropDown {
	btn := &ButtonDropDown{
		Select:  widget.Select{Options: []string{}},
		buttons: buttons,
	}
	btn.OnChanged = btn.GetOnChanged()
	btn.Update()
	btn.ExtendBaseWidget(btn)
	return btn
}
