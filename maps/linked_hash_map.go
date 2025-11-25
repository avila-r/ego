package maps

import (
	"reflect"

	"github.com/avila-r/ego/collection"
	"github.com/avila-r/ego/iterator"
)

type linkedNode[K comparable, V any] struct {
	key   K
	value V
	next  *linkedNode[K, V]
	prev  *linkedNode[K, V]
}

type LinkedHashMap[K comparable, V any] struct {
	elements map[K]*linkedNode[K, V]
	head     *linkedNode[K, V]
	tail     *linkedNode[K, V]
	size     int
}

var _ collection.Map[string, int] = (*LinkedHashMap[string, int])(nil)

func NewLinkedHashMap[K comparable, V any]() *LinkedHashMap[K, V] {
	return &LinkedHashMap[K, V]{
		elements: make(map[K]*linkedNode[K, V]),
		head:     nil,
		tail:     nil,
		size:     0,
	}
}

func EmptyLinkedHashMap[K comparable, V any]() *LinkedHashMap[K, V] {
	return NewLinkedHashMap[K, V]()
}

func (m *LinkedHashMap[K, V]) Get(key K) (V, bool) {
	if node, exists := m.elements[key]; exists {
		return node.value, true
	}
	var zero V
	return zero, false
}

func (m *LinkedHashMap[K, V]) Put(key K, value V) {
	if node, exists := m.elements[key]; exists {
		// Update existing node
		node.value = value
		return
	}

	// Create new node
	newNode := &linkedNode[K, V]{
		key:   key,
		value: value,
	}

	m.elements[key] = newNode

	// Add to linked list
	if m.tail == nil {
		m.head = newNode
		m.tail = newNode
	} else {
		m.tail.next = newNode
		newNode.prev = m.tail
		m.tail = newNode
	}

	m.size++
}

func (m *LinkedHashMap[K, V]) PutIfAbsent(key K, value V) bool {
	if _, exists := m.elements[key]; exists {
		return false
	}
	m.Put(key, value)
	return true
}

func (m *LinkedHashMap[K, V]) Delete(key K) {
	node, exists := m.elements[key]
	if !exists {
		return
	}

	// Remove from linked list
	if node.prev != nil {
		node.prev.next = node.next
	} else {
		m.head = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	} else {
		m.tail = node.prev
	}

	delete(m.elements, key)
	m.size--
}

func (m *LinkedHashMap[K, V]) Clear() {
	m.elements = make(map[K]*linkedNode[K, V])
	m.head = nil
	m.tail = nil
	m.size = 0
}

func (m *LinkedHashMap[K, V]) Len() int {
	return m.size
}

func (m *LinkedHashMap[K, V]) IsEmpty() bool {
	return m.size == 0
}

func (m *LinkedHashMap[K, V]) ContainsKey(key K) bool {
	_, exists := m.elements[key]
	return exists
}

func (m *LinkedHashMap[K, V]) ContainsValue(value V) bool {
	current := m.head
	for current != nil {
		if reflect.DeepEqual(current.value, value) {
			return true
		}
		current = current.next
	}
	return false
}

func (m *LinkedHashMap[K, V]) Filter(predicate func(K, V) bool) collection.Map[K, V] {
	filtered := NewLinkedHashMap[K, V]()
	current := m.head
	for current != nil {
		if predicate(current.key, current.value) {
			filtered.Put(current.key, current.value)
		}
		current = current.next
	}
	return filtered
}

func (m *LinkedHashMap[K, V]) Clone() collection.Map[K, V] {
	cloned := NewLinkedHashMap[K, V]()
	current := m.head
	for current != nil {
		cloned.Put(current.key, current.value)
		current = current.next
	}
	return cloned
}

func (m *LinkedHashMap[K, V]) ToSlice() []collection.Entry[K, V] {
	entries := make([]collection.Entry[K, V], 0, m.size)
	current := m.head
	for current != nil {
		entries = append(entries, collection.Entry[K, V]{
			Key:   current.key,
			Value: current.value,
		})
		current = current.next
	}
	return entries
}

func (m *LinkedHashMap[K, V]) KeySlice() []K {
	keys := make([]K, 0, m.size)
	current := m.head
	for current != nil {
		keys = append(keys, current.key)
		current = current.next
	}
	return keys
}

func (m *LinkedHashMap[K, V]) ValueSlice() []V {
	values := make([]V, 0, m.size)
	current := m.head
	for current != nil {
		values = append(values, current.value)
		current = current.next
	}
	return values
}

func (m *LinkedHashMap[K, V]) Keys() collection.Collection[K] {
	return collection.New(m.KeySlice()...)
}

func (m *LinkedHashMap[K, V]) Values() collection.Collection[V] {
	return collection.New(m.ValueSlice()...)
}

func (m *LinkedHashMap[K, V]) Elements() map[K]V {
	result := make(map[K]V, m.size)
	current := m.head
	for current != nil {
		result[current.key] = current.value
		current = current.next
	}
	return result
}

func (m *LinkedHashMap[K, V]) Entries() collection.Collection[collection.Entry[K, V]] {
	return collection.New(m.ToSlice()...)
}

func (m *LinkedHashMap[K, V]) Iterator() iterator.Iterator[collection.Entry[K, V]] {
	return iterator.From(m.Entries())
}
