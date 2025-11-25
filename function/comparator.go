package function

import "github.com/avila-r/ego/constraint"

type Comparator[T any] interface {
	Compare(a, b T) int
}

type ComparableMethod[T any] interface {
	Compare(T) int
}

type DefaultComparator[T any] struct {
	comparator func(a, b T) int
}

func (c *DefaultComparator[T]) Compare(a T, b T) int {
	return c.comparator(a, b)
}

func (c *DefaultComparator[T]) Equals(other Comparator[T]) bool {
	return c == other
}

func NewComparator[T any](comparator func(a, b T) int) Comparator[T] {
	return &DefaultComparator[T]{comparator: comparator}
}

func (c *DefaultComparator[T]) Reversed() Comparator[T] {
	return NewComparator(func(c1, c2 T) int {
		return -c.Compare(c1, c2)
	})
}

func (c *DefaultComparator[T]) ThenComparing(other Comparator[T]) Comparator[T] {
	return NewComparator(func(c1, c2 T) int {
		res := c.Compare(c1, c2)
		if res != 0 {
			return res
		}
		return other.Compare(c1, c2)
	})
}

func (c *DefaultComparator[T]) ThenComparingInt(keyExtractor ToIntFunction[T]) Comparator[T] {
	return c.ThenComparing(ComparingInt(keyExtractor))
}

func (c *DefaultComparator[T]) ThenComparingFloat64(keyExtractor ToFloatFunction[T]) Comparator[T] {
	return c.ThenComparing(ComparingFloat64(keyExtractor))
}

func NaturalOrder[T ComparableMethod[T]]() Comparator[T] {
	return NewComparator(func(c1, c2 T) int {
		return c1.Compare(c2)
	})
}

func ReverseOrder[T ComparableMethod[T]]() Comparator[T] {
	return NaturalOrder[T]().(*DefaultComparator[T]).Reversed()
}

func Comparing[T, U any](keyExtractor Function[T, U], keyComparator Comparator[U]) Comparator[T] {
	return NewComparator(func(c1, c2 T) int {
		key1 := keyExtractor.Apply(c1)
		key2 := keyExtractor.Apply(c2)
		return keyComparator.Compare(key1, key2)
	})
}

func ComparingNatural[T, U ComparableMethod[U]](keyExtractor Function[T, U]) Comparator[T] {
	keyComparator := NaturalOrder[U]()
	return Comparing(keyExtractor, keyComparator)
}

func ComparingInt[T any](keyExtractor ToIntFunction[T]) Comparator[T] {
	return NewComparator(func(a, b T) int {
		return orderedCompare(keyExtractor.ApplyAsInt(a), keyExtractor.ApplyAsInt(b))
	})
}

func ComparingFloat64[T any](keyExtractor ToFloatFunction[T]) Comparator[T] {
	return NewComparator(func(a, b T) int {
		return orderedCompare(keyExtractor.ApplyAsFloat(a), keyExtractor.ApplyAsFloat(b))
	})
}

func orderedCompare[T constraint.Ordered](a, b T) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}
