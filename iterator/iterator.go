package iterator

import "github.com/avila-r/ego/stream"

type Iterable[T any] interface {
	ForEach(func(T))
	Iterator() Iterator[T]
}

// Iterator represents a stateful iterator over elements
type Iterator[T any] interface {

	// HasNext returns true if there are more elements to iterate
	HasNext() bool

	// Next returns the next element and advances the iterator
	// Panics if HasNext() is false
	Next() T

	// Peek returns the next element without advancing the iterator
	// Panics if HasNext() is false
	Peek() T

	// Reset resets the iterator to the beginning
	Reset()

	Remaining() int

	Collect() []T

	ForEach(action func(T))

	Filter(predicate func(T) bool) Iterator[T]
}

// SliceIterator is a default iterator implementation for slices
type SliceIterator[T any] struct {
	elements []T
	index    int
}

// Ensure SliceIterator implements Iterator
var _ Iterator[int] = (*SliceIterator[int])(nil)

func Of[T any](elements ...T) Iterator[T] {
	return &SliceIterator[T]{
		elements: elements,
		index:    0,
	}
}

func From[T any](collectable stream.Collectable[T]) Iterator[T] {
	return &SliceIterator[T]{
		elements: collectable.Elements(),
		index:    0,
	}
}

// HasNext returns true if there are more elements to iterate
func (it *SliceIterator[T]) HasNext() bool {
	return it.index < len(it.elements)
}

// Next returns the next element and advances the iterator
func (it *SliceIterator[T]) Next() T {
	if !it.HasNext() {
		panic("iterator: no more elements")
	}
	element := it.elements[it.index]
	it.index++
	return element
}

// Peek returns the next element without advancing the iterator
func (it *SliceIterator[T]) Peek() T {
	if !it.HasNext() {
		panic("iterator: no more elements")
	}
	return it.elements[it.index]
}

// Reset resets the iterator to the beginning
func (it *SliceIterator[T]) Reset() {
	it.index = 0
}

// Remaining returns the number of elements left to iterate
func (it *SliceIterator[T]) Remaining() int {
	return len(it.elements) - it.index
}

// Collect collects all remaining elements into a slice
func (it *SliceIterator[T]) Collect() []T {
	remaining := it.elements[it.index:]
	it.index = len(it.elements)
	return remaining
}

// ForEach applies a function to all remaining elements
func (it *SliceIterator[T]) ForEach(action func(T)) {
	for it.HasNext() {
		action(it.Next())
	}
}

// Filter returns a new iterator with only elements matching the predicate
func (it *SliceIterator[T]) Filter(predicate func(T) bool) Iterator[T] {
	filtered := make([]T, 0)
	for it.HasNext() {
		elem := it.Next()
		if predicate(elem) {
			filtered = append(filtered, elem)
		}
	}
	return Of(filtered...)
}

// Map transforms elements to a new type
func Map[T, U comparable](it Iterator[T], mapper func(T) U) Iterator[U] {
	mapped := make([]U, 0)
	for it.HasNext() {
		mapped = append(mapped, mapper(it.Next()))
	}
	return Of(mapped...)
}
