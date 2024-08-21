package memory

import (
	"sync"
)

// Make sure *memoryCell implements MemoryCell
var _ MemoryCell = &memoryCell{}

// MemoryCell is an invisible CanvasObject for storing and retrieving anything.
type MemoryCell interface {
	Get() any
	Set(any) error
}

// *memoryCell will implement MemoryCell
type memoryCell struct {
	sync.Mutex
	value any
}

// Get retrieves the value from the memory cell.
func (m *memoryCell) Get() any {
	m.Lock()
	defer m.Unlock()

	return m.value
}

// Set stores the value in the memory cell.
func (m *memoryCell) Set(value any) error {
	m.Lock()
	defer m.Unlock()

	m.value = value

	return nil
}

// NewMemoryCell creates a new MemoryCell object.
func NewMemoryCell() MemoryCell { return &memoryCell{} }
