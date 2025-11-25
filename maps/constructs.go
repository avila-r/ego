package maps

import (
	"github.com/avila-r/ego/collection"
)

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

func OfLinked[K comparable, V any](entries ...collection.Entry[K, V]) collection.Map[K, V] {
	m := NewLinked[K, V]()
	for _, entry := range entries {
		m.Put(entry.Key, entry.Value)
	}
	return m
}
