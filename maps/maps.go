package maps

import (
	std "maps"

	"github.com/avila-r/ego/collection"
)

type Map[K comparable, V any] map[K]V

// Clone returns a copy of m.
func Clone[M ~map[K]V, K comparable, V any](source M) M {
	return std.Clone(source)
}

// Copy copies all key/value pairs in src adding them to dst.
func Copy[L ~map[K]V, R ~map[K]V, K comparable, V any](destination L, source R) {
	std.Copy(destination, source)
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

func LinkedFrom[K comparable, V any](m Map[K, V]) collection.Map[K, V] {
	base := NewLinked[K, V]()
	for k, v := range m {
		base.Put(k, v)
	}
	return base
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
