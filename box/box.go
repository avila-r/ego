package box

import (
	"fmt"
)

// Box is a generic container that may or may not hold a value
type Box[T any] struct {
	value   *T
	present bool
}

// Get returns the value or nil if empty
func (b Box[T]) Get() *T {
	return b.value
}

// GetOrDefault returns the value or a default if empty
func (b Box[T]) GetOrDefault(defaultValue T) T {
	if b.present {
		return *b.value
	}
	return defaultValue
}

// With sets a new value and returns the Box
func (b Box[T]) With(value T) Box[T] {
	b.value = &value
	b.present = true
	return b
}

// Peek executes the action if a value is present and returns the Box
func (b Box[T]) Peek(action func(T)) Box[T] {
	if b.IsPresent() {
		action(*b.value)
	}
	return b
}

// Then applies a mapper function and updates the Box value
func (b Box[T]) Then(mapper func(T) T) Box[T] {
	if b.present {
		newValue := mapper(*b.value)
		b.value = &newValue
	}
	return b
}

// ThenSupplier sets the value from a supplier function
func (b Box[T]) ThenSupplier(supplier func() T) Box[T] {
	newValue := supplier()
	b.value = &newValue
	b.present = true
	return b
}

// ThenConsumer executes the action if a value is present
func (b Box[T]) ThenConsumer(action func(T)) Box[T] {
	if b.IsPresent() {
		action(*b.value)
	}
	return b
}

// Filter returns this Box if the predicate matches, otherwise returns empty Box
func (b Box[T]) Filter(predicate func(T) bool) Box[T] {
	if b.IsPresent() && predicate(*b.value) {
		return b
	}
	return Empty[T]()
}

// Map transforms the value using the mapper function
func Map[T, R any](b Box[T], mapper func(T) R) Box[R] {
	if b.IsPresent() {
		return Of(mapper(*b.value))
	}
	return Empty[R]()
}

// FlatMap transforms the value and flattens the result
func FlatMap[T, R any](b Box[T], mapper func(T) Box[R]) Box[R] {
	if b.IsPresent() {
		return mapper(*b.value)
	}
	return Empty[R]()
}

// Deflated clears the value and returns the Box
func (b Box[T]) Deflated() Box[T] {
	b.value = nil
	b.present = false
	return b
}

// Copy creates a shallow copy of the Box
func (b Box[T]) Copy() Box[T] {
	if b.IsPresent() {
		return Of(*b.value)
	}
	return Empty[T]()
}

// IsPresent returns true if the Box contains a value
func (b Box[T]) IsPresent() bool {
	return b.present
}

// IsEmpty returns true if the Box is empty
func (b Box[T]) IsEmpty() bool {
	return !b.present
}

// Set updates the Box value
func (b *Box[T]) Set(value T) {
	b.value = &value
	b.present = true
}

// Deflate clears the Box value
func (b *Box[T]) Deflate() {
	b.value = nil
	b.present = false
}

// IfPresent executes the action if a value is present
func (b Box[T]) IfPresent(action func(T)) {
	if b.IsPresent() {
		action(*b.value)
	}
}

// String returns a string representation of the Box
func (b Box[T]) String() string {
	if b.IsEmpty() {
		return "Box(<empty>)"
	}
	return fmt.Sprintf("Box(%v)", *b.value)
}

// Equals checks if two Boxes are equal
func (b Box[T]) Equals(other *Box[T]) bool {
	if other == nil {
		return false
	}

	if b.present != other.present {
		return false
	}

	if !b.present {
		return true
	}

	if b == *other {
		return true
	}

	// Note: This uses == which works for comparable types
	// For non-comparable types, you'd need to pass a custom comparator
	return fmt.Sprintf("%v", *b.value) == fmt.Sprintf("%v", *other.value)
}
