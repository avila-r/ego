package maps

import (
	std "maps"

	"github.com/avila-r/ego/collection"
	"github.com/avila-r/ego/iterator"
)

type Map[K comparable, V any] map[K]V

// Equal reports whether two maps contain the same key/value pairs.
func Equal[L, R ~map[K]V, K, V comparable](l L, r R) bool {
	return std.Equal(l, r)
}

// EqualBy is like Equal, but compares values using eq.
func EqualBy[L ~map[K]First, R ~map[K]Second, K comparable, First, Second any](l L, r R, eq func(First, Second) bool) bool {
	return std.EqualFunc(l, r, eq)
}

// Clone returns a copy of m.
func Clone[M ~map[K]V, K comparable, V any](source M) M {
	return std.Clone(source)
}

// Copy copies all key/value pairs in src adding them to dst.
func Copy[L ~map[K]V, R ~map[K]V, K comparable, V any](destination L, source R) {
	std.Copy(destination, source)
}

// DeleteIf deletes any key/value pairs from m for which predicate returns true.
func DeleteIf[M ~map[K]V, K comparable, V any](m M, predicate func(K, V) bool) {
	std.DeleteFunc(m, predicate)
}

// Entries returns a Collection containing all key/value pairs from the map as Entry objects.
func Entries[M ~map[K]V, K comparable, V any](m M) collection.Collection[collection.Entry[K, V]] {
	c := collection.Sized[collection.Entry[K, V]](len(m))
	for k, v := range m {
		c.Add(collection.Entry[K, V]{Key: k, Value: v})
	}
	return c
}

func Keys[M ~map[K]V, K comparable, V any](m M) collection.Collection[K] {
	c := collection.Sized[K](len(m))
	for k := range m {
		c.Add(k)
	}
	return c
}

// Values returns a Collection containing all values from the map.
func Values[M ~map[K]V, K comparable, V any](m M) collection.Collection[V] {
	c := collection.Empty[V]()
	for _, v := range m {
		c.Add(v)
	}
	return c
}

func Iter[M ~map[K]V, K comparable, V any](m M) iterator.Iterator[collection.Entry[K, V]] {
	return Entries(m).Iterator()
}

func From[K comparable, V any](m Map[K, V]) collection.Map[K, V] {
	base := New[K, V]()
	for k, v := range m {
		base.Put(k, v)
	}
	return base
}

func New[K comparable, V any]() collection.Map[K, V] {
	return EmptyHashMap[K, V]()
}

func Empty[K comparable, V any]() collection.Map[K, V] {
	return New[K, V]()
}

func Of[K comparable, V any](entries ...collection.Entry[K, V]) collection.Map[K, V] {
	m := New[K, V]()
	for _, entry := range entries {
		m.Put(entry.Key, entry.Value)
	}
	return m
}


func NewLinked[K comparable, V any]() collection.Map[K, V] {
	return EmptyLinkedHashMap[K, V]()
}

func EmptyLinked[K comparable, V any]() collection.Map[K, V] {
	return NewLinked[K, V]()
}

func LinkedOf[K comparable, V any](entries ...collection.Entry[K, V]) collection.Map[K, V] {
	m := NewLinked[K, V]()
	for _, entry := range entries {
		m.Put(entry.Key, entry.Value)
	}
	return m
}
