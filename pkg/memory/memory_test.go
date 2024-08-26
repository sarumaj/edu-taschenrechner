package memory

import "testing"

func TestExampleFor_MemoryCell(t *testing.T) {
	// create a new memory cell
	cell := NewGenericMemoryCell[int]()

	// store a value in the memory cell
	if err := cell.Set(42); err != nil {
		t.Errorf("MemoryCell.Set() error = %v, want nil", err)
	}

	// retrieve the value from the memory cell
	if got := cell.Get(); got != 42 {
		t.Errorf("MemoryCell.Get() = %v, want %v", got, 42)
	}
}
