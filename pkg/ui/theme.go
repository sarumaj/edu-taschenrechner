//go:generate sh -c "$(go env GOPATH)/bin/fyne bundle -o bundled.go -package ui icons"
//go:generate sh -c "$(go env GOPATH)/bin/fyne bundle -a -o bundled.go -package ui fonts"
package ui

import (
	"fyne.io/fyne/v2"
)

// DoubleSizeTheme is a custom theme that doubles the size of the default theme.
type DoubleSizeTheme struct{ fyne.Theme }

// Font returns the Asana Math font.
func (t *DoubleSizeTheme) Font(fyne.TextStyle) fyne.Resource { return resourceAsanaMathOtf }

// Size returns the size of the theme multiplied by 2.
func (t *DoubleSizeTheme) Size(name fyne.ThemeSizeName) float32 {
	return t.Theme.Size(name) * 2
}

// NewDoubleSizeTheme creates a new theme that doubles the size of the default theme.
func NewDoubleSizeTheme(theme fyne.Theme) *DoubleSizeTheme { return &DoubleSizeTheme{Theme: theme} }
