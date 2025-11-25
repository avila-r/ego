package set

type HashSet[E comparable] struct {
	data map[E]struct{}
}

func NewHashSet[E comparable]() *HashSet[E] {
	return &HashSet[E]{
		data: make(map[E]struct{}),
	}
}

func (s *HashSet[E]) Add(element E) bool {
	if _, exists := s.data[element]; exists {
		return false
	}
	s.data[element] = struct{}{}
	return true
}

func (s *HashSet[E]) Remove(element E) bool {
	if _, exists := s.data[element]; !exists {
		return false
	}
	delete(s.data, element)
	return true
}

func (s *HashSet[E]) Contains(element E) bool {
	_, exists := s.data[element]
	return exists
}

func (s *HashSet[E]) Size() int {
	return len(s.data)
}

func (s *HashSet[E]) IsEmpty() bool {
	return len(s.data) == 0
}

func (s *HashSet[E]) Clear() {
	s.data = make(map[E]struct{})
}

func (s *HashSet[E]) ToSlice() []E {
	result := make([]E, 0, len(s.data))
	for k := range s.data {
		result = append(result, k)
	}
	return result
}

func (s *HashSet[E]) Union(other Settable[E]) Settable[E] {
	result := NewHashSet[E]()
	for k := range s.data {
		result.Add(k)
	}
	for _, k := range other.ToSlice() {
		result.Add(k)
	}
	return result
}

func (s *HashSet[E]) Intersection(other Settable[E]) Settable[E] {
	result := NewHashSet[E]()
	for k := range s.data {
		if other.Contains(k) {
			result.Add(k)
		}
	}
	return result
}

func (s *HashSet[E]) Difference(other Settable[E]) Settable[E] {
	result := NewHashSet[E]()
	for k := range s.data {
		if !other.Contains(k) {
			result.Add(k)
		}
	}
	return result
}
