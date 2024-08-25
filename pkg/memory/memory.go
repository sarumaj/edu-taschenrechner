/*
Package memory provides a simple memory cell for storing and retrieving big floating point numbers.
Generic memory cells can be created using the NewGenericMemoryCell function.
Memory cells can be used to store and retrieve variables for the calculator.

Example:

	cell := memory.NewMemoryCell()
	cell.Set(big.NewFloat(42))
	fmt.Println(cell.Get()) // prints 42
*/
package memory

import (
	"math/big"
	"sync"
)

// Make sure *memoryCell implements MemoryCell
var _ MemoryCell = (*memoryCell[*big.Float])(nil)

// MemoryCell is an interface for storing and retrieving big floating point numbers.
type MemoryCell = MemoryCellInterface[*big.Float]

// MemoryCellInterface is a generic interface for storing and retrieving anything.
type MemoryCellInterface[T any] interface {
	Get() T
	Set(T) error
}

// *memoryCell will implement MemoryCellInterface for any type T.
type memoryCell[T any] struct {
	sync.Mutex
	value T
}

// Get retrieves the value from the memory cell.
func (m *memoryCell[T]) Get() T {
	m.Lock()
	defer m.Unlock()

	return m.value
}

// Set stores the value in the memory cell.
func (m *memoryCell[T]) Set(value T) error {
	m.Lock()
	defer m.Unlock()

	m.value = value

	return nil
}

// NewMemoryCell creates a new MemoryCell object.
func NewMemoryCell() MemoryCell { return NewGenericMemoryCell[*big.Float]() }

// NewGenericMemoryCell creates a new MemoryCell object with a generic type.
func NewGenericMemoryCell[T any]() MemoryCellInterface[T] { return &memoryCell[T]{} }
