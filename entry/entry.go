package entry

import "github.com/avila-r/ego/collection"

func Of[K comparable, V any](key K, value V) collection.Entry[K, V] {
	return collection.Entry[K, V]{
		Key:   key,
		Value: value,
	}
}
