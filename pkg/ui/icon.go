//go:build !headless

package ui

import (
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Make sure the Icon widget implements all necessary interfaces.
var _ interface {
	widget.ToolbarItem
	mobile.Touchable
	desktop.Cursorable
} = (*Icon)(nil)

// Icon is a custom icon widget that extends the default icon with a hyperlink.
type Icon struct {
	link          widget.Hyperlink
	lightResource fyne.Resource
	darkResource  fyne.Resource
	widget.Button
}

// Cursor returns the pointer cursor.
func (i *Icon) Cursor() desktop.Cursor { return desktop.PointerCursor }

// GetOnTapped returns a function that opens the hyperlink of the icon.
func (i *Icon) GetOnTapped() func() { return func() { i.link.Tapped(&fyne.PointEvent{}) } }

// SetOnTapped sets the OnTapped function of the icon.
func (i *Icon) SetOnTapped(fn func()) *Icon { i.Button.OnTapped = fn; return i }

// ToolbarObject returns the icon widget as a toolbar object.
func (i *Icon) ToolbarObject() fyne.CanvasObject { return i }

// Update updates the icon widget with the given theme variant.
func (i *Icon) Update() *Icon {
	if i.darkResource == nil {
		return i
	}

	variant := fyne.CurrentApp().Settings().ThemeVariant()
	i.Icon = map[fyne.ThemeVariant]fyne.Resource{
		theme.VariantLight: i.lightResource,
		theme.VariantDark:  i.darkResource,
	}[variant]

	return i
}

// TouchCancel cancels the touch event of the icon.
func (*Icon) TouchCancel(*mobile.TouchEvent) {}

// TouchDown triggers the touch event of the icon.
func (i *Icon) TouchDown(e *mobile.TouchEvent) {
	scale := fyne.CurrentApp().Driver().CanvasForObject(i).Scale()
	i.Tapped(&fyne.PointEvent{
		AbsolutePosition: fyne.NewPos(e.AbsolutePosition.X/scale, e.AbsolutePosition.Y/scale),
		Position:         fyne.NewPos(e.Position.X/scale, e.Position.Y/scale),
	})
}

// TouchUp cancels the touch event of the icon.
func (*Icon) TouchUp(*mobile.TouchEvent) {}

// NewIcon creates a new icon with a hyperlink.
func NewIcon(link string, icon, darkIcon fyne.Resource) *Icon {
	i := &Icon{
		Button:        widget.Button{Icon: icon},
		lightResource: icon,
		darkResource:  darkIcon,
	}
	parsed, _ := url.Parse(link)
	i.link = widget.Hyperlink{URL: parsed}
	i.SetOnTapped(i.GetOnTapped())

	i.ExtendBaseWidget(i)
	return i
}
