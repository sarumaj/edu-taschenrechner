package memory

import (
	"sync"

	"fyne.io/fyne/v2"
)

var _ MemoryCell = &memoryCell{}

// MemoryCell is an invisible CanvasObject for storing and retrieving anything.
type MemoryCell interface {
	fyne.CanvasObject
	Get() any
	Set(any)
}

// *memoryCell will implement MemoryCell
type memoryCell struct {
	sync.Mutex
	Value any
}

/*
Implement fyne.CanvasObject
*/

func (*memoryCell) MinSize() fyne.Size      { return fyne.Size{} }
func (*memoryCell) Move(fyne.Position)      {}
func (*memoryCell) Position() fyne.Position { return fyne.Position{} }
func (*memoryCell) Resize(fyne.Size)        {}
func (*memoryCell) Size() fyne.Size         { return fyne.Size{} }
func (*memoryCell) Hide()                   {}
func (*memoryCell) Visible() bool           { return false }
func (*memoryCell) Show()                   {}
func (*memoryCell) Refresh()                {}

/*
Implement MemoryCell
*/

func (m *memoryCell) Get() any {
	m.Lock()
	defer m.Unlock()

	return m.Value
}

func (m *memoryCell) Set(value any) {
	m.Lock()
	m.Value = value
	m.Unlock()
}

// NewMemoryCell creates a new MemoryCell object.
func NewMemoryCell() MemoryCell {
	return &memoryCell{}
}
