package collection

type Map[K comparable, V any] interface {
	Keys() []K
	Values() []V
	Elements() map[K]V
	Entries() Collection[Entry[K, V]]
}

type Entry[K comparable, V any] struct {
	Key   K
	Value V
}
