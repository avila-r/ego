package stream

import (
	"github.com/avila-r/ego/function"
	"github.com/avila-r/ego/optional"
	"sort"
)

type Stream[T comparable] struct {
	elements []T
}

func Of[T comparable](elements ...T) Stream[T] {
	return Stream[T]{elements: elements}
}

func From[T comparable](collectable Collectable[T]) Stream[T] {
	return Stream[T]{elements: collectable.Elements()}
}

func Empty[T comparable]() Stream[T] {
	return Stream[T]{elements: []T{}}
}

func OfNullable[T comparable](value *T) Stream[T] {
	if value == nil {
		return Empty[T]()
	}
	return Of(*value)
}

func (s Stream[T]) Filter(predicate function.Predicate[T]) Stream[T] {
	newElements := make([]T, 0)
	for _, element := range s.elements {
		if predicate.Test(element) {
			newElements = append(newElements, element)
		}
	}
	return Stream[T]{elements: newElements}
}

func Map[T comparable, R comparable](s Stream[T], mapper function.Function[T, R]) Stream[R] {
	newElements := make([]R, 0, len(s.elements))
	for _, element := range s.elements {
		mappedElement := mapper.Apply(element)
		newElements = append(newElements, mappedElement)
	}
	return Stream[R]{elements: newElements}
}

func FlatMap[T comparable, R comparable](s Stream[T], mapper function.Function[T, Stream[R]]) Stream[R] {
	newElements := make([]R, 0)
	for _, element := range s.elements {
		mappedStream := mapper.Apply(element)
		newElements = append(newElements, mappedStream.elements...)
	}
	return Stream[R]{elements: newElements}
}

func MapMulti[T comparable, R comparable](s Stream[T], mapper function.BiConsumer[T, function.Consumer[R]]) Stream[R] {
	return FlatMap(s, function.NewFunction(func(t T) Stream[R] {
		newElements := make([]R, 0)
		consumer := function.NewConsumer[R](func(r R) {
			newElements = append(newElements, r)
		})
		mapper.Accept(t, consumer)
		return Stream[R]{elements: newElements}
	}))
}

func (s Stream[T]) Distinct() Stream[T] {
	seen := make(map[T]bool)
	newElements := make([]T, 0)
	for _, element := range s.elements {
		if _, found := seen[element]; !found {
			seen[element] = true
			newElements = append(newElements, element)
		}
	}
	return Stream[T]{elements: newElements}
}

func (s Stream[T]) Sort() Stream[T] {
	compare := function.DefaultComparator[T]{}

	for i := 0; i < len(s.elements)-1; i++ {
		for j := i + 1; j < len(s.elements); j++ {
			if compare.Compare(s.elements[i], s.elements[j]) > 0 {
				s.elements[i], s.elements[j] = s.elements[j], s.elements[i]
			}
		}
	}

	return Stream[T]{elements: s.elements}
}

func (s Stream[T]) Sorted(comparator function.Comparator[T]) Stream[T] {
	sort.Slice(s.elements, func(i, j int) bool {
		return comparator.Compare(s.elements[i], s.elements[j]) < 0
	})

	return Stream[T]{elements: s.elements}
}

func (s Stream[T]) Peek(consumer function.Consumer[T]) Stream[T] {
	for _, element := range s.elements {
		consumer.Accept(element)
	}
	return s
}

func (s Stream[T]) Limit(number int) Stream[T] {
	if number >= len(s.elements) {
		return s
	}
	return Stream[T]{elements: s.elements[:number]}
}

func (s Stream[T]) Skip(quantity int) Stream[T] {
	if quantity >= len(s.elements) {
		return Stream[T]{elements: []T{}}
	}
	return Stream[T]{elements: s.elements[quantity:]}
}

func (s Stream[T]) TakeWhile(predicate function.Predicate[T]) Stream[T] {
	newElements := make([]T, 0)
	for _, element := range s.elements {
		if !predicate.Test(element) {
			break
		}
		newElements = append(newElements, element)
	}
	return Stream[T]{elements: newElements}
}

func (s Stream[T]) DropWhile(predicate function.Predicate[T]) Stream[T] {
	newElements := make([]T, 0)
	dropping := true
	for _, element := range s.elements {
		if dropping && predicate.Test(element) {
			continue
		}
		dropping = false
		newElements = append(newElements, element)
	}
	return Stream[T]{elements: newElements}
}

func (s Stream[T]) ForEach(consumer function.Consumer[T]) {
	for _, element := range s.elements {
		consumer.Accept(element)
	}
}

func (s Stream[T]) ForEachOrdered(consumer function.Consumer[T]) {
	s.Sort().ForEach(consumer)
}

func (s Stream[T]) toSlice() []T {
	return s.elements
}
func (s Stream[T]) ReduceWithIdentity(identity T, accumulator function.BinaryOperator[T]) T {
	result := identity
	for _, element := range s.elements {
		result = accumulator.Apply(result, element)
	}
	return result
}

func (s Stream[T]) Reduce(accumulator function.BinaryOperator[T]) optional.Optional[T] {
	if len(s.elements) == 0 {
		var zero optional.Optional[T]
		return zero
	}
	result := s.elements[0]
	for _, element := range s.elements[1:] {
		result = accumulator.Apply(result, element)
	}
	return optional.Of(result)
}

func Reduce[U comparable](s Stream[U], identity U, accumulator function.BinaryOperator[U]) U {
	result := identity
	for _, element := range s.elements {
		result = accumulator.Apply(result, element)
	}
	return result
}

func Collect[T, A, R any](stream Stream[any], collector Collector[T, A, R]) R {
	acc := collector.Supplier().Get()

	accumulator := collector.Accumulator()
	for _, e := range stream.elements {
		accumulator.Accept(acc, e)
	}

	return collector.Finisher().Apply(acc)
}

func (s Stream[T]) Min(comparator function.Comparator[T]) optional.Optional[T] {
	return s.findMinimumOrMaximum(true, comparator)
}

func (s Stream[T]) Max(comparator function.Comparator[T]) optional.Optional[T] {
	return s.findMinimumOrMaximum(false, comparator)
}

func (s Stream[T]) findMinimumOrMaximum(isMin bool, comparator function.Comparator[T]) optional.Optional[T] {
	if len(s.elements) == 0 {
		var empty optional.Optional[T]
		return empty
	}
	m := s.elements[0]
	for _, e := range s.elements[1:] {
		compareResult := comparator.Compare(e, m)
		if (isMin && compareResult < 0) || (!isMin && compareResult > 0) {
			m = e
		}
	}
	return optional.Of(m)
}

func (s Stream[T]) Count() int64 {
	return int64(len(s.elements))
}

func (s Stream[T]) AnyMatch(predicate function.Predicate[T]) bool {
	for _, e := range s.elements {
		if predicate.Test(e) {
			return true
		}
	}
	return false
}

func (s Stream[T]) AllMatch(predicate function.Predicate[T]) bool {
	for _, e := range s.elements {
		if !predicate.Test(e) {
			return false
		}
	}
	return true
}

func (s Stream[T]) NoneMatch(predicate function.Predicate[T]) bool {
	for _, e := range s.elements {
		if predicate.Test(e) {
			return false
		}
	}
	return true
}

func (s Stream[T]) FindFirst() optional.Optional[T] {
	if len(s.elements) == 0 {
		var empty optional.Optional[T]
		return empty
	}
	return optional.Of(s.elements[0])
}

func (s Stream[T]) FindAny() optional.Optional[T] {
	return s.FindFirst()
}

func Concat[T comparable](a, b Stream[T]) Stream[T] {
	combined := make([]T, 0, len(a.elements)+len(b.elements))
	combined = append(combined, a.elements...)
	combined = append(combined, b.elements...)
	return Stream[T]{elements: combined}
}

func Generate[T comparable](supplier function.Supplier[T], limit int) Stream[T] {
	elements := make([]T, 0, limit)
	for i := 0; i < limit; i++ {
		elements = append(elements, supplier.Get())
	}
	return Stream[T]{elements: elements}
}

func Iterate[T comparable](seed T, f function.UnaryOperator[T], limit int) Stream[T] {
	elements := make([]T, 0, limit)
	current := seed
	for i := 0; i < limit; i++ {
		elements = append(elements, current)
		current = f.Apply(current)
	}
	return Stream[T]{elements: elements}
}

func IterateWhile[T comparable](
	seed T,
	hasNext function.Predicate[T],
	next function.UnaryOperator[T],
) Stream[T] {
	elements := make([]T, 0)
	current := seed
	for hasNext.Test(current) {
		elements = append(elements, current)
		current = next.Apply(current)
	}
	return Stream[T]{elements: elements}
}
