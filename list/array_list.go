package list

import (
	"slices"

	"github.com/avila-r/ego/collection"
	"github.com/avila-r/ego/iterator"
	"github.com/avila-r/ego/stream"
)

type ArrayList[T comparable] struct {
	elements []T
}

var _ collection.List[int] = (*ArrayList[int])(nil)

func NewArrayList[T comparable](items ...T) *ArrayList[T] {
	return &ArrayList[T]{
		elements: items,
	}
}

func EmptyArrayList[T comparable]() *ArrayList[T] {
	return &ArrayList[T]{
		elements: []T{},
	}
}

func (l *ArrayList[T]) Add(items ...T) {
	l.elements = append(l.elements, items...)
}

func (l *ArrayList[T]) Get(index int) (T, bool) {
	var zero T
	if index < 0 || index >= len(l.elements) {
		return zero, false
	}
	return l.elements[index], true
}

func (l *ArrayList[T]) Set(index int, value T) bool {
	if index < 0 || index >= len(l.elements) {
		return false
	}
	l.elements[index] = value
	return true
}

func (l *ArrayList[T]) Remove(index int) bool {
	if index < 0 || index >= len(l.elements) {
		return false
	}
	l.elements = append(l.elements[:index], l.elements[index+1:]...)
	return true
}

func (l *ArrayList[T]) Contains(value T) bool {
	return slices.Contains(l.elements, value)
}

func (l *ArrayList[T]) Size() int {
	return len(l.elements)
}

func (l *ArrayList[T]) IsEmpty() bool {
	return len(l.elements) == 0
}

func (l *ArrayList[T]) Clear() {
	l.elements = []T{}
}

func (l *ArrayList[T]) Items() []T {
	return l.elements
}

func (l *ArrayList[T]) Elements() []T {
	return l.elements
}

func (l *ArrayList[T]) Stream() stream.Stream[T] {
	return stream.From(l)
}

func (l *ArrayList[T]) ForEach(action func(T)) {
	for _, item := range l.elements {
		action(item)
	}
}

func (l *ArrayList[T]) Iterator() iterator.Iterator[T] {
	return iterator.From(l)
}
