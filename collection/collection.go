package collection

import (
	"github.com/avila-r/ego/iterator"
	"github.com/avila-r/ego/slice"
)

type Collection[T any] interface {
	Add(elements ...T)
	Get(index int) (T, bool)
	Remove(index int) bool
	Size() int
	IsEmpty() bool
	Clear()
	Clone() Collection[T]
	Elements() []T

	iterator.Iterable[T]
}

// DefaultCollection is a simple slice-backed implementation of Collection
type DefaultCollection[T any] struct {
	elements []T
}

// Ensure DefaultCollection implements Collection
var _ Collection[int] = (*DefaultCollection[int])(nil)

// Of creates a new collection with the provided elements
func Of[T any](items ...T) *DefaultCollection[T] {
	return &DefaultCollection[T]{
		elements: items,
	}
}

// Empty creates a new empty collection
func Empty[T any]() *DefaultCollection[T] {
	return &DefaultCollection[T]{
		elements: []T{},
	}
}

// Add appends elements to the collection
func (c *DefaultCollection[T]) Add(elements ...T) {
	c.elements = append(c.elements, elements...)
}

// Get retrieves the element at the specified index
func (c *DefaultCollection[T]) Get(index int) (T, bool) {
	var zero T
	if index < 0 || index >= len(c.elements) {
		return zero, false
	}
	return c.elements[index], true
}

// Remove removes the element at the specified index
func (c *DefaultCollection[T]) Remove(index int) bool {
	if index < 0 || index >= len(c.elements) {
		return false
	}
	c.elements = append(c.elements[:index], c.elements[index+1:]...)
	return true
}

// Size returns the number of elements in the collection
func (c *DefaultCollection[T]) Size() int {
	return len(c.elements)
}

// IsEmpty returns true if the collection has no elements
func (c *DefaultCollection[T]) IsEmpty() bool {
	return len(c.elements) == 0
}

// Clear removes all elements from the collection
func (c *DefaultCollection[T]) Clear() {
	c.elements = []T{}
}

// Elements returns a slice of all elements (implements stream.Collectable)
func (c *DefaultCollection[T]) Elements() []T {
	return slice.Clone(c.elements)
}

// ForEach applies the given action to each element
func (c *DefaultCollection[T]) ForEach(action func(T)) {
	for _, item := range c.elements {
		action(item)
	}
}

// Clone creates a shallow copy of the collection
func (c *DefaultCollection[T]) Clone() Collection[T] {
	cloned := make([]T, len(c.elements))
	copy(cloned, c.elements)
	return &DefaultCollection[T]{
		elements: cloned,
	}
}

// Iterator returns an iterator over the elements in the collection
func (c *DefaultCollection[T]) Iterator() iterator.Iterator[T] {
	return iterator.Of(c.elements...)
}
