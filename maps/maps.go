package maps

import (
	std "maps"
)

// Clone returns a copy of m.
func Clone[M ~map[K]V, K comparable, V any](source M) M {
	return std.Clone(source)
}

// Copy copies all key/value pairs in src adding them to dst.
func Copy[L ~map[K]V, R ~map[K]V, K comparable, V any](destination L, source R) {
	for k, v := range source {
		destination[k] = v
	}
}
