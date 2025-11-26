package slice

import (
	"slices"

	"github.com/avila-r/ego/constraint"
	"github.com/avila-r/ego/optional"
	"github.com/avila-r/ego/stream"
)

func Of[T any](values ...T) []T {
	return values
}

func New[T any](size ...int) []T {
	if len(size) > 0 {
		return make([]T, 0, size[0])
	}

	return []T{}
}

func Empty[T any](size ...int) []T {
	return New[T](size...)
}

func Sized[T any](size int) []T {
	return New[T](size)
}

func Append[T any](at []T, t ...T) []T {
	return append(at, t...)
}

func Add[T any](at *[]T, t ...T) {
	*at = append(*at, t...)
}

func Filter[S ~[]E, E any](values S, filter func(v E) bool) S {
	result := S{}

	for _, v := range values {
		if filter(v) {
			result = append(result, v)
		}
	}

	return result
}

func IsEmpty[T any](t []T) bool {
	return len(t) == 0
}

func First[T any](t []T) optional.Optional[T] {
	return optional.Of(t[0])
}

func Last[T any](t []T) optional.Optional[T] {
	if IsEmpty(t) {
		return optional.Empty[T]()
	}

	v := t[len(t)-1]

	return optional.Of(v)
}

func Size[T any](t []T) int {
	return len(t)
}

func IsNil[T any](t []T) bool {
	return t == nil
}

func Stream[T comparable](t []T) stream.Stream[T] {
	return stream.Of(t...)
}

func ForEach[T any](s []T, f func(T)) {
	for _, v := range s {
		f(v)
	}
}

func Map[T any, R any](s []T, mapper func(T) R) []R {
	result := make([]R, len(s))
	for i, v := range s {
		result[i] = mapper(v)
	}
	return result
}

func Reduce[T any, R any](s []T, initial R, reducer func(R, T) R) R {
	acc := initial
	for _, v := range s {
		acc = reducer(acc, v)
	}
	return acc
}

func Contains[T comparable](s []T, value T) bool {
	return slices.Contains(s, value)
}

func IndexOf[T comparable](s []T, value T) int {
	for i, v := range s {
		if v == value {
			return i
		}
	}
	return -1
}

func Reversed[T any](s []T) []T {
	result := make([]T, len(s))
	for i, v := range s {
		result[len(s)-1-i] = v
	}
	return result
}

func Clone[T any](s []T) []T {
	result := make([]T, len(s))
	copy(result, s)
	return result
}

func Unique[T comparable](s []T) []T {
	seen := map[T]struct{}{}
	result := []T{}
	for _, v := range s {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

func Sort[S ~[]T, T constraint.Comparable](s S) {
	slices.Sort(s)
}

func Sorted[S ~[]T, T constraint.Comparable](s S) S {
	base := Clone(s)
	slices.Sort(base)
	return base
}

func SortBy[S ~[]E, E any](s S, function func(a, b E) int) {
	slices.SortFunc(s, function)
}

func SortedBy[S ~[]E, E any](s S, function func(a, b E) int) S {
	base := Clone(s)
	SortBy(base, function)
	return base
}

func Clear[T any](s *[]T) {
	*s = (*s)[:0]
}
