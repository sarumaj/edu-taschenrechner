package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Make sure the ToolbarItem widget implements all necessary interfaces.
var _ interface {
	widget.ToolbarItem
	mobile.Touchable
	desktop.Cursorable
} = (*ToolbarItem)(nil)

// Toolbar is a custom toolbar widget that extends the default toolbar with a display object.
type Toolbar struct{ widget.Toolbar }

// SetActions sets the actions of the toolbar.
func (t *Toolbar) SetActions(actions ...widget.ToolbarItem) {
	t.Toolbar.Items = make([]widget.ToolbarItem, len(actions))
	_ = copy(t.Toolbar.Items, actions)
}

// ToolbarItem is a custom toolbar item widget that extends the default button with the toolbar object.
type ToolbarItem struct{ widget.Button }

// Cursor returns the pointer cursor.
func (*ToolbarItem) Cursor() desktop.Cursor { return desktop.PointerCursor }

// SetOnTapped sets the OnTapped function of the toolbar item.
func (t *ToolbarItem) SetOnTapped(fn func()) *ToolbarItem { t.Button.OnTapped = fn; return t }

// ToolbarObject returns the toolbar item as a toolbar object.
func (t *ToolbarItem) ToolbarObject() fyne.CanvasObject { return t }

// TouchCancel cancels the touch event of the toolbar item.
func (*ToolbarItem) TouchCancel(*mobile.TouchEvent) {}

// TouchDown triggers the touch event of the toolbar item.
func (t *ToolbarItem) TouchDown(e *mobile.TouchEvent) {
	scale := fyne.CurrentApp().Driver().CanvasForObject(t).Scale()
	t.Button.Tapped(&fyne.PointEvent{Position: fyne.NewPos(e.Position.X/scale, e.Position.Y/scale)})
}

// TouchUp cancels the touch event of the toolbar item.
func (*ToolbarItem) TouchUp(*mobile.TouchEvent) {}

// NewDisplayToolbar creates a new toolbar with the given actions and display object.
func NewDisplayToolbar(display fyne.CanvasObject, actions ...widget.ToolbarItem) *Toolbar {
	if display, ok := display.(*Display); ok {
		if display.Entry.ActionItem == nil {
			actions = append([]widget.ToolbarItem{
				NewToolbarItem(theme.ContentCopyIcon()).SetOnTapped(display.CopyToClipboard),
				NewToolbarItem(theme.SettingsIcon()).SetOnTapped(display.MeasureDisplayCapacity),
			}, actions...)
		} else {
			actions = append([]widget.ToolbarItem{
				NewToolbarItem(theme.ContentCopyIcon()).SetOnTapped(display.CopyToClipboard),
			}, actions...)
		}
	}

	return NewToolbar(actions...)
}

// NewToolbar creates a new toolbar with the given actions.
func NewToolbar(actions ...widget.ToolbarItem) *Toolbar {
	t := &Toolbar{}
	t.SetActions(actions...)
	t.ExtendBaseWidget(t)
	return t
}

// NewToolbarItem creates a new toolbar item with the given icon.
func NewToolbarItem(icon fyne.Resource) *ToolbarItem {
	i := &ToolbarItem{Button: widget.Button{Icon: icon}}
	i.ExtendBaseWidget(i)
	return i
}
