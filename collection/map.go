package collection

import "github.com/avila-r/ego/iterator"

type Map[K comparable, V any] interface {
	Get(key K) (V, bool)
	Put(key K, value V)
	PutIfAbsent(key K, value V) bool
	Delete(key K)
	Clear()
	Len() int
	IsEmpty() bool
	ContainsKey(key K) bool
	ContainsValue(value V) bool
	Filter(func(K, V) bool) Map[K, V]

	Clone() Map[K, V]
	ToSlice() []Entry[K, V]

	KeySlice() []K
	ValueSlice() []V

	Keys() Collection[K]
	Values() Collection[V]
	Elements() map[K]V
	Entries() Collection[Entry[K, V]]

	Iterator() iterator.Iterator[Entry[K, V]]
}

type Entry[K comparable, V any] struct {
	Key   K
	Value V
}
