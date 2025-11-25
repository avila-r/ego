package set

type Set[E comparable] struct {
	elements []E
}

func NewSet[E comparable](elements []E) *Set[E] {
	set := &Set[E]{elements: make([]E, 0)}
	for _, element := range elements {
		set.Add(element)
	}
	return set
}

func (s *Set[E]) Size() int {
	return len(s.elements)
}

func (s *Set[E]) IsEmpty() bool {
	return s.Size() == 0
}

func (s *Set[E]) Add(element E) bool {
	if !s.Contains(element) {
		s.elements = append(s.elements, element)
		return true
	}
	return false
}

func (s *Set[E]) AddAll(elements ...E) {
	for _, element := range elements {
		s.Add(element)
	}
}

func (s *Set[E]) Clear() {
	s.elements = []E{}
}

func (s *Set[E]) Contains(element E) bool {
	for _, e := range s.elements {
		if e == element {
			return true
		}
	}
	return false
}

func (s *Set[E]) ContainsAll(elements ...E) bool {
	for _, element := range elements {
		if !s.Contains(element) {
			return false
		}
	}
	return true
}

func (s *Set[E]) Remove(element E) bool {
	for i, e := range s.elements {
		if e == element {
			s.elements = append(s.elements[:i], s.elements[i+1:]...)
			return true
		}
	}
	return false
}

func (s *Set[E]) RemoveAll(elements ...E) int {
	removed := 0
	for _, element := range elements {
		if s.Remove(element) {
			removed++
		}
	}
	return removed
}

func (s *Set[E]) RetainAll(elements ...E) int {
	retained := make([]E, 0)
	removed := 0

	for _, element := range s.elements {
		keep := false
		keep = s.retainIfPresent(elements, element, keep)

		switch {
		case keep:
			retained = append(retained, element)
		default:
			removed++
		}
	}

	s.elements = retained
	return removed
}

func (s *Set[E]) retainIfPresent(elements []E, element E, keep bool) bool {
	for _, toRetain := range elements {
		if element == toRetain {
			keep = true
			break
		}
	}
	return keep
}

func (s *Set[E]) Equals(other E) bool {
	return s.Size() == 1 && s.Contains(other)
}

func (s *Set[E]) ToSlice() []E {
	result := make([]E, len(s.elements))
	copy(result, s.elements)
	return result
}

func (s *Set[E]) Union(other Settable[E]) Settable[E] {
	result := NewSet([]E{})
	for _, element := range s.elements {
		result.Add(element)
	}
	for _, element := range other.ToSlice() {
		result.Add(element)
	}
	return result
}

func (s *Set[E]) Intersection(other Settable[E]) Settable[E] {
	result := NewSet([]E{})
	for _, element := range s.elements {
		if other.Contains(element) {
			result.Add(element)
		}
	}
	return result
}

func (s *Set[E]) Difference(other Settable[E]) Settable[E] {
	result := NewSet([]E{})
	for _, element := range s.elements {
		if !other.Contains(element) {
			result.Add(element)
		}
	}
	return result
}
