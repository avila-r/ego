package set

type TreeSet[E comparable] struct {
	data []E
	less func(a, b E) bool
}

func NewTreeSet[E comparable](less func(a, b E) bool) *TreeSet[E] {
	return &TreeSet[E]{
		data: make([]E, 0),
		less: less,
	}
}

func (s *TreeSet[E]) binarySearch(element E) (int, bool) {
	left, right := 0, len(s.data)
	for left < right {
		mid := (left + right) / 2
		switch {
		case s.less(s.data[mid], element):
			left = mid + 1
		case s.less(element, s.data[mid]):
			right = mid
		default:
			return mid, true
		}
	}
	return left, false
}

func (s *TreeSet[E]) Add(element E) bool {
	idx, found := s.binarySearch(element)
	if found {
		return false
	}
	s.data = append(s.data[:idx], append([]E{element}, s.data[idx:]...)...)
	return true
}

func (s *TreeSet[E]) Remove(element E) bool {
	idx, found := s.binarySearch(element)
	if !found {
		return false
	}
	s.data = append(s.data[:idx], s.data[idx+1:]...)
	return true
}

func (s *TreeSet[E]) Contains(element E) bool {
	_, found := s.binarySearch(element)
	return found
}

func (s *TreeSet[E]) Size() int {
	return len(s.data)
}

func (s *TreeSet[E]) IsEmpty() bool {
	return len(s.data) == 0
}

func (s *TreeSet[E]) Clear() {
	s.data = make([]E, 0)
}

func (s *TreeSet[E]) ToSlice() []E {
	result := make([]E, len(s.data))
	copy(result, s.data)
	return result
}

func (s *TreeSet[E]) Union(other Settable[E]) Settable[E] {
	result := NewTreeSet[E](s.less)
	for _, v := range s.data {
		result.Add(v)
	}
	for _, v := range other.ToSlice() {
		result.Add(v)
	}
	return result
}

func (s *TreeSet[E]) Intersection(other Settable[E]) Settable[E] {
	result := NewTreeSet[E](s.less)
	for _, v := range s.data {
		if other.Contains(v) {
			result.Add(v)
		}
	}
	return result
}

func (s *TreeSet[E]) Difference(other Settable[E]) Settable[E] {
	result := NewTreeSet[E](s.less)
	for _, v := range s.data {
		if !other.Contains(v) {
			result.Add(v)
		}
	}
	return result
}
