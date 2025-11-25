package maps

import (
	"github.com/avila-r/ego/collection"
	"github.com/avila-r/ego/iterator"
)

type HashMap[K comparable, V any] struct {
	elements map[K]V
}

var _ collection.Map[string, int] = (*HashMap[string, int])(nil)

func NewHashMap[K comparable, V any]() *HashMap[K, V] {
	return &HashMap[K, V]{
		elements: make(map[K]V),
	}
}

func EmptyHashMap[K comparable, V any]() *HashMap[K, V] {
	return NewHashMap[K, V]()
}

func (m *HashMap[K, V]) Get(key K) (V, bool) {
	val, ok := m.elements[key]
	return val, ok
}

func (m *HashMap[K, V]) Put(key K, value V) {
	m.elements[key] = value
}

func (m *HashMap[K, V]) PutIfAbsent(key K, value V) bool {
	if _, exists := m.elements[key]; !exists {
		m.elements[key] = value
		return true
	}
	return false
}

func (m *HashMap[K, V]) Delete(key K) {
	delete(m.elements, key)
}

func (m *HashMap[K, V]) Clear() {
	m.elements = make(map[K]V)
}

func (m *HashMap[K, V]) Len() int {
	return len(m.elements)
}

func (m *HashMap[K, V]) IsEmpty() bool {
	return len(m.elements) == 0
}

func (m *HashMap[K, V]) ContainsKey(key K) bool {
	_, exists := m.elements[key]
	return exists
}

func (m *HashMap[K, V]) ContainsValue(value V) bool {
	for _, v := range m.elements {
		if any(v) == any(value) {
			return true
		}
	}
	return false
}

func (m *HashMap[K, V]) Filter(predicate func(K, V) bool) collection.Map[K, V] {
	filtered := NewHashMap[K, V]()
	for k, v := range m.elements {
		if predicate(k, v) {
			filtered.Put(k, v)
		}
	}
	return filtered
}

func (m *HashMap[K, V]) Clone() collection.Map[K, V] {
	cloned := NewHashMap[K, V]()
	for k, v := range m.elements {
		cloned.Put(k, v)
	}
	return cloned
}

func (m *HashMap[K, V]) ToSlice() []collection.Entry[K, V] {
	entries := make([]collection.Entry[K, V], 0, len(m.elements))
	for k, v := range m.elements {
		entries = append(entries, collection.Entry[K, V]{Key: k, Value: v})
	}
	return entries
}

func (m *HashMap[K, V]) KeySlice() []K {
	keys := make([]K, 0, len(m.elements))
	for k := range m.elements {
		keys = append(keys, k)
	}
	return keys
}

func (m *HashMap[K, V]) ValueSlice() []V {
	values := make([]V, 0, len(m.elements))
	for _, v := range m.elements {
		values = append(values, v)
	}
	return values
}

func (m *HashMap[K, V]) Keys() collection.Collection[K] {
	return collection.New(m.KeySlice()...)
}

func (m *HashMap[K, V]) Values() collection.Collection[V] {
	return collection.New(m.ValueSlice()...)
}

func (m *HashMap[K, V]) Elements() map[K]V {
	return m.elements
}

func (m *HashMap[K, V]) Entries() collection.Collection[collection.Entry[K, V]] {
	return collection.New(m.ToSlice()...)
}

func (m *HashMap[K, V]) Iterator() iterator.Iterator[collection.Entry[K, V]] {
	return iterator.Of(m.ToSlice()...)
}
