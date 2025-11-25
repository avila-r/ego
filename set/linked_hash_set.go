package set

type LinkedHashSet[E comparable] struct {
	data  map[E]*node[E]
	head  *node[E]
	tail  *node[E]
	count int
}

type node[E comparable] struct {
	value E
	prev  *node[E]
	next  *node[E]
}

func NewLinkedHashSet[E comparable]() *LinkedHashSet[E] {
	return &LinkedHashSet[E]{
		data: make(map[E]*node[E]),
	}
}

func (s *LinkedHashSet[E]) Add(element E) bool {
	if _, exists := s.data[element]; exists {
		return false
	}

	n := &node[E]{value: element}
	s.data[element] = n

	if s.tail == nil {
		s.head = n
		s.tail = n
	} else {
		s.tail.next = n
		n.prev = s.tail
		s.tail = n
	}
	s.count++
	return true
}

func (s *LinkedHashSet[E]) Remove(element E) bool {
	n, exists := s.data[element]
	if !exists {
		return false
	}

	if n.prev != nil {
		n.prev.next = n.next
	} else {
		s.head = n.next
	}

	if n.next != nil {
		n.next.prev = n.prev
	} else {
		s.tail = n.prev
	}

	delete(s.data, element)
	s.count--
	return true
}

func (s *LinkedHashSet[E]) Contains(element E) bool {
	_, exists := s.data[element]
	return exists
}

func (s *LinkedHashSet[E]) Size() int {
	return s.count
}

func (s *LinkedHashSet[E]) IsEmpty() bool {
	return s.count == 0
}

func (s *LinkedHashSet[E]) Clear() {
	s.data = make(map[E]*node[E])
	s.head = nil
	s.tail = nil
	s.count = 0
}

func (s *LinkedHashSet[E]) ToSlice() []E {
	result := make([]E, 0, s.count)
	for n := s.head; n != nil; n = n.next {
		result = append(result, n.value)
	}
	return result
}

func (s *LinkedHashSet[E]) Union(other Settable[E]) Settable[E] {
	result := NewLinkedHashSet[E]()
	for n := s.head; n != nil; n = n.next {
		result.Add(n.value)
	}
	for _, v := range other.ToSlice() {
		result.Add(v)
	}
	return result
}

func (s *LinkedHashSet[E]) Intersection(other Settable[E]) Settable[E] {
	result := NewLinkedHashSet[E]()
	for n := s.head; n != nil; n = n.next {
		if other.Contains(n.value) {
			result.Add(n.value)
		}
	}
	return result
}

func (s *LinkedHashSet[E]) Difference(other Settable[E]) Settable[E] {
	result := NewLinkedHashSet[E]()
	for n := s.head; n != nil; n = n.next {
		if !other.Contains(n.value) {
			result.Add(n.value)
		}
	}
	return result
}
