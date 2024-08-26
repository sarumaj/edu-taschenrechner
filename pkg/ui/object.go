//go:build !headless

package ui

import (
	"fyne.io/fyne/v2"
)

// ObjectStorage is a map of canvas objects.
type ObjectStorage map[string]fyne.CanvasObject

// SelectButtons selects the buttons from the object storage.
func (o ObjectStorage) SelectButtons(in ...string) (out []*Button) {
	return selectObjects[*Button](o, in...)
}

// SelectDisplay selects the display from the object storage.
func (o ObjectStorage) SelectDisplay(in string) (out *Display) {
	v, _ := o[in].(*Display)
	return v
}

// SelectDropDowns selects the dropdowns from the object storage.
func (o ObjectStorage) SelectDropDowns(in ...string) (out []*ButtonDropDown) {
	return selectObjects[*ButtonDropDown](o, in...)
}

// SelectIcons selects the icons from the object storage.
func (o ObjectStorage) SelectIcons(in ...string) (out []*Icon) {
	return selectObjects[*Icon](o, in...)
}

// SelectCanvasObjects selects the canvas objects from the object storage.
func (o ObjectStorage) SelectCanvasObjects(in ...string) (out []fyne.CanvasObject) {
	return selectObjects[fyne.CanvasObject](o, in...)
}

// selectObjects selects the objects from the map based on the given keys and type.
func selectObjects[O any](from map[string]fyne.CanvasObject, in ...string) (out []O) {
	for _, c := range in {
		e, ok := from[c]
		if !ok {
			continue
		}

		v, ok := e.(O)
		if !ok {
			continue
		}

		out = append(out, v)
	}

	return
}
