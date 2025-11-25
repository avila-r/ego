package set

type Settable[E comparable] interface {
	Add(element E) bool
	Remove(element E) bool
	Contains(element E) bool
	Size() int
	IsEmpty() bool
	Clear()
	ToSlice() []E
	Union(other Settable[E]) Settable[E]
	Intersection(other Settable[E]) Settable[E]
	Difference(other Settable[E]) Settable[E]
}
